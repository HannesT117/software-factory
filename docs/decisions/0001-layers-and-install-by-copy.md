---
status: accepted
date: 2026-10-26
decision-makers: engineer
---

# Vendor-neutral core with thin adapters, installed into projects by copy

## Context and Problem Statement

The factory consists of skills, scripts, hooks and templates that projects use. Coding agents load these parts in different ways, and agent tools change faster than the projects that use them. How is the factory structured and delivered to projects, so that the agent tool stays replaceable and every project works in cloud sessions and CI?

## Decision Drivers

* Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read skills from `.agents/skills/`. Claude Code reads only `.claude/skills/`, and it loads a `<skill-name>` entry there that is a symlink.
* Hooks have no cross-tool standard. Each tool has its own format.
* Harness changes alone can lower output quality: three Claude Code changes lowered it for weeks in 2026 while the model stayed the same (Anthropic postmortem, 2026-04-23).
* Cloud sessions and CI see only the project repository. Cloud sessions do not load plugins from local settings.
* Each project pins the factory version it uses, and every update is reviewable.

## Considered Options

1. Vendor-neutral core plus thin adapters, installed by copy
2. Claude Code plugin from a marketplace
3. Symlinks from the project to a local factory checkout
4. Git submodule of the factory in each project

## Decision Outcome

Chosen option 1, because it is the only option that works in cloud sessions and CI without extra setup, pins a reviewable version per project, and keeps tool-specific parts in one replaceable place.

* **Core:**
  * Agent Skills with standard frontmatter fields only
  * Node `.mjs` scripts without dependencies
  * git hook, CI template and file templates
* **Adapters:** `adapters/<tool>/` holds settings and hooks that only call core scripts.
* **Install:** `install.mjs` copies the skills into the project's `.agents/skills/` with a version stamp. For Claude Code it also creates relative symlinks `.claude/skills/<name>` → `../../.agents/skills/<name>`. The project commits the result.
* **Check:** `install.mjs --check` fails when vendored files were edited in the project. Fixes go to the factory.

### Consequences

* Good, because cloud sessions and CI work with the project repository alone.
* Good, because each update shows up as a diff in the project.
* Good, because replacing the agent tool means adding one adapter.
* Bad, because every update needs a re-run of `install.mjs`.
* Bad, because the skills exist twice, in the factory and in each project.
* Bad, because Windows clones need `core.symlinks`.
* Neutral, because project cloud threads do not load committed `.claude/skills/`. Other cloud sessions do.

### Confirmation

* A check fails if a file under `core/` names a specific agent tool.
* `install.mjs --check` runs in the CI of every project.
* A test installs into a fixture project. It checks that Claude Code lists the skills through the symlinks.

## Pros and Cons of the Options

### Claude Code plugin

* Good, because it installs with one command and updates through the marketplace.
* Bad, because it works only in Claude Code, and cloud sessions do not load local plugins.

### Symlinks to a local factory checkout

* Good, because factory changes show up in the project at once.
* Bad, because the links break in cloud sessions and CI, and no version is pinned.

### Git submodule

* Good, because it pins a version from one source.
* Bad, because agents still need files at `.agents/skills/` and `.claude/skills/`, so the project needs links into the submodule anyway.
* Bad, because it is not verified that cloud sessions fetch submodules.

## More Information

* **Confidence:** medium. The symlink behaviour comes from the Claude Code docs and is tested in increment 07.
* **Revisit when** a cross-tool hook standard appears, or Claude Code reads `.agents/skills/`.
* **Sources:** agentskills.io specification; vercel-labs/skills agent paths; Claude Code docs on skills, memory and plugins; Anthropic postmortem 2026-04-23. Details are in the [research notes](../research/2026-10-factory-basis.md).
