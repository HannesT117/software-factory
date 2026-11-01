// Package claudecode translates Claude Code hook input and output for the
// factory's guard and review checks.
package claudecode

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/guard"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/review"
)

// writeTools are the tools whose file writes the guard checks.
var writeTools = map[string]bool{"Write": true, "Edit": true, "MultiEdit": true, "NotebookEdit": true}

// skillTools are the tools whose changes to the files of a skill need the
// skill for skills.
var skillTools = map[string]bool{"Write": true, "Edit": true, "MultiEdit": true}

type preToolUseInput struct {
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		FilePath     string `json:"file_path"`
		NotebookPath string `json:"notebook_path"`
		URL          string `json:"url"`
	} `json:"tool_input"`
	AgentType      *string `json:"agent_type"`
	AgentID        *string `json:"agent_id"`
	Cwd            string  `json:"cwd"`
	TranscriptPath string  `json:"transcript_path"`
}

type stopInput struct {
	StopHookActive bool   `json:"stop_hook_active"`
	Cwd            string `json:"cwd"`
}

// PreToolUse reads the hook input from in and writes a deny decision to out
// when the tool call writes a protected path, changes a skill without the
// skill for skills, or reads the web outside a design session.
func PreToolUse(in io.Reader, out io.Writer) int {
	var input preToolUseInput
	if err := json.NewDecoder(in).Decode(&input); err != nil {
		return report.New("hook", out).Abort("stdin", "bad-input", err.Error())
	}
	root := projectRoot(input.Cwd)

	if input.ToolName == "WebSearch" || input.ToolName == "WebFetch" {
		if os.Getenv("FACTORY_SESSION") == guard.DesignSession {
			return report.OK
		}
		hosts, err := guard.ReadHosts(root)
		if err != nil {
			return deny(out, "The guard cannot read "+guard.HostsFile+": "+err.Error())
		}
		host := ""
		if u, err := url.Parse(input.ToolInput.URL); err == nil {
			host = u.Hostname()
		}
		if reason := guard.WebReason(input.ToolName == "WebSearch", host, hosts); reason != "" {
			return deny(out, reason)
		}
		return report.OK
	}

	if !writeTools[input.ToolName] {
		return report.OK
	}
	file := input.ToolInput.FilePath
	if file == "" {
		file = input.ToolInput.NotebookPath
	}
	if file == "" {
		return report.OK
	}
	if !filepath.IsAbs(file) {
		file = filepath.Join(input.Cwd, file)
	}
	file = filepath.Clean(file)
	if rel, err := filepath.Rel(root, file); err == nil {
		// A session started with --agent has agent_type and no agent_id. A
		// subagent has both.
		testWriter := input.AgentType != nil && *input.AgentType == guard.TestWriter && input.AgentID == nil
		if reason := guard.Check(filepath.ToSlash(rel), testWriter); reason != "" {
			return deny(out, reason)
		}
	}
	if skillTools[input.ToolName] && guard.IsSkillFile(file) && !loadedSkillForSkills(input.TranscriptPath) {
		return deny(out, guard.SkillReason(file))
	}
	return report.OK
}

// loadedSkillForSkills reports whether the transcript, in JSON lines, holds
// a call of the tool Skill for the skill for skills, or a Read of its
// SKILL.md.
func loadedSkillForSkills(transcript string) bool {
	if transcript == "" {
		return false
	}
	f, err := os.Open(transcript)
	if err != nil {
		return false
	}
	defer f.Close()
	rd := bufio.NewReaderSize(f, 1<<20)
	marker := []byte(guard.SkillForSkills)
	for {
		line, err := rd.ReadBytes('\n')
		if bytes.Contains(line, marker) && lineLoads(line) {
			return true
		}
		if err != nil {
			return false
		}
	}
}

type toolUse struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Input struct {
		Skill    string `json:"skill"`
		FilePath string `json:"file_path"`
	} `json:"input"`
}

func lineLoads(line []byte) bool {
	var entry struct {
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(line, &entry) != nil {
		return false
	}
	var uses []toolUse
	if json.Unmarshal(entry.Message.Content, &uses) != nil {
		return false
	}
	for _, u := range uses {
		if u.Type != "tool_use" {
			continue
		}
		if u.Name == "Skill" && u.Input.Skill == guard.SkillForSkills {
			return true
		}
		if u.Name == "Read" && filepath.Base(u.Input.FilePath) == "SKILL.md" &&
			filepath.Base(filepath.Dir(u.Input.FilePath)) == guard.SkillForSkills {
			return true
		}
	}
	return false
}

// Stop reads the hook input from in and writes a block decision to out
// while the latest code commit has no review result.
func Stop(in io.Reader, out io.Writer) int {
	var input stopInput
	if err := json.NewDecoder(in).Decode(&input); err != nil {
		return report.New("hook", out).Abort("stdin", "bad-input", err.Error())
	}
	if input.StopHookActive {
		return report.OK
	}
	root := projectRoot(input.Cwd)
	commit, err := review.Pending(root)
	if err != nil {
		// Exit code 1 shows the error to the person and lets the session
		// end; exit code 2 would block it.
		report.New("hook", out).Error(root, 0, "git-failed", err.Error())
		return report.Fail
	}
	if commit == "" {
		return report.OK
	}
	return writeJSON(out, map[string]any{
		"decision": "block",
		"reason": "Commit " + commit + " changes code and has no review result yet. Start the reviewer: run the subagent " +
			"increment-reviewer, which uses the skill reviewing-increments. It writes evidence/<nn>-<slug>/review.md, " +
			"which names the full hash of the reviewed commit; commit and push that file.",
	})
}

func deny(out io.Writer, reason string) int {
	return writeJSON(out, map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":            "PreToolUse",
			"permissionDecision":       "deny",
			"permissionDecisionReason": reason,
		},
	})
}

// projectRoot returns the root of the git working tree that contains cwd,
// or cwd itself.
func projectRoot(cwd string) string {
	if cwd == "" {
		cwd = "."
	}
	abs, err := filepath.Abs(cwd)
	if err != nil {
		return cwd
	}
	if _, err := os.Stat(abs); err != nil {
		return abs
	}
	if top, err := gitx.Toplevel(abs); err == nil {
		return top
	}
	return abs
}

func writeJSON(out io.Writer, v any) int {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return report.Usage
	}
	return report.OK
}
