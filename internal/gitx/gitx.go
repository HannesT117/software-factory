// Package gitx runs git commands in a repository.
package gitx

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Repo is a git working tree.
type Repo struct {
	Dir string
}

// Run runs git with args and returns its stdout.
func (r Repo) Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return stdout.String(), nil
}

// Lines runs git with args and returns its stdout split into lines,
// without empty lines.
func (r Repo) Lines(args ...string) ([]string, error) {
	out, err := r.Run(args...)
	if err != nil {
		return nil, err
	}
	return splitNonEmpty(out, "\n"), nil
}

// Paths runs git with args, which must include -z, and returns the
// NUL-separated paths of its stdout.
func (r Repo) Paths(args ...string) ([]string, error) {
	out, err := r.Run(args...)
	if err != nil {
		return nil, err
	}
	return splitNonEmpty(out, "\x00"), nil
}

// Toplevel returns the root of the working tree that contains dir.
func Toplevel(dir string) (string, error) {
	out, err := Repo{Dir: dir}.Run("rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func splitNonEmpty(s, sep string) []string {
	var out []string
	for _, part := range strings.Split(s, sep) {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
