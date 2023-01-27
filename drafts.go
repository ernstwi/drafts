package drafts

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// -----------------------------------------------------------------------------

type Draft struct {
	UUID    string
	Content string
}

func (d *Draft) String() string {
	sep := '│'
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return fmt.Sprintf("%s %c %s", d.UUID, sep, d.Content)
	}

	width, _, err := term.GetSize(fd)
	if err != nil {
		log.Fatal(err)
	}
	if width < 45 {
		return d.UUID
	}

	// Best effort hope that len([]rune) covers utf8 better than len(string)
	r := []rune(d.Content)
	if len(r) > width-39 {
		r = r[:width-39-3]
		return fmt.Sprintf("%s %c %s...", d.UUID, sep, string(r))
	}
	return fmt.Sprintf("%s %c %s", d.UUID, sep, d.Content)
}

// -----------------------------------------------------------------------------

type Folder int

const (
	FolderInbox Folder = iota
	FolderArchive
)

func (f Folder) String() string {
	return [...]string{"inbox", "archive"}[f]
}

// -----------------------------------------------------------------------------

type Filter int

const (
	FilterInbox Filter = iota
	FilterFlagged
	FilterArchive
	FilterTrash
	FilterAll
)

func (f Filter) String() string {
	return [...]string{"inbox", "flagged", "archive", "trash", "all"}[f]
}

// -----------------------------------------------------------------------------

type Sort int

const (
	SortCreated Sort = iota
	SortModified
	SortAccessed
)

func (s Sort) String() string {
	return [...]string{"created", "modified", "accessed"}[s]
}

// -----------------------------------------------------------------------------

type CreateOptions struct {
	Tags       []string
	Folder     Folder
	Flagged    bool
	Action     string
	AllowEmpty bool
	// Omitted: RetParam
}

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
	SortDescending   bool
	SortFlaggedToTop bool
}

//go:embed query.js
var queryjs string

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

// -----------------------------------------------------------------------------

//go:embed trash.js
var trashjs string

// Trash a draft.
func Trash(uuid string) {
	JS(trashjs, uuid)
}

// -----------------------------------------------------------------------------

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

// Run action with `text` without creating a new draft.
// https://docs.getdrafts.com/docs/automation/urlschemes#runaction
func RunAction(action, text string) url.Values {
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

	// Send socket address to open(), which sends it to the callback handler via
	// Drafts. The callback handler then uses the socket address to forward the
	// reply from Drafts to open(). This also signals to open() that the server
	// is ready to accept connections.
	ch <- sockAddr

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

func mustJSON(a any) string {
	js, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	return string(js)
}
