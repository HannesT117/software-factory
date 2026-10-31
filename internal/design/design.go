// Package design reads and checks design notes.
package design

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HannesT117/software-factory/internal/report"
)

var (
	// ScenarioLine is the format of a scenario line.
	ScenarioLine = regexp.MustCompile(`^- \*\*scenario (\d{2})\.([1-9]\d*)\*\* \((automated|demonstrated|manual)\) \*\*(.+?)[.?]\*\*`)
	title        = regexp.MustCompile(`^# Increment (\d{2}): \S`)
	folder       = regexp.MustCompile(`^(\d{2})-[^/]+$`)
)

// RequiredSections are the headings that every note has.
var RequiredSections = []string{"## Goal", "## Out of scope", "## Existing tools", "## Interface", "## Scenarios", "## Evidence"}

// Scenario is one scenario line of a note.
type Scenario struct {
	ID    string // for example "02.1"
	Kind  string // automated, demonstrated or manual
	Title string // without the final period or question mark
	Line  int
}

// FolderNumber returns the increment number of a folder name such as
// "02-trust-core", or false when the name has another form.
func FolderNumber(dir string) (string, bool) {
	m := folder.FindStringSubmatch(filepath.Base(dir))
	if m == nil {
		return "", false
	}
	return m[1], true
}

// Lines calls fn for each line outside fenced code blocks, with its number.
func Lines(text string, fn func(n int, line string)) {
	inFence := false
	for i, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		t := strings.TrimLeft(line, " ")
		if strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence {
			fn(i+1, line)
		}
	}
}

// Scenarios returns the scenario lines of a note, and the numbers of lines
// that start with "- **scenario" but do not match the format.
func Scenarios(text string) (scenarios []Scenario, bad []int) {
	Lines(text, func(n int, line string) {
		if !strings.HasPrefix(line, "- **scenario") {
			return
		}
		m := ScenarioLine.FindStringSubmatch(line)
		if m == nil {
			bad = append(bad, n)
			return
		}
		scenarios = append(scenarios, Scenario{ID: m[1] + "." + m[2], Kind: m[3], Title: m[4], Line: n})
	})
	return scenarios, bad
}

// Check checks the note text at file, in the folder with increment number
// nn, and returns the number of scenarios.
func Check(r *report.Report, file, nn, text string) int {
	first, _, _ := strings.Cut(text, "\n")
	if m := title.FindStringSubmatch(strings.TrimRight(first, "\r")); m == nil {
		r.Error(file, 1, "bad-title", "line 1 is not # Increment "+nn+": <title>")
	} else if m[1] != nn {
		r.Error(file, 1, "bad-title", fmt.Sprintf("increment %s in folder %s", m[1], nn))
	}

	headings := map[string]bool{}
	Lines(text, func(_ int, line string) {
		headings[strings.TrimRight(line, " \r")] = true
	})
	for _, h := range RequiredSections {
		if !headings[h] {
			r.Error(file, 0, "missing-section", h)
		}
	}

	scenarios, bad := Scenarios(text)
	for _, n := range bad {
		r.Error(file, n, "bad-scenario-line", "does not match the scenario format")
	}
	seen := map[string]bool{}
	for _, s := range scenarios {
		if !strings.HasPrefix(s.ID, nn+".") {
			r.Error(file, s.Line, "wrong-increment", "scenario "+s.ID+" in increment "+nn)
		}
		if seen[s.ID] {
			r.Error(file, s.Line, "duplicate-id", "scenario "+s.ID)
		}
		seen[s.ID] = true
	}
	return len(scenarios)
}
