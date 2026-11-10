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
	// ReplacedLine is the format of a line under Replaced scenarios.
	ReplacedLine = regexp.MustCompile(`^- (\d{2})\.([1-9]\d*): (.+)$`)
	title        = regexp.MustCompile(`^# Increment (\d{2}): \S`)
	folder       = regexp.MustCompile(`^(\d{2})-[^/]+$`)
)

// RequiredSections are the headings that every note has.
var RequiredSections = []string{"## Goal", "## Out of scope", "## Existing tools", "## Interface", "## Scenarios", "## Evidence", "## Rollout"}

// replacedHeading starts the optional list of replaced scenarios.
const replacedHeading = "## Replaced scenarios"

// unchangedHeading starts the optional list of scenarios whose behaviour
// the code before the increment already has.
const unchangedHeading = "## Unchanged behaviour"

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

// Replaced is one line under Replaced scenarios.
type Replaced struct {
	ID     string
	Reason string
	Line   int
}

// ReplacedScenarios returns the lines under Replaced scenarios, and the
// numbers of lines there that do not match the format.
func ReplacedScenarios(text string) (replaced []Replaced, bad []int) {
	return listed(text, replacedHeading)
}

// UnchangedScenarios returns the lines under Unchanged behaviour, and the
// numbers of lines there that do not match the format. Reason holds why
// the code before the increment already behaves so.
func UnchangedScenarios(text string) (unchanged []Replaced, bad []int) {
	return listed(text, unchangedHeading)
}

// listed returns the lines "- <nn>.<n>: <text>" under heading.
func listed(text, heading string) (list []Replaced, bad []int) {
	in := false
	Lines(text, func(n int, line string) {
		line = strings.TrimRight(line, " \r")
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ") {
			in = line == heading
			return
		}
		if !in || line == "" {
			return
		}
		m := ReplacedLine.FindStringSubmatch(line)
		if m == nil {
			bad = append(bad, n)
			return
		}
		list = append(list, Replaced{ID: m[1] + "." + m[2], Reason: m[3], Line: n})
	})
	return list, bad
}

// Check checks the note text at file, in the folder with increment number
// nn, and returns the number of scenarios. earlier maps the increment
// number of each earlier note to its scenarios.
func Check(r *report.Report, file, nn, text string, earlier map[string][]Scenario) int {
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

	replaced, badReplaced := ReplacedScenarios(text)
	for _, n := range badReplaced {
		r.Error(file, n, "bad-replaced-line", "does not match - <nn>.<n>: <reason>")
	}
	for _, rep := range replaced {
		inc, _, _ := strings.Cut(rep.ID, ".")
		if inc >= nn {
			r.Error(file, rep.Line, "unknown-scenario", "scenario "+rep.ID+" is not from an earlier increment")
			continue
		}
		known := false
		for _, s := range earlier[inc] {
			known = known || s.ID == rep.ID
		}
		if !known {
			r.Error(file, rep.Line, "unknown-scenario", "no scenario "+rep.ID+" in the note of increment "+inc)
		}
	}

	kinds := map[string]string{}
	for _, s := range scenarios {
		kinds[s.ID] = s.Kind
	}
	unchanged, badUnchanged := UnchangedScenarios(text)
	for _, n := range badUnchanged {
		r.Error(file, n, "bad-unchanged-line", "does not match - <nn>.<n>: <reason>")
	}
	for _, u := range unchanged {
		if kinds[u.ID] != "automated" {
			r.Error(file, u.Line, "unknown-scenario", "no automated scenario "+u.ID+" in this note")
		}
	}
	return len(scenarios)
}
