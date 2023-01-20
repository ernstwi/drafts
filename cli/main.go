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
	UUID string `arg:"positional"`
}

type QueryCmd struct {
	QueryString string   `arg:"positional,required" help:"search string"`
	Filter      string   `arg:"positional" help:"filter: inbox, flagged, archive, or trash" default:"all"`
	Tags        []string `arg:"-t,--tag,separate" help:"filter by tag"`
	OmitTags    []string `arg:"-T,--omit-tag,separate" help:"filter out by tag"`
	Sort        string   `arg:"-s" help:"sort method: created, modified, or accessed"`
	// TODO: Replace with SortAscending? Make descending default.
	SortDescending   bool `arg:"-d,--descending" help:"sort descending"`
	SortFlaggedToTop bool `arg:"-f,--flagged-first" help:"sort flagged drafts to top"`
}

func main() {
	var args struct {
		New   *NewCmd   `arg:"subcommand:new"`
		Get   *GetCmd   `arg:"subcommand:get"`
		Query *QueryCmd `arg:"subcommand:query"`
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
		Flagged: drafts.Bool(param.Flagged),
		Action:  param.Action,
	}

	if param.Archive {
		opt.Folder = drafts.Archive
	}

	uuid := drafts.Create(text, opt)
	return uuid
}

func get(p *arg.Parser, uuid string) string {
	return drafts.Get(uuid)
}

func query(p *arg.Parser, param *QueryCmd) []string {
	opt := drafts.QueryOptions{
		Tags:             param.Tags,
		OmitTags:         param.OmitTags,
		SortDescending:   drafts.Bool(param.SortDescending),
		SortFlaggedToTop: drafts.Bool(param.SortFlaggedToTop),
	}

	// TODO: Custom parsing
	// https://github.com/alexflint/go-arg#custom-parsing
	if param.Filter != "all" && !contains([]string{"inbox", "flagged", "archive", "trash"}, param.Filter) {
		p.Fail("filter must be inbox, flagged, archive, or trash")
	}

	switch param.Sort {
	case "":
	case "created":
		opt.Sort = drafts.Created
	case "modified":
		opt.Sort = drafts.Modified
	case "accessed":
		opt.Sort = drafts.Accessed
	default:
		p.Fail("sort must be created, modified, or accessed")
	}

	uuids := drafts.Query(param.QueryString, param.Filter, opt)
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
