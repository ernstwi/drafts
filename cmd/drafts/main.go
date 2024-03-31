package main

import (
	"fmt"
	"strings"

	arg "github.com/alexflint/go-arg"

	"github.com/ernstwi/drafts/pkg/drafts"
)

const linebreak = " ~ "

// ---- Commands ---------------------------------------------------------------

type NewCmd struct {
	Message string   `arg:"positional" help:"draft content (omit to use stdin)"`
	Tag     []string `arg:"-t,separate" help:"tag"`
	Archive bool     `arg:"-a" help:"create draft in archive"`
	Flagged bool     `arg:"-f" help:"create flagged draft"`
}

func new(param *NewCmd) string {
	// Input
	text := orStdin(param.Message)

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

type PrependCmd struct {
	Message string `arg:"positional" help:"text to prepend (omit to use stdin)"`
	UUID    string `arg:"-u" help:"UUID (omit to use active draft)"`
}

func prepend(param *PrependCmd) string {
	text := orStdin(param.Message)
	uuid := orActive(param.UUID)
	drafts.Prepend(uuid, text)
	return drafts.Get(uuid).Content
}

type AppendCmd struct {
	Message string `arg:"positional" help:"text to append (omit to use stdin)"`
	UUID    string `arg:"-u" help:"UUID (omit to use active draft)"`
}

func append(param *AppendCmd) string {
	text := orStdin(param.Message)
	uuid := orActive(param.UUID)
	drafts.Append(uuid, text)
	return drafts.Get(uuid).Content
}

type ReplaceCmd struct {
	Message string `arg:"positional" help:"text to append (omit to use stdin)"`
	UUID    string `arg:"-u" help:"UUID (omit to use active draft)"`
}

func replace(param *ReplaceCmd) string {
	text := orStdin(param.Message)
	uuid := orActive(param.UUID)
	drafts.Replace(uuid, text)
	return drafts.Get(uuid).Content
}

type EditCmd struct {
	UUID string `arg:"positional" help:"UUID (omit to use active draft)"`
}

func edit(param *EditCmd) string {
	uuid := orActive(param.UUID)
	new := editor(drafts.Get(uuid).Content)
	drafts.Replace(uuid, new)
	return new
}

type GetCmd struct {
	UUID string `arg:"positional" help:"UUID (omit to use active draft)"`
}

func get(param *GetCmd) string {
	uuid := orActive(param.UUID)
	return drafts.Get(uuid).Content
}

type SelectCmd struct{}

func _select() {
	ds := drafts.Query("", drafts.FilterInbox, drafts.QueryOptions{})
	var b strings.Builder
	for _, d := range ds {
		fmt.Fprintf(&b, "%s %c %s\n", d.UUID, drafts.Separator, strings.ReplaceAll(d.Content, "\n", linebreak))
	}
	uuid := fzfUUID(b.String())
	drafts.Select(uuid)
}

// ---- Main -------------------------------------------------------------------

func main() {
	var args struct {
		New     *NewCmd     `arg:"subcommand:new" help:"create new draft"`
		Prepend *PrependCmd `arg:"subcommand:prepend" help:"prepend to draft"`
		Append  *AppendCmd  `arg:"subcommand:append" help:"append to draft"`
		Replace *ReplaceCmd `arg:"subcommand:replace" help:"append to draft"`
		Edit    *EditCmd    `arg:"subcommand:edit" help:"edit draft in vim"`
		Get     *GetCmd     `arg:"subcommand:get" help:"get content of draft"`
		Select  *SelectCmd  `arg:"subcommand:select" help:"select active draft using fzf"`
	}
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}
	switch {
	case args.New != nil:
		fmt.Println(new(args.New))
	case args.Prepend != nil:
		fmt.Println(prepend(args.Prepend))
	case args.Append != nil:
		fmt.Println(append(args.Append))
	case args.Replace != nil:
		fmt.Println(replace(args.Replace))
	case args.Edit != nil:
		fmt.Println(edit(args.Edit))
	case args.Get != nil:
		fmt.Println(get(args.Get))
	case args.Select != nil:
		_select()
	}
}
