// Package crossreview asks a model of another vendor to review a pull
// request, through an OpenAI-compatible chat completions API.
package crossreview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/paths"
	"github.com/HannesT117/software-factory/internal/report"
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

const answerFormat = `Answer with only one JSON object and no other text: {"findings": [{"severity": "blocking" or "note", "path": "<path in the repository>", "text": "<one or two sentences>"}]}. An empty list means no findings.`

var incrementDir = regexp.MustCompile(`^(?:test/acceptance|evidence)/(\d{2}-[^/]+)/`)

// Run reviews HEAD of the repo at root against base and returns the exit
// code: 0 without a blocking finding, 1 with one, 2 for a usage or
// environment problem. A model that cannot be reached, or that answers in
// another form, counts as blocking.
func Run(r *report.Report, root, base, key string) int {
	repo := gitx.Repo{Dir: root}
	cfgText, err := repo.Run("show", base+":"+ConfigFile)
	if err != nil {
		return blocked(r, ConfigFile, "no-config", "no "+ConfigFile+" on "+base, report.Usage)
	}
	var cfg Config
	if err := json.Unmarshal([]byte(cfgText), &cfg); err != nil || cfg.Endpoint == "" || cfg.Models.Design == "" || cfg.Models.Code == "" {
		return blocked(r, ConfigFile, "bad-config", `needs {"endpoint": "...", "models": {"design": "...", "code": "..."}}`, report.Usage)
	}
	if key == "" {
		return blocked(r, KeyVar, "no-key", "the API key is missing", report.Usage)
	}

	mb, err := repo.Run("merge-base", base, "HEAD")
	if err != nil {
		return blocked(r, base, "git-failed", err.Error(), report.Usage)
	}
	mb = strings.TrimSpace(mb)
	changed, err := repo.Paths("diff", "--name-only", "--no-renames", "-z", mb, "HEAD")
	if err != nil {
		return blocked(r, root, "git-failed", err.Error(), report.Usage)
	}

	model, instructions, content, err := request(repo, cfg, mb, changed)
	if err != nil {
		return blocked(r, root, "git-failed", err.Error(), report.Usage)
	}
	findings, err := ask(cfg.Endpoint, key, model, instructions+"\n\n"+answerFormat, content)
	if err != nil {
		return blocked(r, cfg.Endpoint, "no-review", err.Error(), report.Fail)
	}

	blocking := false
	for _, f := range findings {
		blocking = blocking || f.Severity == "blocking"
		r.Line(fmt.Sprintf("finding %s %s: %s", f.Severity, oneLine(f.Path), oneLine(f.Text)))
	}
	if blocking {
		r.Line("blocking")
		return report.Fail
	}
	r.Line("ok")
	return report.OK
}

// request builds the model, the instructions and the content: the changed
// design notes for a design pull request; otherwise the notes of the
// increment, its scenario titles and the diff.
func request(repo gitx.Repo, cfg Config, mb string, changed []string) (model, instructions, content string, err error) {
	var designFiles []string
	folders := map[string]bool{}
	for _, p := range changed {
		if paths.GroupOf(p) == paths.DesignNotes {
			designFiles = append(designFiles, p)
		}
		if m := incrementDir.FindStringSubmatch(p); m != nil {
			folders[m[1]] = true
		}
	}

	var b strings.Builder
	if len(designFiles) > 0 {
		for _, f := range designFiles {
			text, err := repo.Run("show", "HEAD:"+f)
			if err != nil {
				text = "(deleted)"
			}
			fmt.Fprintf(&b, "## File %s\n\n%s\n\n", f, text)
		}
		return cfg.Models.Design, designInstructions, b.String(), nil
	}

	notes, err := design.ReadNotes(repo.Dir, "HEAD")
	if err != nil {
		return "", "", "", err
	}
	var chosen []design.Note
	for _, n := range notes {
		if folders[n.Folder] {
			chosen = append(chosen, n)
		}
	}
	if len(chosen) == 0 && len(notes) > 0 {
		chosen = notes[len(notes)-1:]
	}
	for _, n := range chosen {
		fmt.Fprintf(&b, "## Design note %s\n\n%s\n\n", n.Path, n.Text)
	}
	b.WriteString("## Scenario titles\n\n")
	for _, s := range design.InForce(chosen) {
		fmt.Fprintf(&b, "- %s (%s) %s\n", s.ID, s.Kind, s.Title)
	}
	diff, err := repo.Run("diff", "--no-renames", mb, "HEAD")
	if err != nil {
		return "", "", "", err
	}
	fmt.Fprintf(&b, "\n## Diff\n\n%s\n", diff)
	return cfg.Models.Code, codeInstructions, b.String(), nil
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

// blocked reports a problem, the result blocking, and returns code.
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
