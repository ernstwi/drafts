package drafts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// -----------------------------------------------------------------------------

type Folder int

const (
	Inbox Folder = iota
	Archive
)

func (f Folder) String() string {
	return [...]string{"inbox", "archive"}[f]
}

// -----------------------------------------------------------------------------

type Sort int

const (
	Created Sort = iota
	Modified
	Accessed
)

func (s Sort) String() string {
	return [...]string{"created", "modified", "accessed"}[s]
}

// -----------------------------------------------------------------------------

// TODO: Use json.Marshal on regular bool instead
type Bool bool

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

// -----------------------------------------------------------------------------

type CreateOptions struct {
	Tags       []string
	Folder     Folder
	Flagged    Bool
	Action     string
	AllowEmpty Bool
	// Omitted: RetParam
}

// Create a new draft. Return new draft's UUID.
// https://docs.getdrafts.com/docs/automation/urlschemes#create
func Create(text string, opt CreateOptions) string {
	v := url.Values{
		"text":       []string{text},
		"folder":     []string{opt.Folder.String()},
		"flagged":    []string{opt.Flagged.String()},
		"allowEmpty": []string{opt.AllowEmpty.String()},
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

// -----------------------------------------------------------------------------

// Get content of draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#get
func Get(uuid string) string {
	res := open("get", url.Values{
		"uuid": []string{uuid},
	})
	return res.Get("text")
}

// -----------------------------------------------------------------------------

type QueryOptions struct {
	Tags             []string
	OmitTags         []string
	Sort             Sort
	SortDescending   Bool
	SortFlaggedToTop Bool
}

// Query for drafts, return UUIDs.
// https://scripting.getdrafts.com/classes/Draft#query
func Query(queryString, filter string, opt QueryOptions) []string {
	args := []any{
		queryString,
		filter,
		opt.Tags,
		opt.OmitTags,
		opt.Sort.String(),
		opt.SortDescending,
		opt.SortFlaggedToTop,
	}
	json, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}
	return queryJS(string(json))
}

// Query for drafts using JS, return UUIDs.
// https://scripting.getdrafts.com/classes/Draft#query
// NOTE: Minimum required params are `queryString` and `filter`.
func queryJS(params string) []string {
	v := RunAction(params, "query")
	if v.Has("uuids") {
		return strings.Split(v.Get("uuids"), ",")
	}
	return []string{}
}

// -----------------------------------------------------------------------------

// Trash a draft.
func Trash(uuid string) {
	json, err := json.Marshal(uuid)
	if err != nil {
		log.Fatal(err)
	}
	RunAction(string(json), "trash")
}

// -----------------------------------------------------------------------------

// Run action with `text` without creating a new draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#runaction
func RunAction(text, action string) url.Values {
	res := open("runAction", url.Values{
		"text":   []string{text},
		"action": []string{action},
	})
	return res
}

// -----------------------------------------------------------------------------

func open(action string, v url.Values) url.Values {
	ch := make(chan string)
	go server(ch)
	sockAddr := <-ch // Wait for ready signal
	v.Add("x-success", "ernst://"+sockAddr)
	err := exec.Command("open", "-g", draftsURL(action, v)).Run()
	if err != nil {
		log.Fatal(err)
	}

	return urlValues(<-ch)
}

func urlValues(urlstr string) url.Values {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	return u.Query()
}

func draftsURL(action string, v url.Values) string {
	return fmt.Sprintf("drafts://x-callback-url/%s?%s", action, strings.ReplaceAll(v.Encode(), "+", "%20"))
}

// Start a server, listen for one message, send it over ch
func server(ch chan string) {
	// Create a temp file to use as socket address
	f, err := os.CreateTemp("", "*.sock")
	if err != nil {
		log.Fatal(err)
	}
	sockAddr := f.Name()

	// We don't actually need the file, just the filename
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	os.Remove(f.Name())

	// To delete the socket after communication is done
	defer os.Remove(f.Name())

	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	ch <- sockAddr // Signal ready

	c, err := l.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}

	msg, err := io.ReadAll(c)
	if err != nil {
		log.Fatal(err)
	}

	ch <- string(msg)
}
