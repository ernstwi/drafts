package drafts

import (
	"encoding/json"
	"net/url"
)

// Create a new draft. Return new draft's UUID.
// https://docs.getdrafts.com/docs/automation/urlschemes#create
func Create(text string, opt CreateOptions) string {
	v := url.Values{
		"text":       []string{text},
		"folder":     []string{opt.Folder.String()},
		"flagged":    []string{mustJSON(opt.Flagged)},
		"allowEmpty": []string{mustJSON(opt.AllowEmpty)},
	}
	if len(opt.Tags) > 0 {
		v["tag"] = opt.Tags
	}
	if opt.Action != "" {
		v.Add("action", opt.Action)
	}
	res := open("create", v)
	return res.Get("uuid")
}

// Prepend to an existing draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#prepend
func Prepend(uuid, text string) {
	open("prepend", url.Values{
		"uuid": []string{uuid},
		"text": []string{text},
	})
}

// Append to an existing draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#prepend
func Append(uuid, text string) {
	open("append", url.Values{
		"uuid": []string{uuid},
		"text": []string{text},
	})
}

// Get content of draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#get
func Get(uuid string) string {
	res := open("get", url.Values{
		"uuid": []string{uuid},
	})
	return res.Get("text")
}

// Query for drafts.
// https://scripting.getdrafts.com/classes/Draft#query
func Query(queryString string, filter Filter, opt QueryOptions) []Draft {
	args := []any{
		queryString,
		filter.String(),
		opt.Tags,
		opt.OmitTags,
		opt.Sort.String(),
		opt.SortDescending,
		opt.SortFlaggedToTop,
	}
	js := JS(queryjs, args...)
	var ds []Draft
	json.Unmarshal([]byte(js), &ds)
	return ds
}

// Trash a draft.
func Trash(uuid string) {
	JS(trashjs, uuid)
}

// Run action with `text` without creating a new draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#runaction
func RunAction(action, text string) url.Values {
	res := open("runAction", url.Values{
		"text":   []string{text},
		"action": []string{action},
	})
	return res
}

// Run JavaScript program in Drafts. Params are available as an array `input`.
// Returns any JSON added as `result` using context.addSuccessParameter.
func JS(program string, params ...any) string {
	js := mustJSON(struct {
		Program string `json:"program"`
		Input   []any  `json:"input"`
	}{
		program,
		params,
	})
	v := RunAction("Drafts CLI Helper", string(js))
	if v.Has("result") {
		return v.Get("result")
	}
	return ""
}
