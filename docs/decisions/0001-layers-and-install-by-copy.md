# 0001: Vendor-neutral core, thin adapters, install by copy

- Status: accepted
- Date: 2026-10-26

## Context

The engineer wants no vendor lock-in. Coding agents differ in how they load instructions, skills and hooks:

- Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read skills from `.agents/skills/`.
- Claude Code reads skills from `.claude/skills/` and nothing under `.agents/`. A `<skill-name>` entry there may be a symlink.
- Hooks have no shared standard. Each tool has its own format.
- Claude Code plugins work only in Claude Code, and cloud sessions do not load plugins from local settings.

## Decision

1. **Core**, usable by any agent:
   - skills in the Agent Skills format, with standard frontmatter fields only
   - Node `.mjs` scripts without dependencies
   - git hook, CI template and file templates
2. **Adapters** are thin. `adapters/claude-code/` holds settings and hooks that only call core scripts. Claude-only features such as `/advisor` are documented there, never required.
3. **Install by copy.** `install.mjs` copies the skills into the project's `.agents/skills/` with a version stamp. For each skill it creates a relative symlink `.claude/skills/<name>` → `../../.agents/skills/<name>`. The project commits the result.
4. **No link to the factory checkout.** `install.mjs --check` fails in CI when vendored files were edited in the project. Fixes go to this repo.

## Consequences

Good:

- Cloud sessions and CI work with the project repo alone.
- Each project pins one factory version. Every update shows up as a reviewable diff.
- Switching to another agent means writing a new adapter only.

Bad:

- Updating a project needs a re-run of `install.mjs`. A `--watch` mode is in the backlog.
- Windows clones need `core.symlinks`.
- Project cloud threads do not load committed `.claude/skills/`.

## Re-check when

A cross-tool standard for hooks appears, or Claude Code starts reading `.agents/skills/`.

## Sources

agentskills.io specification; vercel-labs/skills agent paths; Claude Code docs on skills, memory and plugins. Details in [research](../research/2026-10-factory-basis.md).
