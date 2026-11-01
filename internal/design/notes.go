package design

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/HannesT117/software-factory/internal/gitx"
)

var noteFile = regexp.MustCompile(`^docs/increments/((\d{2})-[^/]+)/design\.md$`)

// Note is a design note with the folder of its increment.
type Note struct {
	Folder string // for example "02-trust-core"
	Number string // for example "02"
	Path   string // relative to the repository root
	Text   string
}

// ReadNotes reads the notes under docs/increments of the repository at
// root: from the commit base, or from the working tree when base is empty.
// The notes come in the order of their folders.
func ReadNotes(root, base string) ([]Note, error) {
	var notes []Note
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
			notes = append(notes, Note{Folder: m[1], Number: m[2], Path: rel, Text: string(data)})
		}
	} else {
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
			notes = append(notes, Note{Folder: m[1], Number: m[2], Path: rel, Text: text})
		}
	}
	sort.Slice(notes, func(i, j int) bool { return notes[i].Folder < notes[j].Folder })
	return notes, nil
}

// ReplacedIDs returns the ids of the scenarios that any of the notes
// replaces.
func ReplacedIDs(notes []Note) map[string]bool {
	ids := map[string]bool{}
	for _, n := range notes {
		replaced, _ := ReplacedScenarios(n.Text)
		for _, r := range replaced {
			ids[r.ID] = true
		}
	}
	return ids
}

// InForce returns the scenarios of the notes that no note replaces, in the
// order of their increments and numbers.
func InForce(notes []Note) []Scenario {
	replaced := ReplacedIDs(notes)
	var list []Scenario
	for _, n := range notes {
		scenarios, _ := Scenarios(n.Text)
		for _, s := range scenarios {
			if !replaced[s.ID] {
				list = append(list, s)
			}
		}
	}
	sort.SliceStable(list, func(i, j int) bool { return less(list[i].ID, list[j].ID) })
	return list
}

// less orders scenario ids by increment, then by number.
func less(a, b string) bool {
	ai, an, _ := strings.Cut(a, ".")
	bi, bn, _ := strings.Cut(b, ".")
	if ai != bi {
		return ai < bi
	}
	x, _ := strconv.Atoi(an)
	y, _ := strconv.Atoi(bn)
	return x < y
}
