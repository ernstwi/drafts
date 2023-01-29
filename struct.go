package drafts

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

const Separator = '|'

type Draft struct {
	UUID    string
	Content string
}

func (d *Draft) String() string {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return fmt.Sprintf("%s %c %s", d.UUID, Separator, d.Content)
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
		return fmt.Sprintf("%s %c %s...", d.UUID, Separator, string(r))
	}
	return fmt.Sprintf("%s %c %s", d.UUID, Separator, d.Content)
}

// ---- Enums ------------------------------------------------------------------

type Folder int

const (
	FolderInbox Folder = iota
	FolderArchive
)

func (f Folder) String() string {
	return [...]string{"inbox", "archive"}[f]
}

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

type Sort int

const (
	SortCreated Sort = iota
	SortModified
	SortAccessed
)

func (s Sort) String() string {
	return [...]string{"created", "modified", "accessed"}[s]
}
