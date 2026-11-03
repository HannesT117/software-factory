package setup

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

// gh runs the GitHub command-line tool at path.
type gh struct {
	path string
}

// callError is a call through gh that exited with a status other than 0.
type callError struct {
	// call is the method and path of a gh api call, or "gh secret set".
	call string
	// first is the first line of gh's error output.
	first string
	// notFound is true for an HTTP 404.
	notFound bool
}

func (e *callError) Error() string {
	if e.first == "" {
		return e.call
	}
	return e.call + ": " + e.first
}

// run runs gh with args and stdin, and returns its standard output. The
// command never passes on the standard input of factory.
func (g gh) run(call string, stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command(g.path, args...)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return stdout.Bytes(), nil
	}
	first, _, _ := strings.Cut(strings.TrimSpace(stderr.String()), "\n")
	first = strings.TrimSpace(first)
	if first == "" {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			first = err.Error()
		}
	}
	return nil, &callError{call: call, first: first, notFound: strings.Contains(first, "(HTTP 404)")}
}

// api calls the REST endpoint path with method. A body other than nil goes
// to gh as JSON on standard input. The result is the decoded JSON answer,
// or nil for an empty answer.
func (g gh) api(method, path string, body any) (any, error) {
	call := method + " " + path
	args := []string{"api", "-X", method, "-H", "Accept: application/vnd.github+json", path}
	var input []byte
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, &callError{call: call, first: err.Error()}
		}
		input = data
		args = append(args, "--input", "-")
	}
	out, err := g.run(call, input, args...)
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(out)) == 0 {
		return nil, nil
	}
	dec := json.NewDecoder(bytes.NewReader(out))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, &callError{call: call, first: "the answer is no JSON: " + err.Error()}
	}
	return v, nil
}

// isNotFound reports whether err is a call that GitHub answered with 404.
func isNotFound(err error) bool {
	var ce *callError
	return errors.As(err, &ce) && ce.notFound
}
