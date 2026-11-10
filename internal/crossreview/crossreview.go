// Package crossreview asks a model of another vendor to review a pull
// request, through an OpenAI-compatible chat completions API.
package crossreview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HannesT117/software-factory/internal/branch"
	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/review"
)

// ConfigFile holds the endpoint and the models.
const ConfigFile = ".factory/cross-review.json"

// KeyVar is the environment variable with the API key.
const KeyVar = "FACTORY_MODEL_KEY"

// Timeout is how long the command waits for the answer.
const Timeout = 10 * time.Minute

// Config is the content of ConfigFile.
type Config struct {
	Endpoint string `json:"endpoint"`
	Models   struct {
		Design string `json:"design"`
		Code   string `json:"code"`
	} `json:"models"`
}

// Finding is one finding of the model.
type Finding struct {
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Text     string `json:"text"`
}

const designInstructions = `You review the design note of a software increment before any test or code exists. Each scenario becomes an acceptance test written by a separate agent that sees only the note and the public interface. Look for scenarios that are ambiguous, cannot be tested from the interface, contradict each other or the rest of the note, or leave out a case that the goal needs. A finding is blocking only when a test writer could not write a correct test from the note.`

const codeInstructions = `You review a code change of a software increment against its design note. The acceptance tests were written first, from the note alone. Look for code that detects a test run or hard-codes test cases, tests that do not check what their scenario says, special cases that the scenarios do not cover, work outside the scenarios, and defects that a reader at 2 am would need to debug. A finding is blocking only when the change is wrong or a test does not test its scenario.`

const testsInstructions = `You review the acceptance tests of a software increment before its code exists. A separate agent wrote them from the design note alone; after the first code commit they can no longer change. Check that each test checks what its scenario says, including each case where something must not happen, and that it would fail without the behaviour of its scenario. A finding is blocking only when a test would pass without the behaviour of its scenario, or would fail for a correct implementation.`

const ownInstructions = `The change belongs to no increment, so no design note comes with it. Judge the change by itself.`

const answerFormat = `Answer with only one JSON object and no other text: {"findings": [{"severity": "blocking" or "note", "path": "<path in the repository>", "text": "<one or two sentences>"}]}. An empty list means no findings.`

// Kind is the review that a pull request gets.
type Kind string

// The kinds of review.
const (
	Design Kind = "design"
	Tests  Kind = "tests"
	Code   Kind = "code"
	None   Kind = "none"
)

// KindOf chooses the review from the paths that a pull request changes.
func KindOf(changed []string, code paths.Code) Kind {
	kind := None
	for _, p := range changed {
		if paths.GroupOf(p) == paths.DesignNotes {
			return Design
		}
		if !code.Is(p) {
			continue
		}
		if paths.GroupOf(p) == paths.Tests {
			if kind == None {
				kind = Tests
			}
		} else {
			kind = Code
		}
	}
	return kind
}

// Run reviews HEAD of the repo at root against base and returns the exit
// code: 0 without a blocking finding, 1 with one, 2 for a usage or
// environment problem. A model that cannot be reached, or that answers in
// another form, counts as blocking. A design review never fails for a
// finding, and a pull request without design notes and code is not sent.
func Run(r *report.Report, root, base, key string) int {
	repo := gitx.Repo{Dir: root}
	if _, err := repo.Run("rev-parse", "--verify", "--quiet", base+"^{commit}"); err != nil {
		return r.Abort(base, "unknown-ref", "no commit "+base)
	}
	code, err := review.ReadCodeAt(repo, base)
	if err != nil {
		return r.Abort(paths.CodePathsFile, "unreadable", short(err.Error()))
	}
	b, err := branch.Read(repo, base)
	if err != nil {
		return r.Abort(base, "git-failed", short(err.Error()))
	}
	kind := KindOf(b.Changed, code)
	switch kind {
	case None:
		r.Line("does not apply: no design note and no code change")
		r.Line("ok")
		return report.OK
	case Code:
		if h := b.LatestCode(code); h != "" {
			file, err := review.Find(root, h)
			if err != nil {
				return r.Abort(root, "unreadable", short(err.Error()))
			}
			if file == "" {
				r.Line("waiting for a review result for " + h)
				r.Line("blocking")
				return report.Fail
			}
		}
	}

	cfgText, err := repo.Run("show", base+":"+ConfigFile)
	if err != nil {
		return r.Abort(ConfigFile, "no-config", "no "+ConfigFile+" on "+base)
	}
	var cfg Config
	if err := json.Unmarshal([]byte(cfgText), &cfg); err != nil || cfg.Endpoint == "" || cfg.Models.Design == "" || cfg.Models.Code == "" {
		return r.Abort(ConfigFile, "no-config", `on `+base+` it is not {"endpoint": "...", "models": {"design": "...", "code": "..."}}`)
	}
	if key == "" {
		return r.Abort(KeyVar, "no-key", "the API key is missing or empty")
	}

	model, instructions, content, err := request(repo, cfg, kind, b.MergeBase, b.Changed)
	if err != nil {
		return r.Abort(root, "git-failed", short(err.Error()))
	}
	r.Line("model " + model)
	findings, err := ask(cfg.Endpoint, key, model, instructions+"\n\n"+answerFormat, content)
	if err != nil {
		return blocked(r, cfg.Endpoint, "no-review", err.Error(), report.Fail)
	}

	blocking := false
	for _, f := range findings {
		blocking = blocking || f.Severity == "blocking"
		r.Line(fmt.Sprintf("finding %s %s: %s", f.Severity, oneLine(f.Path), oneLine(f.Text)))
	}
	if blocking && kind != Design {
		r.Line("blocking")
		return report.Fail
	}
	r.Line("ok")
	return report.OK
}

