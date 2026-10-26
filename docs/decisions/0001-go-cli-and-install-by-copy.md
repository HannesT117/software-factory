---
status: proposed
date: 2026-10-26
decision-makers: engineer
---

# Factory tooling as a Go CLI with thin adapters; skills copied into projects

## Context and Problem Statement

The factory consists of skills, checks, hooks and templates that projects use. Coding agents load these parts in different ways, and agent tools change faster than the projects that use them. The checks run in three places: in agent hooks on every tool call, on the engineer's machine, and in CI. How is the factory built and delivered, so that the agent tool stays replaceable and every project works in cloud sessions and CI?

## Decision Drivers

* Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read skills from `.agents/skills/`. Claude Code reads only `.claude/skills/`, and it loads a `<skill-name>` entry there that is a symlink.
* Hooks have no cross-tool standard.
* Harness changes alone can lower output quality: three Claude Code changes lowered it for weeks in 2026 while the model stayed the same (Anthropic postmortem, 2026-04-23).
* Cloud sessions and CI see only the project repository. Cloud sessions do not load plugins from local settings.
* Startup time, measured in a Linux container: a Go or Rust binary about 2 ms, a Node script about 30 ms, a Node script that imports Effect 3 about 390 ms. Hooks start a process on every tool call.
* Prebuilt binaries install without a language toolchain. Widely used CLIs are written in Go, for example `gh`, `kubectl` and Terraform.
* Each project pins the factory version it uses, and every update is reviewable.
* The project is also a learning project for Go.

## Considered Options

1. Go CLI `factory`, skills copied into projects
2. Node or TypeScript scripts copied into projects
3. Rust CLI, skills copied into projects
4. Claude Code plugin from a marketplace
5. Links from the project to a factory checkout (symlink or git submodule)

## Decision Outcome

Chosen option 1, because it gives the fastest hooks, needs no runtime in the project, and keeps tool-specific parts in one replaceable place.

* **CLI:** one Go binary `factory`. It uses only the Go standard library at runtime. Test-only dependencies are allowed. The policy logic lives in tool-neutral packages.
* **Skills:** embedded in the binary. `factory install` copies them into the project's `.agents/skills/` and pins the factory version in `.factory/config.json`. The project commits the copies.
* **Adapters:** `adapters/<tool>/` translates the tool's hook input and output and calls `factory` subcommands. The Claude Code adapter also creates relative symlinks `.claude/skills/<name>` → `../../.agents/skills/<name>`.
* **Distribution:** GitHub Releases with checksums. The engineer's machine installs the binary with an install script. Cloud sessions install it in their setup script, and CI installs it in a pinned download step.
* **Check:** `factory install --check` fails when copied files were edited in the project.

### Consequences

* Good, because hooks start in milliseconds.
* Good, because one tool serves the engineer, the hooks and CI.
* Good, because replacing the agent tool means adding one adapter.
* Bad, because releases need a pipeline, and every environment needs a pinned install step.
* Bad, because the skills exist twice, in the factory and in each project.
* Bad, because Windows clones need `core.symlinks`.
* Neutral, because project cloud threads do not load committed `.claude/skills/`. Other cloud sessions do.

### Confirmation

* A check fails if a file under `core/` or in the tool-neutral packages names a specific agent tool.
* `factory install --check` runs in the CI of every project.
* Each release publishes a checksum for every binary.
* An end-to-end test installs into a fixture project and checks that Claude Code lists the skills through the symlinks (increment 07).

## Pros and Cons of the Options

### Node or TypeScript scripts

* Good, because a TypeScript project already has Node, so nothing needs to be released.
* Bad, because hooks start more slowly, and much more slowly with Effect.

### Rust CLI

* Good, because it has the same startup time and stronger types.
* Bad, because Rust CLIs usually depend on crates, and Rust is not on the learning list.

### Claude Code plugin

* Good, because it installs with one command.
* Bad, because it works only in Claude Code, and cloud sessions do not load local plugins.

### Links to a factory checkout

* Good, because changes show up in the project at once.
* Bad, because the links break in cloud sessions and CI, and no version is pinned.

## More Information

* **Confidence:** medium. The symlink behaviour comes from the Claude Code docs. The release pipeline is not built yet (increment 03).
* **Revisit when** a cross-tool hook standard appears, Claude Code reads `.agents/skills/`, or keeping releases running costs more than it saves.
* **History:** first accepted with option 2 as the choice. The engineer reopened it on 2026-10-26, before anything was merged to main, and chose the Go CLI.
* **Sources:** agentskills.io specification; vercel-labs/skills agent paths; Claude Code docs on skills, memory and plugins; Anthropic postmortem 2026-04-23; Swamp source repository (a TypeScript CLI compiled to one binary). Details are in the [research notes](../research/2026-10-factory-basis.md).
