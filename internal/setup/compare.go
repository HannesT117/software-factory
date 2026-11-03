package setup

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// lookup follows keys through nested JSON objects. It returns nil when a
// key is missing.
func lookup(v any, keys ...string) any {
	for _, k := range keys {
		m, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		v = m[k]
	}
	return v
}

// equal reports whether two JSON values are equal. Numbers compare by
// their text, as json.Number keeps it.
func equal(a, b any) bool {
	return encode(a) == encode(b)
}

// contains reports whether the current value holds the target value: an
// object holds every key of the target with a value that holds the
// target's value, and may have more keys, which GitHub adds as defaults;
// an array has the same length and each target element is held by its
// own element of the current array, in any order; any other value is
// equal.
func contains(current, target any) bool {
	switch t := target.(type) {
	case map[string]any:
		c, ok := current.(map[string]any)
		if !ok {
			return false
		}
		for k, tv := range t {
			cv, ok := c[k]
			if !ok || !contains(cv, tv) {
				return false
			}
		}
		return true
	case []any:
		c, ok := current.([]any)
		if !ok || len(c) != len(t) {
			return false
		}
		used := make([]bool, len(c))
		for _, tv := range t {
			found := false
			for i, cv := range c {
				if !used[i] && contains(cv, tv) {
					used[i], found = true, true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	return equal(current, target)
}

// diff formats "<field> <current> -> <target>". A missing current value
// is "missing".
func diff(field string, current any, present bool, target any) string {
	cur := "missing"
	if present {
		cur = show(current)
	}
	return fmt.Sprintf("%s %s -> %s", field, cur, show(target))
}

// show prints a string as it is and any other value as compact JSON.
func show(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return encode(v)
}

func encode(v any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return fmt.Sprint(v)
	}
	return string(bytes.TrimSpace(buf.Bytes()))
}
