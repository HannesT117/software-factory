// Package report writes the output of a factory command: one line per
// finding, each with the command's prefix, and a last line with the result.
package report

import (
	"fmt"
	"io"
)

// Exit codes of every command.
const (
	OK    = 0
	Fail  = 1
	Usage = 2
)

// Report collects the error lines of one command run.
type Report struct {
	cmd    string
	w      io.Writer
	errors int
}

// New returns a report that writes lines for cmd to w.
func New(cmd string, w io.Writer) *Report {
	return &Report{cmd: cmd, w: w}
}

// Error writes "<cmd>: error <where>: <code> <detail>". where is a path, or
// path:line when line is greater than 0.
func (r *Report) Error(where string, line int, code, detail string) {
	r.errors++
	if line > 0 {
		where = fmt.Sprintf("%s:%d", where, line)
	}
	if detail != "" {
		code += " " + detail
	}
	fmt.Fprintf(r.w, "%s: error %s: %s\n", r.cmd, where, code)
}

// Line writes "<cmd>: <text>".
func (r *Report) Line(text string) {
	fmt.Fprintf(r.w, "%s: %s\n", r.cmd, text)
}

// Errors returns the number of error lines so far.
func (r *Report) Errors() int {
	return r.errors
}

// Finish writes the last line and returns the exit code: OK when there were
// no errors, Fail otherwise.
func (r *Report) Finish(checked int) int {
	if r.errors == 0 {
		fmt.Fprintf(r.w, "%s: ok %d checked\n", r.cmd, checked)
		return OK
	}
	fmt.Fprintf(r.w, "%s: %d errors\n", r.cmd, r.errors)
	return Fail
}

// Abort writes an error line for a usage or environment problem and the
// last line, and returns Usage.
func (r *Report) Abort(where, code, detail string) int {
	r.Error(where, 0, code, detail)
	fmt.Fprintf(r.w, "%s: %d errors\n", r.cmd, r.errors)
	return Usage
}
