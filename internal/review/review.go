// Package review finds the latest code commit that still needs a review
// result.
package review

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
)

// resultFile is the review result of an increment.
var resultFile = regexp.MustCompile(`^evidence/[^/]+/review\.md$`)

// Bases are the refs, in this order, whose merge base with HEAD starts the
// commits of the branch.
var Bases = []string{"origin/HEAD", "origin/main", "main"}

// Pending returns the latest commit of the branch at root that changes an
// approval-free path, when no review result in evidence/*/review.md names
// its full hash. It returns "" when there is no such commit. Without the
// list of approval-free paths, every path outside the groups counts.
// Changes to review results themselves do not count.
func Pending(root string) (string, error) {
	repo := gitx.Repo{Dir: root}
	var list paths.List
	data, err := os.ReadFile(filepath.Join(root, paths.ApprovalFreeFile))
	switch {
	case err == nil:
		list = paths.ParseList(string(data))
	case os.IsNotExist(err):
		list = paths.List{"**"}
	default:
		return "", err
	}

	rangeArg := "HEAD"
	for _, b := range Bases {
		mb, err := repo.Run("merge-base", b, "HEAD")
		if err == nil {
			rangeArg = strings.TrimSpace(mb) + "..HEAD"
			break
		}
	}
	hashes, err := repo.Lines("rev-list", "--first-parent", rangeArg)
	if err != nil {
		return "", err
	}

	for _, h := range hashes {
		ps, err := repo.Paths("diff-tree", "--root", "-r", "-m", "--first-parent", "--name-only", "--no-renames", "--no-commit-id", "-z", h)
		if err != nil {
			return "", err
		}
		if !changesCode(ps, list) {
			continue
		}
		reviewed, err := named(root, h)
		if err != nil || reviewed {
			return "", err
		}
		return h, nil
	}
	return "", nil
}

func changesCode(ps []string, list paths.List) bool {
	for _, p := range ps {
		if !resultFile.MatchString(p) && list.ApprovalFree(p) {
			return true
		}
	}
	return false
}

// named reports whether a review result in the working tree names hash.
func named(root, hash string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(root, "evidence", "*", "review.md"))
	if err != nil {
		return false, err
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return false, err
		}
		if strings.Contains(string(data), hash) {
			return true, nil
		}
	}
	return false, nil
}
