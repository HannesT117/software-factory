// Package skills checks skills against the fields of the Agent Skills
// specification and against names of specific agent tools.
package skills

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/HannesT117/software-factory/internal/report"
)

var (
	nameFormat = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	// toolName finds a specific agent tool, in any case.
	toolName = regexp.MustCompile(`(?i)\b(claude|cursor|codex|copilot|gemini)\b`)
)

var allowed = map[string]bool{
	"name": true, "description": true, "license": true,
	"compatibility": true, "metadata": true, "allowed-tools": true,
}

// Validate checks every SKILL.md file under dirs and returns the exit code.
func Validate(r *report.Report, dirs []string) int {
	var files []string
	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && d.Name() == ".git" {
				return filepath.SkipDir
			}
			if !d.IsDir() && d.Name() == "SKILL.md" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return r.Abort(dir, "unreadable", err.Error())
		}
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return r.Abort(f, "unreadable", err.Error())
		}
		check(r, filepath.ToSlash(f), filepath.Base(filepath.Dir(f)), string(data))
	}
	return r.Finish(len(files))
}

// check checks one skill file. folder is the name of its folder.
func check(r *report.Report, path, folder, text string) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		r.Error(path, 0, "no-frontmatter", "the file does not start with ---")
		return
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		r.Error(path, 0, "no-frontmatter", "no closing ---")
		return
	}
	fields, err := parseFrontmatter(lines[1:end], 2)
	if err != nil {
		var ye *yamlError
		if errors.As(err, &ye) {
			r.Error(path, ye.line, "unsupported-yaml", ye.detail)
		} else {
			r.Error(path, 0, "unsupported-yaml", err.Error())
		}
		return
	}

	byKey := map[string]field{}
	for _, f := range fields {
		byKey[f.key] = f
		if !allowed[f.key] {
			r.Error(path, f.line, "unknown-field", f.key)
			continue
		}
		if f.isMap && f.key != "metadata" {
			r.Error(path, f.line, "unsupported-yaml", "only metadata may be a map")
		}
	}

	name, ok := byKey["name"]
	switch {
	case !ok || name.value == "":
		r.Error(path, name.line, "name-format", "name is missing")
	case utf8.RuneCountInString(name.value) > 64:
		r.Error(path, name.line, "name-too-long", "more than 64 characters")
	case !nameFormat.MatchString(name.value):
		r.Error(path, name.line, "name-format", "use lowercase letters, digits and single hyphens")
	}
	if ok && name.value != "" && name.value != folder {
		r.Error(path, name.line, "name-mismatch", "name "+name.value+" differs from folder "+folder)
	}

	desc, ok := byKey["description"]
	switch {
	case !ok || strings.TrimSpace(desc.value) == "":
		r.Error(path, desc.line, "description-missing", "")
	case utf8.RuneCountInString(desc.value) > 1024:
		r.Error(path, desc.line, "description-too-long", "more than 1024 characters")
	}

	if c, ok := byKey["compatibility"]; ok && utf8.RuneCountInString(c.value) > 500 {
		r.Error(path, c.line, "compatibility-too-long", "more than 500 characters")
	}
	if m, ok := byKey["metadata"]; ok && !m.isMap {
		r.Error(path, m.line, "metadata-not-map", "")
	}

	for _, key := range []string{"name", "description"} {
		if f, ok := byKey[key]; ok {
			if m := toolName.FindString(f.value); m != "" {
				r.Error(path, f.line, "tool-name", key+" names "+m)
			}
		}
	}
	for i := end + 1; i < len(lines); i++ {
		if m := toolName.FindString(lines[i]); m != "" {
			r.Error(path, i+1, "tool-name", "body names "+m)
		}
	}
}
