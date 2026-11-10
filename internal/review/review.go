// Package review finds the latest code commit of a branch, the review
// result that names it, and reads that result.
package review

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/HannesT117/software-factory/internal/branch"
	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
)

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
	code, err := ReadCode(root)
	if err != nil {
		return "", err
	}
	b, err := branch.ReadAny(gitx.Repo{Dir: root})
	if err != nil {
		return "", err
	}
	h := b.LatestCode(code)
	if h == "" {
		return "", nil
	}
	file, err := Find(root, h)
	if err != nil || file != "" {
		return "", err
	}
	return h, nil
}

// Find returns the path, relative to root and slash-separated, of the
// first review result in the working tree that names hash, or "".
func Find(root, hash string) (string, error) {
	files, err := FindAll(root, hash)
	if err != nil || len(files) == 0 {
		return "", err
	}
	return files[0], nil
}

// FindAll returns the paths, relative to root and slash-separated, of
// every review result in the working tree that names hash, in path order.
func FindAll(root, hash string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(root, "evidence", "*", "review.md"))
	if err != nil {
		return nil, err
	}
	var found []string
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if strings.Contains(string(data), hash) {
			rel, err := filepath.Rel(root, f)
			if err != nil {
				return nil, err
			}
			found = append(found, filepath.ToSlash(rel))
		}
	}
	return found, nil
}

// ResultLine is the form of the Result line of a review result.
var ResultLine = regexp.MustCompile(`^Result: (\d+) blocking, (\d+) notes$`)

// Result is what a review result says about its findings.
type Result struct {
	// Line is the number of the Result line, or 0 without one.
	Line int
	// Valid is true when the Result line has the form of ResultLine.
	Valid bool
	// Blocking and Notes are the counts of the Result line.
	Blocking, Notes int
	// FoundBlocking and FoundNotes count the finding lines under
	// ## Findings.
	FoundBlocking, FoundNotes int
}

// Parse reads a review result.
func Parse(text string) Result {
	var res Result
	inFindings := false
	design.Lines(text, func(n int, line string) {
		line = strings.TrimRight(line, " \r")
		if strings.HasPrefix(line, "#") {
			inFindings = line == "## Findings"
			return
		}
		if strings.HasPrefix(line, "Result:") && res.Line == 0 {
			res.Line = n
			if m := ResultLine.FindStringSubmatch(line); m != nil {
				res.Valid = true
				res.Blocking, _ = strconv.Atoi(m[1])
				res.Notes, _ = strconv.Atoi(m[2])
			}
			return
		}
		if !inFindings {
			return
		}
		switch {
		case strings.HasPrefix(line, "- **blocking**"):
			res.FoundBlocking++
		case strings.HasPrefix(line, "- **note**"):
			res.FoundNotes++
		}
	})
	return res
}

// ReadCodeAt reads .factory/code-paths from the commit ref of repo.
func ReadCodeAt(repo gitx.Repo, ref string) (paths.Code, error) {
	if _, err := repo.Run("cat-file", "-e", ref+":"+paths.CodePathsFile); err != nil {
		return paths.Code{}, nil
	}
	text, err := repo.Run("show", ref+":"+paths.CodePathsFile)
	if err != nil {
		return paths.Code{}, err
	}
	return paths.Code{List: paths.ParseList(text), Set: true}, nil
}
