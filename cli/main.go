package main

import (
	"fmt"
	"io"
	"log"
	"os"

	arg "github.com/alexflint/go-arg"
	"github.com/ktr0731/go-fuzzyfinder"

	"github.com/ernstwi/drafts"
)

type NewCmd struct {
	Message string   `arg:"-m" help:"draft content (omit to use stdin)"`
	Tag     []string `arg:"-t,separate" help:"tag"`
	Archive bool     `arg:"-a" help:"create draft in archive"`
	Flagged bool     `arg:"-f" help:"create flagged draft"`
}

type GetCmd struct {
	UUID string `arg:"positional" help:"UUID (omit to use active draft)"`
}

type SelectCmd struct{}

func main() {
	var args struct {
		New    *NewCmd    `arg:"subcommand:new" help:"create new draft"`
		Get    *GetCmd    `arg:"subcommand:get" help:"get content of draft"`
		Select *SelectCmd `arg:"subcommand:select" help:"select active draft using fzf"`
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
	case args.Select != nil:
		_select()
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
	}

	if param.Archive {
		opt.Folder = drafts.FolderArchive
	}

	uuid := drafts.Create(text, opt)
	return uuid
}

func get(p *arg.Parser, uuid string) string {
	if uuid == "" {
		return drafts.Get(drafts.Active())
	}
	return drafts.Get(uuid)
}

func _select() {
	ds := drafts.Query("", drafts.FilterInbox, drafts.QueryOptions{})
	idx, _ := fuzzyfinder.Find(
		ds,
		func(i int) string {
			return ds[i].Content
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return ds[i].Content
		}))
	drafts.Load(ds[idx].UUID)
}
