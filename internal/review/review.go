// Package review finds the latest code commit that still needs a review
// result.
package review

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
)

// Bases are the refs, in this order, whose merge base with HEAD starts the
// commits of the branch.
var Bases = []string{"origin/HEAD", "origin/main", "main"}

// ReadCode reads .factory/code-paths from the working tree at root.
func ReadCode(root string) (paths.Code, error) {
	data, err := os.ReadFile(filepath.Join(root, paths.CodePathsFile))
	switch {
	case err == nil:
		return paths.Code{List: paths.ParseList(string(data)), Set: true}, nil
	case os.IsNotExist(err):
		return paths.Code{}, nil
	}
	return paths.Code{}, err
}

// Pending returns the latest commit of the branch at root that changes
// code, when no review result in evidence/*/review.md names its full hash.
// It returns "" when there is no such commit.
func Pending(root string) (string, error) {
	repo := gitx.Repo{Dir: root}
	code, err := ReadCode(root)
	if err != nil {
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
		if !changesCode(ps, code) {
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

func changesCode(ps []string, code paths.Code) bool {
	for _, p := range ps {
		if code.Is(p) {
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
