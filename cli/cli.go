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
	Message string   `arg:"positional" help:"draft content (omit to use stdin)"`
	Tag     []string `arg:"-t,separate" help:"tag"`
	Archive bool     `arg:"-a" help:"create draft in archive"`
	Flagged bool     `arg:"-f" help:"create flagged draft"`
}

type PrependCmd struct {
	Message string `arg:"positional" help:"text to prepend (omit to use stdin)"`
	UUID    string `arg:"-u" help:"UUID (omit to use active draft)"`
}

type AppendCmd struct {
	Message string `arg:"positional" help:"text to append (omit to use stdin)"`
	UUID    string `arg:"-u" help:"UUID (omit to use active draft)"`
}

// TODO: Add `update`, then `edit`

type GetCmd struct {
	UUID string `arg:"positional" help:"UUID (omit to use active draft)"`
}

type SelectCmd struct{}

const linebreak = " ~ "

func main() {
	var args struct {
		New     *NewCmd     `arg:"subcommand:new" help:"create new draft"`
		Prepend *PrependCmd `arg:"subcommand:prepend" help:"prepend to draft"`
		Append  *AppendCmd  `arg:"subcommand:append" help:"append to draft"`
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
	case args.Get != nil:
		fmt.Println(get(args.Get))
	case args.Select != nil:
		_select()
	}
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

func prepend(param *PrependCmd) string {
	text := orStdin(param.Message)
	uuid := orActive(param.UUID)
	drafts.Prepend(uuid, text)
	return drafts.Get(uuid).Content
}

func append(param *AppendCmd) string {
	text := orStdin(param.Message)
	uuid := orActive(param.UUID)
	drafts.Append(uuid, text)
	return drafts.Get(uuid).Content
}

func get(param *GetCmd) string {
	uuid := orActive(param.UUID)
	return drafts.Get(uuid).Content
}

func _select() {
	ds := drafts.Query("", drafts.FilterInbox, drafts.QueryOptions{})
	var b strings.Builder
	for _, d := range ds {
		fmt.Fprintf(&b, "%s %c %s\n", d.UUID, drafts.Separator, strings.Replace(d.Content, "\n", linebreak, -1))
	}
	uuid, err := fzfUUID(b.String())
	if err != nil {
		log.Fatal(err)
	}
	drafts.Select(uuid)
}

// --- Helpers -----------------------------------------------------------------

func orStdin(text string) string {
	if text != "" {
		return text
	}
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	return string(stdin)
}

func orActive(uuid string) string {
	if uuid != "" {
		return uuid
	}
	return drafts.Active()
}
