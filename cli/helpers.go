package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ernstwi/drafts"
)

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
	cmd := exec.Command("fzf")
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
