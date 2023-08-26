package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ernstwi/drafts"
)

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

// Run FZF on input, return UUID.
func fzfUUID(input string) (string, error) {
	line, err := fzf(input)
	if err != nil {
		return "", err
	}
	return strings.Split(line, fmt.Sprintf(" %c ", drafts.Separator))[0], nil
}

// Run FZF on input, return line.
func fzf(input string) (string, error) {
	var result strings.Builder
	cmd := exec.Command("fzf", "--delimiter", "\\|", "--with-nth", "2")
	cmd.Stdout = &result
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input)

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result.String()), nil
}

func editor(input string) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte(input)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		log.Fatal(err)
	}
	return string(data), nil
}
