// Package prcheck checks the paths that a pull request changes: whether it
// mixes path groups, and whether it needs the engineer's approval.
package prcheck

import (
	"sort"
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
)

// Result is what Check found out about a pull request.
type Result struct {
	// Changed maps each path that the pull request changes to its group.
	Changed map[string]paths.Group
	// Mixed is true when a path in the design-notes or CI-and-config group
	// comes with a path outside that group.
	Mixed bool
	// ApprovalRequired is true unless the commits keep the order: first
	// only new tests, then only approval-free paths.
	ApprovalRequired bool
	// TestsEnd is the last commit of the test phase, or "".
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
		return r.Abort(dir, "shallow-clone", "fetch the full history, for example with fetch-depth: 0")
	}
	if _, err := repo.Run("rev-parse", "--verify", "--quiet", base+"^{commit}"); err != nil {
		return r.Abort(base, "unknown-ref", "no commit "+base)
	}

	res, err := Check(repo, base)
	if err != nil {
		return r.Abort(dir, "git-failed", err.Error())
	}

	if res.Mixed {
		var changed []string
		for p := range res.Changed {
			changed = append(changed, p)
		}
		sort.Strings(changed)
		for _, p := range changed {
			r.Error(p, 0, "mixed-pr", "group "+string(res.Changed[p]))
		}
	}
	if res.ApprovalRequired {
		r.Line("approval required")
	} else {
		r.Line("approval not required")
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
	res.Changed = map[string]paths.Group{}
	kinds := map[paths.Group]bool{}
	for _, p := range overall {
		g := paths.GroupOf(p)
		res.Changed[p] = g
		kinds[g] = true
	}
	res.Mixed = len(kinds) > 1 && (kinds[paths.DesignNotes] || kinds[paths.CIConfig])

	commits, err := readCommits(repo, mb, res.Changed)
	if err != nil {
		return res, err
	}
	list, hasList, err := readList(repo, base)
	if err != nil {
		return res, err
	}
	onBase, err := testFilesOn(repo, base)
	if err != nil {
		return res, err
	}

	res.ApprovalRequired = res.Mixed || !hasList
	testPhase := kinds[paths.Tests]
	for _, c := range commits {
		if testPhase && allInGroup(c.paths, paths.Tests) {
			for _, p := range c.paths {
				if strings.HasPrefix(p, paths.TestDir) && onBase[p] {
					res.ApprovalRequired = true
				}
			}
			res.TestsEnd = c.hash
			continue
		}
		testPhase = false
		for _, p := range c.paths {
			if !list.ApprovalFree(p) {
				res.ApprovalRequired = true
			}
		}
	}
	return res, nil
}

// readCommits returns the commits on the first-parent line from mb to HEAD,
// oldest first, each with the paths that it changes against its first
// parent. A merge commit counts only the paths that the pull request
// changes.
func readCommits(repo gitx.Repo, mb string, changed map[string]paths.Group) ([]commit, error) {
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
				if _, ok := changed[p]; ok {
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

// readList reads the list of approval-free paths from base.
func readList(repo gitx.Repo, base string) (paths.List, bool, error) {
	if _, err := repo.Run("cat-file", "-e", base+":"+paths.ApprovalFreeFile); err != nil {
		return nil, false, nil
	}
	text, err := repo.Run("show", base+":"+paths.ApprovalFreeFile)
	if err != nil {
		return nil, false, err
	}
	return paths.ParseList(text), true, nil
}

// testFilesOn returns the test files that exist on base.
func testFilesOn(repo gitx.Repo, base string) (map[string]bool, error) {
	files, err := repo.Paths("ls-tree", "-r", "-z", "--name-only", base, "--", strings.TrimSuffix(paths.TestDir, "/"))
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for _, f := range files {
		set[f] = true
	}
	return set, nil
}

func allInGroup(ps []string, g paths.Group) bool {
	for _, p := range ps {
		if paths.GroupOf(p) != g {
			return false
		}
	}
	return true
}
