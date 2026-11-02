# Claude Code adapter

This folder translates the factory for Claude Code. `factory install` copies it into a project from increment 03 on.

- `settings.json`: the hooks. In a project, they go into `.claude/settings.json`.
  - `PreToolUse` runs `factory hook pre-tool-use`. It denies writes to the tests group and the CI-and-config group, except test files from a session started as `acceptance-test-writer`. It denies changes to the files of a skill until the session has loaded `authoring-skills`. Outside a design session (`FACTORY_SESSION=design`), it denies web searches and pages from hosts that `.factory/allowed-hosts` does not list.
  - `Stop` runs `factory hook stop`. It keeps the session from ending while the latest code commit has no review result.
- `user-settings.json`: a template for the Claude Code user settings of the engineer's computer, for local agent sessions. The sandbox keeps shell commands from reading `~/.ssh` and `~/.config/gh`, and from reaching hosts outside its list; the deny rules do the same for the Read and Edit tools, which the sandbox does not cover. The SSH agent socket stays closed, because the template allows no Unix socket.
- `agents/acceptance-test-writer.md`: the fixed prompt of the test writer. The engineer starts a session with `claude --agent acceptance-test-writer`. Only such a session may write under `test/acceptance/`.
- `agents/increment-reviewer.md`: the fixed prompt of the reviewer. The implementing session starts it as a subagent when the stop hook blocks. It publishes the review result as a Claude artifact.
- `authoring-skills/`: the parts of the test loop of `authoring-skills` that call Claude Code.
- `hooks.go`: reads the hook input of Claude Code, including the session transcript, and writes its decisions. The rules themselves live in the tool-neutral packages under `internal/`.

In a project, the agents go into `.claude/agents/`, and the skills from `core/skills/` reach `.claude/skills/` as symlinks to `.agents/skills/` (ADR 0002).
