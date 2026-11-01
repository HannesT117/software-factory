// Package guard decides which file writes and web requests an agent may
// make.
package guard

import (
	"os"
	"path"
	"path/filepath"
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

// HostsFile lists the hosts that a session outside a design session may
// fetch pages from.
const HostsFile = ".factory/allowed-hosts"

// DesignSession is the value of FACTORY_SESSION in a design session, the
// only session that reads the open web.
const DesignSession = "design"

// SkillForSkills is the skill that a session loads before it changes the
// files of a skill.
const SkillForSkills = "authoring-skills"

// ReadHosts reads the allowed hosts of the project at root. Without the
// file, no host is allowed.
func ReadHosts(root string) ([]string, error) {
	data, err := os.ReadFile(filepath.Join(root, HostsFile))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return paths.ParseList(string(data)), nil
}

// HostAllowed reports whether host is in the list. An entry *.example.com
// matches every subdomain of example.com.
func HostAllowed(host string, list []string) bool {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if host == "" {
		return false
	}
	for _, entry := range list {
		entry = strings.ToLower(entry)
		if suffix, ok := strings.CutPrefix(entry, "*."); ok {
			if strings.HasSuffix(host, "."+suffix) {
				return true
			}
		} else if host == entry {
			return true
		}
	}
	return false
}

// WebReason returns why a web tool call is denied outside a design
// session, or "" when it is allowed. search is true for a web search;
// otherwise host is the host of the page to fetch.
func WebReason(search bool, host string, hosts []string) string {
	if search {
		return "Only a design session searches the web (FACTORY_SESSION=design). Ask the engineer, or use a page from a host in " + HostsFile + "."
	}
	if HostAllowed(host, hosts) {
		return ""
	}
	return "The host " + host + " is not in " + HostsFile + ". Only a design session (FACTORY_SESSION=design) reads pages from other hosts. Ask the engineer to add the host, or to look the page up."
}

// IsSkillFile reports whether file, an absolute path, belongs to a skill:
// it is a SKILL.md, or a folder above it holds a SKILL.md.
func IsSkillFile(file string) bool {
	if filepath.Base(file) == "SKILL.md" {
		return true
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}

// SkillReason is why a change to the files of a skill is denied before the
// session has loaded the skill for skills.
func SkillReason(file string) string {
	return file + " belongs to a skill. Load the skill " + SkillForSkills + " first, and follow it to create or change a skill."
}
