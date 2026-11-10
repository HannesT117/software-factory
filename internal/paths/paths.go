// Package paths knows the path groups of a project and matches paths
// against patterns with * and **.
package paths

import (
	"bufio"
	"strings"
)

// Group is a set of paths that a pull request changes only on its own.
type Group string

// The path groups. Other is every path outside a group.
const (
	Tests       Group = "tests"
	DesignNotes Group = "design-notes"
	CIConfig    Group = "ci-and-config"
	Other       Group = "other"
)

var groups = []struct {
	group    Group
	patterns []string
}{
	{Tests, []string{"test/acceptance/**", "go.mod", "go.sum"}},
	{DesignNotes, []string{"docs/increments/*/design.md", "docs/domain/**"}},
	{CIConfig, []string{".github/**", ".factory/**"}},
}

// TestDir is the folder of the acceptance tests.
const TestDir = "test/acceptance/"

// CodePathsFile lists the paths that count as code.
const CodePathsFile = ".factory/code-paths"

// EvidenceDir holds review results and recordings. A path under it goes
// with any group and never counts as code.
const EvidenceDir = "evidence/"

// GroupOf returns the group of a slash-separated path relative to the
// repository root.
func GroupOf(path string) Group {
	for _, g := range groups {
		for _, p := range g.patterns {
			if Match(p, path) {
				return g.group
			}
		}
	}
	return Other
}

// Match reports whether path matches pattern. In a pattern, * matches any
// part of one path segment, and a segment ** matches zero or more segments.
func Match(pattern, path string) bool {
	return matchSegments(strings.Split(pattern, "/"), strings.Split(path, "/"))
}

func matchSegments(pattern, path []string) bool {
	if len(pattern) == 0 {
		return len(path) == 0
	}
	if pattern[0] == "**" {
		for i := 0; i <= len(path); i++ {
			if matchSegments(pattern[1:], path[i:]) {
				return true
			}
		}
		return false
	}
	if len(path) == 0 || !matchSegment(pattern[0], path[0]) {
		return false
	}
	return matchSegments(pattern[1:], path[1:])
}

// matchSegment matches one segment, where * matches any run of characters.
func matchSegment(pattern, s string) bool {
	star := strings.IndexByte(pattern, '*')
	if star < 0 {
		return pattern == s
	}
	prefix := pattern[:star]
	if !strings.HasPrefix(s, prefix) {
		return false
	}
	rest := pattern[star+1:]
	for i := len(prefix); i <= len(s); i++ {
		if matchSegment(rest, s[i:]) {
			return true
		}
	}
	return false
}

// List is a list of path patterns, as in .factory/code-paths.
type List []string

// ParseList reads one pattern per line. Empty lines and lines that start
// with # are skipped.
func ParseList(text string) List {
	var list List
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		list = append(list, line)
	}
	return list
}

// Matches reports whether path matches a pattern of the list.
func (l List) Matches(path string) bool {
	for _, p := range l {
		if Match(p, path) {
			return true
		}
	}
	return false
}

// IsEvidence reports whether path lies under evidence/.
func IsEvidence(path string) bool {
	return strings.HasPrefix(path, EvidenceDir)
}

// Code tells which paths count as code. A path in the tests group counts,
// and a path that matches the list of .factory/code-paths. Without the
// list, every path outside the groups counts. A path under evidence/ never
// counts.
type Code struct {
	List List
	Set  bool // whether .factory/code-paths exists
}

// Is reports whether a change to path counts as code.
func (c Code) Is(path string) bool {
	if IsEvidence(path) {
		return false
	}
	switch GroupOf(path) {
	case Tests:
		return true
	case Other:
		return !c.Set || c.List.Matches(path)
	}
	return false
}

// IncrementOf returns the folder <nn>-<slug> of an increment when path lies
// under test/acceptance/<nn>-<slug>/ or evidence/<nn>-<slug>/. A pull
// request belongs to the increments of the paths that it changes.
func IncrementOf(path string) (string, bool) {
	for _, dir := range []string{TestDir, EvidenceDir} {
		rest, ok := strings.CutPrefix(path, dir)
		if !ok {
			continue
		}
		folder, _, ok := strings.Cut(rest, "/")
		if ok && len(folder) > 3 && folder[0] >= '0' && folder[0] <= '9' && folder[1] >= '0' && folder[1] <= '9' && folder[2] == '-' {
			return folder, true
		}
	}
	return "", false
}

// Increments returns the folders of the increments that the paths belong
// to.
func Increments(ps []string) map[string]bool {
	set := map[string]bool{}
	for _, p := range ps {
		if f, ok := IncrementOf(p); ok {
			set[f] = true
		}
	}
	return set
}
