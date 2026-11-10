// Package testcheck checks that the acceptance tests match the automated
// scenarios of the design notes.
package testcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HannesT117/software-factory/internal/branch"
	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
)

var testFile = regexp.MustCompile(`^scenario-(\d{2}\.[1-9]\d*)\.txtar$`)

// Check compares the notes, read from base or, when base is empty, from the
// working tree at root, with the test files in the working tree at root.
// A scenario that a note replaces needs no test, and its test file is an
// error. With a base, a missing test counts only for the increments that
// HEAD belongs to against base; each other increment with missing tests
// gets a line that it waits.
func Check(r *report.Report, root, base string) int {
	notes, err := design.ReadNotes(root, base)
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}
	// own holds the increments whose missing tests are errors; nil means
	// all of them.
	var own map[string]bool
	if base != "" {
		b, err := branch.Read(gitx.Repo{Dir: root}, base)
		if err != nil {
			return r.Abort(root, "git-failed", err.Error())
		}
		own = paths.Increments(b.Changed)
	}
	replaced := design.ReplacedIDs(notes)

	// expected maps folder and id to the scenario that needs a test file.
	expected := map[string]map[string]design.Scenario{}
	checked := 0
	for _, n := range notes {
		ids := map[string]design.Scenario{}
		scenarios, _ := design.Scenarios(n.Text)
		for _, s := range scenarios {
			if s.Kind == "automated" && !replaced[s.ID] {
				ids[s.ID] = s
				checked++
			}
		}
		expected[n.Folder] = ids
	}

	found := map[string]map[string]bool{}
	dirs, _ := filepath.Glob(filepath.Join(root, "test", "acceptance", "*"))
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		folder := filepath.Base(dir)
		if _, ok := design.FolderNumber(folder); !ok {
			continue
		}
		files, _ := filepath.Glob(filepath.Join(dir, "*.txtar"))
		sort.Strings(files)
		for _, f := range files {
			rel := paths.TestDir + folder + "/" + filepath.Base(f)
			m := testFile.FindStringSubmatch(filepath.Base(f))
			if m == nil {
				r.Error(rel, 0, "unknown-id", "the file name is not scenario-<nn>.<n>.txtar")
				continue
			}
			id := m[1]
			if replaced[id] {
				r.Error(rel, 0, "replaced-test", "a later note replaces scenario "+id+"; remove its test")
				continue
			}
			s, ok := expected[folder][id]
			if !ok {
				r.Error(rel, 0, "unknown-id", "no automated scenario "+id+" in docs/increments/"+folder+"/design.md")
				continue
			}
			if found[folder] == nil {
				found[folder] = map[string]bool{}
			}
			found[folder][id] = true
			data, err := os.ReadFile(f)
			if err != nil {
				return r.Abort(rel, "unreadable", err.Error())
			}
			checkFile(r, rel, s, string(data))
		}
	}

	for _, n := range notes {
		var ids []string
		for id := range expected[n.Folder] {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		waiting := 0
		for _, id := range ids {
			if found[n.Folder][id] {
				continue
			}
			if own != nil && !own[n.Folder] {
				waiting++
				continue
			}
			r.Error(paths.TestDir+n.Folder+"/scenario-"+id+".txtar", 0, "missing-test", "scenario "+id+" of "+n.Path)
		}
		if waiting > 0 {
			r.Line(fmt.Sprintf("waiting %s %d scenarios without tests", n.Folder, waiting))
		}
	}
	return r.Finish(checked)
}

// checkFile checks line 1 and the commands of one test file.
func checkFile(r *report.Report, rel string, s design.Scenario, text string) {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	want := "# scenario " + s.ID + ": " + s.Title
	if lines[0] != want {
		r.Error(rel, 1, "title-mismatch", "line 1 must be: "+want)
	}
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "-- ") && strings.HasSuffix(t, " --") {
			break // the files of the archive start here
		}
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		fields := strings.Fields(t)
		if fields[0] == "!" || fields[0] == "?" {
			fields = fields[1:]
		}
		if len(fields) == 0 {
			continue
		}
		cmd := strings.TrimLeft(fields[0], "!?")
		switch {
		case strings.HasPrefix(cmd, "["):
			r.Error(rel, i+1, "forbidden-command", "condition "+cmd+" in scenario "+s.ID)
		case cmd == "skip" || cmd == "stop":
			r.Error(rel, i+1, "forbidden-command", cmd+" in scenario "+s.ID)
		}
	}
}
