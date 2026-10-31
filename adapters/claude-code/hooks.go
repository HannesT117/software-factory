// Package claudecode translates Claude Code hook input and output for the
// factory's guard and review checks.
package claudecode

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/guard"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/review"
)

// writeTools are the tools whose file writes the guard checks.
var writeTools = map[string]bool{"Write": true, "Edit": true, "MultiEdit": true, "NotebookEdit": true}

type preToolUseInput struct {
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		FilePath     string `json:"file_path"`
		NotebookPath string `json:"notebook_path"`
	} `json:"tool_input"`
	AgentType *string `json:"agent_type"`
	AgentID   *string `json:"agent_id"`
	Cwd       string  `json:"cwd"`
}

type stopInput struct {
	StopHookActive bool   `json:"stop_hook_active"`
	Cwd            string `json:"cwd"`
}

// PreToolUse reads the hook input from in and writes a deny decision to out
// when the tool call writes a protected path.
func PreToolUse(in io.Reader, out io.Writer) int {
	var input preToolUseInput
	if err := json.NewDecoder(in).Decode(&input); err != nil {
		return report.New("hook", out).Abort("stdin", "bad-input", err.Error())
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
	root := projectRoot(input.Cwd)
	if !filepath.IsAbs(file) {
		file = filepath.Join(input.Cwd, file)
	}
	rel, err := filepath.Rel(root, filepath.Clean(file))
	if err != nil {
		return report.OK
	}
	// A session started with --agent has agent_type and no agent_id. A
	// subagent has both.
	testWriter := input.AgentType != nil && *input.AgentType == guard.TestWriter && input.AgentID == nil
	reason := guard.Check(filepath.ToSlash(rel), testWriter)
	if reason == "" {
		return report.OK
	}
	return writeJSON(out, map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":            "PreToolUse",
			"permissionDecision":       "deny",
			"permissionDecisionReason": reason,
		},
	})
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
		return report.New("hook", out).Abort(root, "git-failed", err.Error())
	}
	if commit == "" {
		return report.OK
	}
	return writeJSON(out, map[string]any{
		"decision": "block",
		"reason": "Commit " + commit + " has no review result yet. Start the reviewer: run the subagent " +
			"increment-reviewer with the skill reviewing-increments. It writes evidence/<nn>-<slug>/review.md, " +
			"which names the full hash of the reviewed commit.",
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
