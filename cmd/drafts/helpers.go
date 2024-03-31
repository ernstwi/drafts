package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ernstwi/drafts/pkg/drafts"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func orStdin(text string) string {
	if text != "" {
		return text
	}
	stdin, err := io.ReadAll(os.Stdin)
	fatal(err)
	return string(stdin)
}

func orActive(uuid string) string {
	if uuid != "" {
		return uuid
	}
	return drafts.Active()
}

// Run FZF on input, return UUID.
func fzfUUID(input string) string {
	line := fzf(input)
	return strings.Split(line, fmt.Sprintf(" %c ", drafts.Separator))[0]
}

// Run FZF on input, return line.
func fzf(input string) string {
	var result strings.Builder
	cmd := exec.Command("fzf", "--delimiter", "\\|", "--with-nth", "2")
	cmd.Stdout = &result
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input)

	err := cmd.Start()
	fatal(err)

	err = cmd.Wait()
	fatal(err)

	return strings.TrimSpace(result.String())
}

func editor(input string) string {
	f, err := os.CreateTemp("", "")
	fatal(err)
	defer os.Remove(f.Name()) // clean up

	_, err = f.Write([]byte(input))
	fatal(err)

	err = f.Close()
	fatal(err)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	fatal(err)

	data, err := os.ReadFile(f.Name())
	fatal(err)

	// Trim trailing newline
	res := strings.TrimSuffix(string(data), "\n")

	return res
}
