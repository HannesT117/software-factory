// Package branch reads the commits of a branch against its base, as a pull
// request shows them.
package branch

import (
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
)

// Bases are the refs, in this order, whose merge base with HEAD starts the
// commits of the branch when no base is given.
var Bases = []string{"origin/HEAD", "origin/main", "main"}

// Commit is one commit of the branch with the paths it changes.
type Commit struct {
	Hash  string
	Paths []string
}

// Branch is HEAD against the merge base with a base ref.
type Branch struct {
	// MergeBase is the merge base, or "" when no base was found.
	MergeBase string
	// Changed are the paths that HEAD changes against the merge base.
	Changed []string
	// Commits are the commits on the first-parent line from the merge base
	// to HEAD, oldest first, each with the paths that it changes against its
	// first parent. A merge commit counts only with the paths in Changed,
	// so that a merge from the base brings in nothing. Commits without such
	// paths are left out.
	Commits []Commit
}

// Read reads HEAD of repo against base.
func Read(repo gitx.Repo, base string) (Branch, error) {
	var b Branch
	mb, err := repo.Run("merge-base", base, "HEAD")
	if err != nil {
		return b, err
	}
	return read(repo, strings.TrimSpace(mb))
}

// ReadAny reads HEAD against the first of Bases that has a merge base with
// HEAD. Without one, every commit of HEAD counts, and Changed is nil.
func ReadAny(repo gitx.Repo) (Branch, error) {
	for _, base := range Bases {
		mb, err := repo.Run("merge-base", base, "HEAD")
		if err == nil {
			return read(repo, strings.TrimSpace(mb))
		}
	}
	return read(repo, "")
}

func read(repo gitx.Repo, mb string) (Branch, error) {
	b := Branch{MergeBase: mb}
	rangeArg := "HEAD"
	var changed map[string]bool
	if mb != "" {
		rangeArg = mb + "..HEAD"
		ps, err := repo.Paths("diff", "--name-only", "--no-renames", "-z", mb, "HEAD")
		if err != nil {
			return b, err
		}
		b.Changed = ps
		changed = map[string]bool{}
		for _, p := range ps {
			changed[p] = true
		}
	}
	lines, err := repo.Lines("rev-list", "--reverse", "--first-parent", "--parents", rangeArg)
	if err != nil {
		return b, err
	}
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
			return b, err
		}
		if len(fields) > 2 && changed != nil {
			var kept []string
			for _, p := range ps {
				if changed[p] {
					kept = append(kept, p)
				}
			}
			ps = kept
		}
		if len(ps) > 0 {
			b.Commits = append(b.Commits, Commit{Hash: hash, Paths: ps})
		}
	}
	return b, nil
}

// LatestCode returns the latest commit that changes a path that counts as
// code, or "" when there is none.
func (b Branch) LatestCode(code paths.Code) string {
	for i := len(b.Commits) - 1; i >= 0; i-- {
		for _, p := range b.Commits[i].Paths {
			if code.Is(p) {
				return b.Commits[i].Hash
			}
		}
	}
	return ""
}

// Touching returns the number of commits that change path.
func (b Branch) Touching(path string) int {
	n := 0
	for _, c := range b.Commits {
		for _, p := range c.Paths {
			if p == path {
				n++
				break
			}
		}
	}
	return n
}
