// Package status derives where each increment stands from the files of
// the repo: design notes, test files, review results, results files and
// retros.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/review"
)

var (
	// ResultLine is a result of a manual or demonstrated scenario.
	ResultLine = regexp.MustCompile(`^- (\d{2}\.[1-9]\d*) (passed|failed) (\d{4}-\d{2}-\d{2}): (.+)$`)
	// BreakGlassLine is a merge while a check of the pull request failed.
	BreakGlassLine = regexp.MustCompile(`^- break-glass (\d{4}-\d{2}-\d{2}) (\S+): (.+)$`)
)

// The stages of an increment, in order.
const (
	StageTests   = "tests"
	StageReview  = "review"
	StageResults = "results"
	StageRetro   = "retro"
	StageDone    = "done"
)

// Run prints the stage and the open points of each increment of the repo
// at root, and the code commit that waits for a review, and returns the
// exit code.
func Run(r *report.Report, root string) int {
	notes, err := design.ReadNotes(root, "")
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}
	byNumber := map[string][]design.Scenario{}
	for _, s := range design.InForce(notes) {
		nn, _, _ := strings.Cut(s.ID, ".")
		byNumber[nn] = append(byNumber[nn], s)
	}
	for _, n := range notes {
		increment(r, root, n, byNumber[n.Number])
	}
	h, err := review.Pending(root)
	if err != nil {
		return r.Abort(root, "git-failed", err.Error())
	}
	if h != "" {
		r.Line("review-open " + h)
	}
	return r.Finish(len(notes))
}

// increment prints the lines of one increment.
func increment(r *report.Report, root string, n design.Note, scenarios []design.Scenario) {
	var open []string
	stage := ""
	reach := func(s string) {
		if stage == "" {
			stage = s
		}
	}

	for _, s := range scenarios {
		if s.Kind != "automated" {
			continue
		}
		test := filepath.Join(root, "test", "acceptance", n.Folder, "scenario-"+s.ID+".txtar")
		if _, err := os.Stat(test); err != nil {
			open = append(open, "missing-test "+s.ID)
			reach(StageTests)
		}
	}

	data, err := os.ReadFile(filepath.Join(root, "evidence", n.Folder, "review.md"))
	if err != nil {
		open = append(open, "missing-review")
		reach(StageReview)
	} else if res := review.Parse(string(data)); res.Valid && res.Blocking > 0 {
		open = append(open, fmt.Sprintf("blocking-review %d", res.Blocking))
		reach(StageReview)
	}

	latest := results(r, root, n, scenarios)
	for _, s := range scenarios {
		if s.Kind == "automated" {
			continue
		}
		switch latest[s.ID] {
		case "":
			open = append(open, "waits-for-engineer "+s.ID+" "+s.Kind)
			reach(StageResults)
		case "failed":
			open = append(open, "failed "+s.ID)
			reach(StageResults)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "docs", "increments", n.Folder, "retro.md")); err != nil {
		open = append(open, "missing-retro")
		reach(StageRetro)
	}
	reach(StageDone)

	r.Line(n.Folder + " " + stage)
	for _, o := range open {
		r.Line(n.Folder + " " + o)
	}
}

// results reads the results file of the increment, reports its malformed
// lines, and returns the latest result of each scenario.
func results(r *report.Report, root string, n design.Note, scenarios []design.Scenario) map[string]string {
	latest := map[string]string{}
	rel := "evidence/" + n.Folder + "/results.md"
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return latest
	}
	kinds := map[string]string{}
	for _, s := range scenarios {
		kinds[s.ID] = s.Kind
	}
	design.Lines(string(data), func(i int, line string) {
		line = strings.TrimRight(line, " \r")
		if !strings.HasPrefix(line, "- ") {
			return
		}
		if BreakGlassLine.MatchString(line) {
			return
		}
		m := ResultLine.FindStringSubmatch(line)
		if m == nil {
			r.Error(rel, i, "bad-result-line", "is not - <id> passed|failed <yyyy-mm-dd>: <text>, or - break-glass <yyyy-mm-dd> <link>: <reason>")
			return
		}
		if k := kinds[m[1]]; k != "manual" && k != "demonstrated" {
			r.Error(rel, i, "unknown-scenario", m[1])
			return
		}
		latest[m[1]] = m[2]
	})
	return latest
}
