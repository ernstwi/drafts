package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"

	"github.com/ernstwi/drafts"
)

type NewCmd struct {
	Message string   `arg:"-m" help:"draft content (omit to use stdin)"`
	Tag     []string `arg:"-t,separate" help:"tag"`
	Archive bool     `arg:"-a" help:"create draft in archive"`
	Flagged bool     `arg:"-f" help:"create flagged draft"`
	Action  string   `help:"run action on input as draft"`
	// Omitted: AllowEmpty
}

type GetCmd struct {
	UUID string `arg:"positional" help:"UUID (omit to read from stdin)"`
}

type QueryCmd struct {
	QueryString string   `arg:"positional" help:"search string"`
	Filter      string   `arg:"-f" help:"filter: inbox, flagged, archive, or trash" default:"all"`
	Tags        []string `arg:"-t,--tag,separate" help:"filter by tag"`
	OmitTags    []string `arg:"-T,--omit-tag,separate" help:"filter out by tag"`
	Sort        string   `arg:"-s" help:"sort method: created, modified, or accessed"`
	// TODO: Replace with SortAscending? Make descending default.
	SortDescending   bool `arg:"-d,--descending" help:"sort descending"`
	SortFlaggedToTop bool `arg:"-F,--flagged-first" help:"sort flagged drafts to top"`
}

func main() {
	var args struct {
		New   *NewCmd   `arg:"subcommand:new" help:"create new draft"`
		Get   *GetCmd   `arg:"subcommand:get" help:"get content of draft"`
		Query *QueryCmd `arg:"subcommand:query" help:"search for drafts"`
	}
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}
	switch {
	case args.New != nil:
		fmt.Println(new(p, args.New))
	case args.Get != nil:
		fmt.Println(get(p, args.Get.UUID))
	case args.Query != nil:
		fmt.Println(strings.Join(query(p, args.Query), "\n"))
	}
}

func new(p *arg.Parser, param *NewCmd) string {
	// Input
	text := param.Message
	if text == "" {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		text = string(stdin)
	}

	// Params -> Options
	opt := drafts.CreateOptions{
		Tags:    param.Tag,
		Flagged: param.Flagged,
		Action:  param.Action,
	}

	if param.Archive {
		opt.Folder = drafts.FolderArchive
	}

	uuid := drafts.Create(text, opt)
	return uuid
}

func get(p *arg.Parser, uuid string) string {
	if uuid == "" {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		uuid = strings.TrimSpace(string(stdin))
	}
	return drafts.Get(uuid)
}

func query(p *arg.Parser, param *QueryCmd) []string {
	opt := drafts.QueryOptions{
		Tags:             param.Tags,
		OmitTags:         param.OmitTags,
		SortDescending:   param.SortDescending,
		SortFlaggedToTop: param.SortFlaggedToTop,
	}

	// TODO: Custom parsing
	// https://github.com/alexflint/go-arg#custom-parsing

	var filter drafts.Filter
	switch param.Filter {
	case "all":
		filter = drafts.FilterAll
	case "inbox":
		filter = drafts.FilterInbox
	case "archive":
		filter = drafts.FilterArchive
	case "trash":
		filter = drafts.FilterTrash
	default:
		p.Fail("filter must be inbox, flagged, archive, or trash")
	}

	switch param.Sort {
	case "":
	case "created":
		opt.Sort = drafts.SortCreated
	case "modified":
		opt.Sort = drafts.SortModified
	case "accessed":
		opt.Sort = drafts.SortAccessed
	default:
		p.Fail("sort must be created, modified, or accessed")
	}

	uuids := drafts.Query(param.QueryString, filter, opt)
	return uuids
}

func contains[T comparable](s []T, x T) bool {
	for _, y := range s {
		if y == x {
			return true
		}
	}
	return false
}
