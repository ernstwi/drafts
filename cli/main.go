package main

import (
	"fmt"
	"io"
	"log"
	"os"

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

func main() {
	var args struct {
		New *NewCmd `arg:"subcommand:new"`
		Get *GetCmd `arg:"subcommand:get"`
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
