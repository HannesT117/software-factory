// Package testcheck checks that the acceptance tests match the automated
// scenarios of the design notes.
package testcheck

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
)

var (
	noteFile = regexp.MustCompile(`^docs/increments/(\d{2}-[^/]+)/design\.md$`)
	testFile = regexp.MustCompile(`^scenario-(\d{2}\.[1-9]\d*)\.txtar$`)
)

// note is a design note with the folder name of its increment.
type note struct {
	folder string
	path   string
	text   string
}

// Check compares the notes, read from base or, when base is empty, from the
// working tree at root, with the test files in the working tree at root.
func Check(r *report.Report, root, base string) int {
	notes, err := readNotes(root, base)
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}

	// expected maps folder and id to the scenario that needs a test file.
	expected := map[string]map[string]design.Scenario{}
	checked := 0
	for _, n := range notes {
		ids := map[string]design.Scenario{}
		scenarios, _ := design.Scenarios(n.text)
		for _, s := range scenarios {
			if s.Kind == "automated" {
				ids[s.ID] = s
				checked++
			}
		}
		expected[n.folder] = ids
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
		for id := range expected[n.folder] {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			if !found[n.folder][id] {
				r.Error(paths.TestDir+n.folder+"/scenario-"+id+".txtar", 0, "missing-test", "scenario "+id+" of "+n.path)
			}
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

func readNotes(root, base string) ([]note, error) {
	var notes []note
	if base == "" {
		files, err := filepath.Glob(filepath.Join(root, "docs", "increments", "*", "design.md"))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			rel := "docs/increments/" + filepath.Base(filepath.Dir(f)) + "/design.md"
			m := noteFile.FindStringSubmatch(rel)
			if m == nil {
				continue
			}
			data, err := os.ReadFile(f)
			if err != nil {
				return nil, err
			}
			notes = append(notes, note{folder: m[1], path: rel, text: string(data)})
		}
		return notes, nil
	}
	repo := gitx.Repo{Dir: root}
	files, err := repo.Paths("ls-tree", "-r", "-z", "--name-only", base, "--", "docs/increments")
	if err != nil {
		return nil, err
	}
	for _, rel := range files {
		m := noteFile.FindStringSubmatch(path.Clean(rel))
		if m == nil {
			continue
		}
		text, err := repo.Run("show", base+":"+rel)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note{folder: m[1], path: rel, text: text})
	}
	return notes, nil
}
