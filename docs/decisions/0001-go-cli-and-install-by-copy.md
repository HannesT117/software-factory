---
status: proposed
date: 2026-10-26
decision-makers: engineer
---

# Factory tooling as a Go CLI with thin adapters; skills copied into projects

## Context and Problem Statement

The factory consists of skills, checks, hooks and templates that projects use. Coding agents load these parts in different ways, and agent tools change faster than the projects that use them. The checks run in three places: in agent hooks, on the engineer's machine, and in CI. How is the factory built and delivered, so that the agent tool stays replaceable and every project works in cloud sessions and CI?

## Decision Drivers

* Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read skills from `.agents/skills/`. Claude Code reads only `.claude/skills/`, and it loads a `<skill-name>` entry there that is a symlink. Cursor also loads `.claude/skills/`.
* Hooks have no cross-tool standard.
* Harness changes alone can lower output quality: three Claude Code changes lowered it for weeks in 2026 while the model stayed the same (Anthropic postmortem, 2026-04-23).
* Cloud sessions and CI see only the project repository. Cloud sessions do not load plugins from local settings.
* A compiled binary needs no language runtime in the project. A script needs its runtime in every project, cloud session and CI job.
* Startup time, measured in a Linux container: a Go binary about 2 ms, a Node script about 30 ms. Model latency is seconds, so this difference alone does not decide.
* Each project pins the factory version it uses, and every update is reviewable.
* The project is also a learning project for Go.

## Considered Options

1. Go CLI `factory`, skills copied into projects
2. Node or TypeScript scripts copied into projects
3. TypeScript compiled into one binary (for example `deno compile`, as Swamp does)
4. Rust CLI, skills copied into projects
5. Claude Code plugin from a marketplace
6. Links from the project to a factory checkout (symlink or git submodule)

## Decision Outcome

Chosen option 1, because it needs no runtime in the project, keeps tool-specific parts in one replaceable place, and serves the learning goal.

* **CLI:** one Go binary `factory`. It uses only the Go standard library at runtime, and it calls `git` and `ssh-keygen` for signatures. Test-only dependencies are allowed. The policy logic lives in tool-neutral packages.
* **Skills:** embedded in the binary. `factory install` copies them into the project's `.agents/skills/` and pins the factory version in `.factory/config.json`. The project commits the copies.
* **Adapters:** `adapters/<tool>/` translates the tool's hook input and output and calls `factory`. The Claude Code adapter also creates relative symlinks `.claude/skills/<name>` → `../../.agents/skills/<name>`.
* **Distribution:** GitHub Releases with binaries for macOS, Linux and Windows, built with `CGO_ENABLED=0 go build -trimpath`, with checksums and GitHub artifact attestations. An install step verifies the attestation before it installs, on the engineer's machine, in the cloud setup script, and in CI.
* **Finding the binary:** hooks call `factory` from `PATH`. `factory` refuses to run if its version differs from the pin in `.factory/config.json`.

### Consequences

* Good, because the project needs no runtime for the factory.
* Good, because one tool serves the engineer, the hooks and CI.
* Good, because replacing the agent tool means adding one adapter.
* Good, because attestations let every install check where a binary was built.
* Bad, because releases need a pipeline, and every environment needs an install step.
* Bad, because the skills exist twice, in the factory and in each project.
* Bad, because Windows clones need `core.symlinks`, and Cursor may list each skill twice.
* Neutral, because project cloud threads do not load committed `.claude/skills/`. Other cloud sessions do.

### Confirmation

* A check fails if a file under `core/` or in the tool-neutral packages names a specific agent tool.
* `factory install --check` runs in the CI of every project.
* Every release carries checksums and an attestation, and the install step refuses a binary whose attestation does not verify.
* An end-to-end test installs into a fixture project and checks that Claude Code lists the skills through the symlinks.

## Pros and Cons of the Options

### Node or TypeScript scripts

* Good, because a TypeScript project already has Node, so nothing needs to be released.
* Bad, because every project, cloud session and CI job needs Node.

### TypeScript compiled into one binary

* Good, because the language matches TypeScript projects such as `health`.
* Bad, because the binary embeds a full runtime, and it does not serve the learning goal.

### Rust CLI

* Good, because it has the same distribution model and stronger types.
* Bad, because Rust CLIs usually depend on crates, and Rust is not on the learning list.

### Claude Code plugin

* Good, because it installs with one command.
* Bad, because it works only in Claude Code, and cloud sessions do not load local plugins.

### Links to a factory checkout

* Good, because changes show up in the project at once.
* Bad, because the links break in cloud sessions and CI, and no version is pinned.

## More Information

* **Confidence:** medium. The symlink behaviour comes from the Claude Code docs. The release pipeline is built in increment 03.
* **Revisit when** a cross-tool hook standard appears, Claude Code reads `.agents/skills/`, or keeping releases running costs more than it saves.
* **History:** first accepted with option 2. The engineer reopened it on 2026-10-26, before anything was merged to main, and chose the Go CLI.
* **Sources:** agentskills.io specification; vercel-labs/skills agent paths; Cursor skills docs; Claude Code docs on skills, memory and plugins; Anthropic postmortem 2026-04-23; Swamp source repository; GitHub docs on artifact attestations; Go blog on reproducible builds. Details are in the [research notes](../research/2026-10-factory-basis.md).