// request builds the model, the instructions and the content: the changed
// design notes for a design review; otherwise the notes of the increments
// that the pull request belongs to, their scenario titles and the diff.
func request(repo gitx.Repo, cfg Config, kind Kind, mb string, changed []string) (model, instructions, content string, err error) {
	var b strings.Builder
	if kind == Design {
		for _, f := range changed {
			if paths.GroupOf(f) != paths.DesignNotes {
				continue
			}
			text, err := repo.Run("show", "HEAD:"+f)
			if err != nil {
				text = "(deleted)"
			}
			fmt.Fprintf(&b, "## File %s\n\n%s\n\n", f, text)
		}
		return cfg.Models.Design, designInstructions, b.String(), nil
	}

	instructions = codeInstructions
	if kind == Tests {
		instructions = testsInstructions
	}
	notes, err := design.ReadNotes(repo.Dir, "HEAD")
	if err != nil {
		return "", "", "", err
	}
	own := paths.Increments(changed)
	var chosen []design.Note
	for _, n := range notes {
		if own[n.Folder] {
			chosen = append(chosen, n)
		}
	}
	if len(chosen) == 0 {
		instructions += " " + ownInstructions
	}
	for _, n := range chosen {
		fmt.Fprintf(&b, "## Design note %s\n\n%s\n\n", n.Path, n.Text)
	}
	if len(chosen) > 0 {
		b.WriteString("## Scenario titles\n\n")
		for _, s := range design.InForce(notes) {
			if !own[folderOf(notes, s.ID)] {
				continue
			}
			fmt.Fprintf(&b, "- %s (%s) %s\n", s.ID, s.Kind, s.Title)
		}
	}
	diff, err := repo.Run("diff", "--no-renames", mb, "HEAD")
	if err != nil {
		return "", "", "", err
	}
	fmt.Fprintf(&b, "\n## Diff\n\n%s\n", diff)
	return cfg.Models.Code, instructions, b.String(), nil
}

// folderOf returns the folder of the note whose increment number starts id.
func folderOf(notes []design.Note, id string) string {
	nn, _, _ := strings.Cut(id, ".")
	for _, n := range notes {
		if n.Number == nn {
			return n.Folder
		}
	}
	return ""
}

// ask sends the request and returns the findings of a well-formed answer.
func ask(endpoint, key, model, system, user string) ([]Finding, error) {
	body, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": user},
		},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(endpoint, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := (&http.Client{Timeout: Timeout}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, short(string(data)))
	}

	var completion struct {
		Choices []struct {
			Message struct {
				Content *string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &completion); err != nil || len(completion.Choices) == 0 || completion.Choices[0].Message.Content == nil {
		return nil, fmt.Errorf("the answer is not a chat completion")
	}
	var answer struct {
		Findings *[]Finding `json:"findings"`
	}
	content := strings.TrimSpace(*completion.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &answer); err != nil || answer.Findings == nil {
		return nil, fmt.Errorf("the answer is not the JSON object with findings: %s", short(content))
	}
	for _, f := range *answer.Findings {
		if f.Severity != "blocking" && f.Severity != "note" {
			return nil, fmt.Errorf("a finding has the severity %q, not blocking or note", f.Severity)
		}
	}
	return *answer.Findings, nil
}

// blocked reports why the review did not happen and the result blocking.
func blocked(r *report.Report, where, code, detail string, exit int) int {
	r.Error(where, 0, code, short(detail))
	r.Line("blocking")
	return exit
}

// oneLine joins the lines of s, so that each output line stays grep-able.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// short is oneLine, cut to 300 bytes, for the detail of a problem.
func short(s string) string {
	s = oneLine(s)
	if len(s) > 300 {
		s = s[:300] + "..."
	}
	return s
}
