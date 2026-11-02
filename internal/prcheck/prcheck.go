// Package prcheck checks the commits of a pull request: whether it mixes
// path groups, and where its tests end.
package prcheck

import (
	"sort"
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
)

// Change is one path that a pull request changes, with its group.
type Change struct {
	Path  string
	Group paths.Group
}

// Mixed is a path that breaks the rules of the path groups.
type Mixed struct {
	Path   string
	Detail string
}

// Result is what Check found out about a pull request.
type Result struct {
	// Changed are the paths that the pull request changes, in order.
	Changed []Change
	// Mixed are the paths that break the rules of the path groups.
	Mixed []Mixed
	// TestsEnd is the last commit that changes test files, or "".
	TestsEnd string
}

// commit is one commit of the pull request with the paths it changes.
type commit struct {
	hash  string
	paths []string
}

// Run checks HEAD of the repo at dir against base and returns the exit code.
func Run(r *report.Report, dir, base string) int {
	repo := gitx.Repo{Dir: dir}
	shallow, err := repo.Run("rev-parse", "--is-shallow-repository")
	if err != nil {
		return r.Abort(dir, "not-a-repo", err.Error())
	}
	if strings.TrimSpace(shallow) == "true" {
		return r.Abort(dir, "shallow-clone", "the clone lacks history; fetch it in full, for example with fetch-depth: 0")
	}
	if _, err := repo.Run("rev-parse", "--verify", "--quiet", base+"^{commit}"); err != nil {
		return r.Abort(base, "unknown-ref", "no commit "+base)
	}

	res, err := Check(repo, base)
	if err != nil {
		return r.Abort(dir, "git-failed", err.Error())
	}
	for _, m := range res.Mixed {
		r.Error(m.Path, 0, "mixed-pr", m.Detail)
	}
	if res.TestsEnd != "" {
		r.Line("tests end at " + res.TestsEnd)
	}
	return r.Finish(len(res.Changed))
}

// Check reads the commits between the merge base with base and HEAD.
func Check(repo gitx.Repo, base string) (Result, error) {
	var res Result
	mb, err := repo.Run("merge-base", base, "HEAD")
	if err != nil {
		return res, err
	}
	mb = strings.TrimSpace(mb)

	overall, err := repo.Paths("diff", "--name-only", "--no-renames", "-z", mb, "HEAD")
	if err != nil {
		return res, err
	}
	sort.Strings(overall)
	changed := map[string]bool{}
	kinds := map[paths.Group]bool{}
	for _, p := range overall {
		g := paths.GroupOf(p)
		res.Changed = append(res.Changed, Change{Path: p, Group: g})
		changed[p] = true
		if !paths.IsEvidence(p) {
			kinds[g] = true
		}
	}

	commits, err := readCommits(repo, mb, changed)
	if err != nil {
		return res, err
	}

	switch {
	case len(kinds) <= 1:
	case len(kinds) == 2 && kinds[paths.Tests] && kinds[paths.Other]:
		// Test files first, then only paths outside the groups.
	default:
		for _, c := range res.Changed {
			if !paths.IsEvidence(c.Path) {
				res.Mixed = append(res.Mixed, Mixed{Path: c.Path, Detail: "group " + string(c.Group) + "; a pull request changes one group, or first test files and then only paths outside the groups"})
			}
		}
		return res, nil
	}

	testPhase := true
	for _, c := range commits {
		groups := map[paths.Group]bool{}
		for _, p := range c.paths {
			if !paths.IsEvidence(p) {
				groups[paths.GroupOf(p)] = true
			}
		}
		if len(groups) == 0 {
			continue
		}
		if testPhase && len(groups) == 1 && groups[paths.Tests] {
			res.TestsEnd = c.hash
			continue
		}
		testPhase = false
		if !groups[paths.Tests] {
			continue
		}
		for _, p := range c.paths {
			if paths.GroupOf(p) == paths.Tests {
				res.Mixed = append(res.Mixed, Mixed{Path: p, Detail: "group tests, changed in commit " + c.hash[:12] + " after the first commit outside the tests group"})
			}
		}
	}
	return res, nil
}

// readCommits returns the commits on the first-parent line from mb to HEAD,
// oldest first, each with the paths that it changes against its first
// parent. A merge commit counts only the paths that the pull request
// changes.
func readCommits(repo gitx.Repo, mb string, changed map[string]bool) ([]commit, error) {
	lines, err := repo.Lines("rev-list", "--reverse", "--first-parent", "--parents", mb+"..HEAD")
	if err != nil {
		return nil, err
	}
	var commits []commit
	for _, line := range lines {
		fields := strings.Fields(line)
		hash := fields[0]
		var ps []string
		if len(fields) == 1 {
			ps, err = repo.Paths("diff-tree", "--root", "-r", "--name-only", "--no-renames", "--no-commit-id", "-z", hash)
		} else {
			ps, err = repo.Paths("diff", "--name-only", "--no-renames", "-z", fields[1], hash)
		}
		if err != nil {
			return nil, err
		}
		if len(fields) > 2 {
			var kept []string
			for _, p := range ps {
				if changed[p] {
					kept = append(kept, p)
				}
			}
			ps = kept
		}
		if len(ps) > 0 {
			commits = append(commits, commit{hash: hash, paths: ps})
		}
	}
	return commits, nil
}
