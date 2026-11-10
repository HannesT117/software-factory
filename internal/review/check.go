package review

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/HannesT117/software-factory/internal/branch"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/report"
)

// Rounds is the number of review rounds with a blocking finding after
// which the engineer decides about another one.
const Rounds = 2

// Check checks that the review result of the latest code commit of HEAD
// against base exists and holds no blocking finding, and returns the exit
// code.
func Check(r *report.Report, root, base string) int {
	repo := gitx.Repo{Dir: root}
	if _, err := repo.Run("rev-parse", "--verify", "--quiet", base+"^{commit}"); err != nil {
		return r.Abort(base, "unknown-ref", "no commit "+base)
	}
	code, err := ReadCode(root)
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}
	b, err := branch.Read(repo, base)
	if err != nil {
		return r.Abort(root, "git-failed", err.Error())
	}
	h := b.LatestCode(code)
	if h == "" {
		r.Line("no code change")
		return r.Finish(0)
	}
	file, err := Find(root, h)
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}
	if file == "" {
		r.Error(h, 0, "missing-review", "no review result in evidence/*/review.md names the latest code commit")
		return r.Finish(1)
	}
	data, err := os.ReadFile(filepath.Join(root, file))
	if err != nil {
		return r.Abort(file, "unreadable", err.Error())
	}
	round := b.Touching(file)
	r.Line(fmt.Sprintf("reviewed %s in %s, round %d", h, file, round))

	res := Parse(string(data))
	switch {
	case !res.Valid:
		r.Error(file, res.Line, "bad-result", "the Result line is not: Result: <b> blocking, <n> notes")
	case res.Blocking != res.FoundBlocking || res.Notes != res.FoundNotes:
		r.Error(file, res.Line, "result-mismatch", fmt.Sprintf("the Result line says %d blocking, %d notes; the findings hold %d blocking, %d notes",
			res.Blocking, res.Notes, res.FoundBlocking, res.FoundNotes))
	}
	blocking := res.Blocking
	if !res.Valid || res.FoundBlocking > blocking {
		blocking = res.FoundBlocking
	}
	if blocking > 0 {
		detail := fmt.Sprintf("%d in round %d", blocking, round)
		if round >= Rounds {
			detail += fmt.Sprintf("; after %d rounds the engineer decides", Rounds)
		}
		r.Error(file, 0, "blocking", detail)
	}
	return r.Finish(1)
}
