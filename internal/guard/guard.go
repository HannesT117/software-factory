// Package guard decides which file writes an agent may make.
package guard

import (
	"path"
	"strings"

	"github.com/HannesT117/software-factory/internal/paths"
)

// TestWriter is the role of the session that writes acceptance tests.
const TestWriter = "acceptance-test-writer"

// Check returns why a write to rel, a slash-separated path relative to the
// repository root, is denied, or "" when it is allowed. testWriter is true
// for a session started in the role TestWriter, not for a helper that such
// a session started.
func Check(rel string, testWriter bool) string {
	rel = path.Clean(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return ""
	}
	switch paths.GroupOf(rel) {
	case paths.Tests:
		if testWriter && strings.HasPrefix(rel, paths.TestDir) {
			return ""
		}
		return rel + " is in the tests group. Only a session started as " + TestWriter +
			" writes acceptance tests, and nobody else changes go.mod or go.sum. If a test looks wrong, stop and ask the engineer."
	case paths.CIConfig:
		return rel + " is in the CI-and-config group. Only the engineer changes it, in a pull request of its own."
	}
	return ""
}
