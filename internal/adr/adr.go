// Package adr checks decision records against the MADR subset of the
// factory.
package adr

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/report"
)

var categories = map[string]bool{"architecture": true, "process": true}

var (
	superseded = regexp.MustCompile(`^superseded by ADR-(\d{4})$`)
	recordNum  = regexp.MustCompile(`^(\d{4})-`)
)

var (
	required = []string{"Context and Problem Statement", "Considered Options", "Decision Outcome"}
	optional = []string{"Decision Drivers", "Pros and Cons of the Options", "More Information"}
	outcome  = map[string]bool{"Consequences": true, "Confirmation": true}
)

// Check checks every .md file under dir and returns the exit code.
func Check(r *report.Report, dir string) int {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return r.Abort(dir, "unreadable", err.Error())
	}
	sort.Strings(files)

	// numbers holds the record numbers per category.
	numbers := map[string]map[string]bool{}
	for _, f := range files {
		cat := category(dir, f)
		if m := recordNum.FindStringSubmatch(filepath.Base(f)); m != nil {
			if numbers[cat] == nil {
				numbers[cat] = map[string]bool{}
			}
			numbers[cat][m[1]] = true
		}
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return r.Abort(f, "unreadable", err.Error())
		}
		where := filepath.ToSlash(f)
		cat := category(dir, f)
		if !categories[cat] {
			r.Error(where, 0, "unknown-category", "the folder is not architecture or process")
		}
		checkRecord(r, where, string(data), numbers[cat])
	}
	return r.Finish(len(files))
}

// category returns the first folder of file below dir, or "" for a file
// directly in dir.
func category(dir, file string) string {
	rel, err := filepath.Rel(dir, file)
	if err != nil {
		return ""
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

func checkRecord(r *report.Report, where, text string, numbers map[string]bool) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	body := 0
	if strings.TrimSpace(lines[0]) == "---" {
		end := -1
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				end = i
				break
			}
		}
		if end < 0 {
			r.Error(where, 1, "unknown-field", "the frontmatter has no closing ---")
			return
		}
		checkFrontmatter(r, where, lines[1:end], numbers)
		body = end + 1
	}

	rest := strings.Join(lines[body:], "\n")
	offset := body
	sawTitle, sawHeading := false, false
	sections := map[string]bool{}
	section := ""
	design.Lines(rest, func(n int, line string) {
		n += offset
		level, name := heading(line)
		if level == 0 {
			return
		}
		if !sawHeading {
			sawHeading = true
			sawTitle = level == 1
		}
		switch level {
		case 2:
			section = name
			sections[name] = true
			if !contains(required, name) && !contains(optional, name) {
				r.Error(where, n, "unknown-section", "## "+name)
			}
		case 3:
			if section == "Decision Outcome" && !outcome[name] {
				r.Error(where, n, "unknown-section", "### "+name+" under Decision Outcome")
			}
		}
	})
	if !sawTitle {
		r.Error(where, 0, "missing-title", "the first heading is not # <title>")
	}
	for _, s := range required {
		if !sections[s] {
			r.Error(where, 0, "missing-section", "## "+s)
		}
	}
}

func checkFrontmatter(r *report.Report, where string, lines []string, numbers map[string]bool) {
	for i, line := range lines {
		n := i + 2
		if strings.TrimSpace(line) == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		switch {
		case !ok:
			r.Error(where, n, "unknown-field", "not a key: value line")
		case key == "date":
			if _, err := time.Parse("2006-01-02", value); err != nil {
				r.Error(where, n, "bad-date", "use YYYY-MM-DD")
			}
		case key == "status":
			if value == "deprecated" {
				continue
			}
			m := superseded.FindStringSubmatch(value)
			if m == nil {
				r.Error(where, n, "bad-status", "use superseded by ADR-<nnnn> or deprecated")
			} else if !numbers[m[1]] {
				r.Error(where, n, "unknown-adr", "no record "+m[1]+" in this category")
			}
		default:
			r.Error(where, n, "unknown-field", key)
		}
	}
}

// heading returns the level and text of a Markdown heading line, or 0.
func heading(line string) (int, string) {
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > 6 || level >= len(line) || line[level] != ' ' {
		return 0, ""
	}
	return level, strings.TrimSpace(line[level:])
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
