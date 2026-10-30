package skills

import (
	"fmt"
	"strings"
)

// field is one top-level frontmatter field.
type field struct {
	key  string
	line int // line in the file, starting at 1
	// value is the string value. isMap is true for a nested map, which is
	// then in entries.
	value   string
	isMap   bool
	entries map[string]string
}

// yamlError is a frontmatter line outside the supported subset.
type yamlError struct {
	line   int
	detail string
}

func (e *yamlError) Error() string { return e.detail }

// parseFrontmatter reads the supported subset of YAML: key: value, quoted
// values, | and > blocks, and nested maps of plain or quoted values.
// first is the file line of lines[0].
func parseFrontmatter(lines []string, first int) ([]field, error) {
	var fields []field
	seen := map[string]bool{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		n := first + i
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			return nil, &yamlError{n, "unexpected indentation"}
		}
		key, rest, ok := strings.Cut(line, ":")
		if !ok || !isKey(key) {
			return nil, &yamlError{n, "not a key: value line"}
		}
		if seen[key] {
			return nil, &yamlError{n, "duplicate key " + key}
		}
		seen[key] = true
		if rest != "" && rest[0] != ' ' {
			return nil, &yamlError{n, "no space after the colon"}
		}
		rest = strings.TrimSpace(rest)
		f := field{key: key, line: n}
		block := indentedBlock(lines[i+1:])
		switch {
		case rest == "" && len(block) > 0:
			entries, err := parseMap(block, n+1)
			if err != nil {
				return nil, err
			}
			f.isMap, f.entries = true, entries
		case isBlockIndicator(rest):
			f.value = foldBlock(block, rest[0] == '>', rest)
		default:
			if len(block) > 0 {
				return nil, &yamlError{n + 1, "unexpected indentation"}
			}
			if rest == "" {
				break
			}
			v, err := scalar(rest)
			if err != nil {
				return nil, &yamlError{n, err.Error()}
			}
			f.value = v
		}
		i += len(block)
		fields = append(fields, f)
	}
	return fields, nil
}

// indentedBlock returns the lines that belong to the field before them: the
// indented lines and the empty lines between them.
func indentedBlock(lines []string) []string {
	end := 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0] != ' ' && line[0] != '\t' {
			break
		}
		end = i + 1
	}
	return lines[:end]
}

func parseMap(lines []string, first int) (map[string]string, error) {
	entries := map[string]string{}
	indent := -1
	for i, line := range lines {
		n := first + i
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(trimmed, "\t") {
			return nil, &yamlError{n, "tab in indentation"}
		}
		if indent < 0 {
			indent = len(line) - len(trimmed)
		}
		if len(line)-len(trimmed) != indent {
			return nil, &yamlError{n, "nested deeper than one map"}
		}
		key, rest, ok := strings.Cut(trimmed, ":")
		if !ok || !isKey(key) || (rest != "" && rest[0] != ' ') {
			return nil, &yamlError{n, "not a key: value line"}
		}
		if _, dup := entries[key]; dup {
			return nil, &yamlError{n, "duplicate key " + key}
		}
		rest = strings.TrimSpace(rest)
		if rest == "" || isBlockIndicator(rest) {
			return nil, &yamlError{n, "nested value is not a string"}
		}
		v, err := scalar(rest)
		if err != nil {
			return nil, &yamlError{n, err.Error()}
		}
		entries[key] = v
	}
	return entries, nil
}

func isKey(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

func isBlockIndicator(s string) bool {
	switch s {
	case "|", "|-", "|+", ">", ">-", ">+":
		return true
	}
	return false
}

// foldBlock returns the text of a | or > block. A > block joins the lines of
// a paragraph with spaces.
func foldBlock(lines []string, folded bool, indicator string) string {
	indent := -1
	var body []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			body = append(body, "")
			continue
		}
		trimmed := strings.TrimLeft(line, " \t")
		if indent < 0 {
			indent = len(line) - len(trimmed)
		}
		if len(line) >= indent {
			body = append(body, line[indent:])
		} else {
			body = append(body, trimmed)
		}
	}
	for len(body) > 0 && body[len(body)-1] == "" {
		body = body[:len(body)-1]
	}
	var text string
	if folded {
		var paras []string
		var cur []string
		for _, l := range body {
			if l == "" {
				paras = append(paras, strings.Join(cur, " "))
				cur = nil
				continue
			}
			cur = append(cur, l)
		}
		paras = append(paras, strings.Join(cur, " "))
		text = strings.Join(paras, "\n")
	} else {
		text = strings.Join(body, "\n")
	}
	if !strings.HasSuffix(indicator, "-") && text != "" {
		text += "\n"
	}
	return text
}

// scalar reads a plain, single-quoted or double-quoted value.
func scalar(s string) (string, error) {
	switch s[0] {
	case '\'':
		if len(s) < 2 || s[len(s)-1] != '\'' {
			return "", fmt.Errorf("unclosed quote")
		}
		inner := s[1 : len(s)-1]
		if strings.Contains(strings.ReplaceAll(inner, "''", ""), "'") {
			return "", fmt.Errorf("unescaped quote")
		}
		return strings.ReplaceAll(inner, "''", "'"), nil
	case '"':
		if len(s) < 2 || s[len(s)-1] != '"' {
			return "", fmt.Errorf("unclosed quote")
		}
		return unescapeDouble(s[1 : len(s)-1])
	case '[', '{', '&', '*', '!', '|', '>', '%', '@', '`':
		return "", fmt.Errorf("value starts with %q", s[0])
	}
	if s == "-" || strings.HasPrefix(s, "- ") {
		return "", fmt.Errorf("list item")
	}
	if strings.Contains(s, ": ") {
		return "", fmt.Errorf("colon and space in a plain value")
	}
	if i := strings.Index(s, " #"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return s, nil
}

func unescapeDouble(s string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			return "", fmt.Errorf("unescaped quote")
		}
		if c != '\\' {
			b.WriteByte(c)
			continue
		}
		i++
		if i == len(s) {
			return "", fmt.Errorf("backslash at the end")
		}
		switch s[i] {
		case '\\', '"', '/':
			b.WriteByte(s[i])
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		default:
			return "", fmt.Errorf("unsupported escape \\%c", s[i])
		}
	}
	return b.String(), nil
}
