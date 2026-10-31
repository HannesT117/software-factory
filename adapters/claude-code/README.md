# Claude Code adapter

This folder translates the factory for Claude Code. `factory install` copies it into a project from increment 03 on.

- `settings.json`: the hooks. `PreToolUse` runs `factory hook pre-tool-use`, which denies writes to the tests group and the CI-and-config group. `Stop` runs `factory hook stop`, which keeps the session from ending while the latest code commit has no review result. In a project, they go into `.claude/settings.json`.
- `agents/acceptance-test-writer.md`: the fixed prompt of the test writer. The engineer starts a session with `claude --agent acceptance-test-writer`. Only such a session may write under `test/acceptance/`.
- `agents/increment-reviewer.md`: the fixed prompt of the reviewer. The implementing session starts it as a subagent when the stop hook blocks. It publishes the review result as a Claude artifact.
- `hooks.go`: reads the hook input of Claude Code and writes its decisions. The rules themselves live in the tool-neutral packages under `internal/`.

In a project, the agents go into `.claude/agents/`, and the skills from `core/skills/` reach `.claude/skills/` as symlinks to `.agents/skills/` (ADR 0002).
