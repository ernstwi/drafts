package drafts

import (
	"encoding/json"
	"net/url"
)

// ---- Writing drafts ---------------------------------------------------------

// Create a new draft. Return new draft's UUID.
// https://docs.getdrafts.com/docs/automation/urlschemes#create
func Create(text string, opt CreateOptions) string {
	v := url.Values{
		"text":    []string{text},
		"folder":  []string{opt.Folder.String()},
		"flagged": []string{mustJSON(opt.Flagged)},
	}
	if len(opt.Tags) > 0 {
		v["tag"] = opt.Tags
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

// Update content of an existing draft.
func Update(uuid, text string) {
	// replaceRange URL requires a range, so using JS is simpler
	JS(updatejs, uuid, text)
}

// Trash a draft.
func Trash(uuid string) {
	JS(trashjs, uuid)
}

// TODO:
// - Archive
// - Tag

// ---- Reading drafts ---------------------------------------------------------

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

// ---- App state --------------------------------------------------------------

// Set active draft.
func Select(uuid string) {
	JS(loadjs, uuid)
}

// Get UUID of active draft.
func Active() string {
	res := open("getCurrentDraft", url.Values{})
	return res.Get("uuid")
}

// ---- Helpers ----------------------------------------------------------------

// Run action with `text` without creating a new draft.
// TODO: Add option to run on s Draft (using "open" URL)
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
