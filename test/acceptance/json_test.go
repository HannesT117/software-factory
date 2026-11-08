package acceptance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/rogpeppe/go-internal/testscript"
)

// jsonsub file json checks that the JSON in file contains json: an object
// holds at least the given keys, each with a matching value; an array has
// the same length and matches element by element; any other value is
// equal. With !, it checks that the JSON does not contain it.
func jsonsub(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: jsonsub file json")
	}
	actual := decodeJSON(ts, args[0], []byte(ts.ReadFile(args[0])))
	want := decodeJSON(ts, "argument", []byte(args[1]))
	where := contains(actual, want, "")
	switch {
	case neg && where == "":
		ts.Fatalf("%s contains %s", args[0], args[1])
	case !neg && where != "":
		ts.Fatalf("%s does not contain %s: differs at %s", args[0], args[1], where)
	}
}

// jsoneq actual expected field... checks that each named top-level field
// of the object in actual equals the same field in expected. A field must
// be present in both.
func jsoneq(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! jsoneq")
	}
	if len(args) < 3 {
		ts.Fatalf("usage: jsoneq actual expected field...")
	}
	actual, ok := decodeJSON(ts, args[0], []byte(ts.ReadFile(args[0]))).(map[string]any)
	if !ok {
		ts.Fatalf("%s holds no JSON object", args[0])
	}
	expected, ok := decodeJSON(ts, args[1], []byte(ts.ReadFile(args[1]))).(map[string]any)
	if !ok {
		ts.Fatalf("%s holds no JSON object", args[1])
	}
	for _, field := range args[2:] {
		a, inActual := actual[field]
		e, inExpected := expected[field]
		if !inActual || !inExpected {
			ts.Fatalf("field %s: in %s %v, in %s %v", field, args[0], inActual, args[1], inExpected)
		}
		if !reflect.DeepEqual(a, e) {
			ts.Fatalf("field %s differs:\n%s: %s\n%s: %s", field, args[0], encodeJSON(a), args[1], encodeJSON(e))
		}
	}
}

func decodeJSON(ts *testscript.TestScript, name string, data []byte) any {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		ts.Fatalf("%s: no JSON: %v", name, err)
	}
	return v
}

func encodeJSON(v any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	return string(bytes.TrimSuffix(buf.Bytes(), []byte("\n")))
}

// contains reports where actual does not contain want, or "" if it does.
func contains(actual, want any, at string) string {
	switch w := want.(type) {
	case map[string]any:
		a, ok := actual.(map[string]any)
		if !ok {
			return orRoot(at)
		}
		keys := make([]string, 0, len(w))
		for k := range w {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v, ok := a[k]
			if !ok {
				return at + "." + k
			}
			if where := contains(v, w[k], at+"."+k); where != "" {
				return where
			}
		}
		return ""
	case []any:
		a, ok := actual.([]any)
		if !ok || len(a) != len(w) {
			return orRoot(at)
		}
		for i := range w {
			if where := contains(a[i], w[i], fmt.Sprintf("%s[%d]", at, i)); where != "" {
				return where
			}
		}
		return ""
	default:
		if !reflect.DeepEqual(actual, want) {
			return orRoot(at)
		}
		return ""
	}
}

func orRoot(at string) string {
	if at == "" {
		return "."
	}
	return at
}
