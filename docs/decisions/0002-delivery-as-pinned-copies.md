---
date: 2026-10-27
---

# Factory parts reach projects as committed copies with a pinned version

## Context and Problem Statement

Projects use the factory's skills, adapter files and templates, and the `factory` binary. Agent tools look for skills in different folders. Cloud sessions and CI see only the project repository. How do the factory's parts reach a project, so that every environment has them and every update can be reviewed?

## Decision Drivers

* Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read skills from `.agents/skills/`. Claude Code reads only `.claude/skills/`, and it loads a `<skill-name>` entry there that is a symlink. Cursor also loads `.claude/skills/`.
* Cloud sessions and CI see only the project repository. Cloud sessions do not load plugins from local settings.
* Harness changes alone can lower output quality: three Claude Code changes lowered it for weeks in 2026 while the model stayed the same (Anthropic postmortem, 2026-04-23). So a project must control when it takes a new factory version.
* Every update must be visible as a diff in a pull request.
* A binary from the internet must be traceable to the build that made it.

## Considered Options

1. Copy and commit: `factory install` writes the parts into the project, the project commits them and pins the version
2. Write at start without committing, as Swamp does: the parts go into the agents' skill folders, which are gitignored
3. Links from the project to a factory checkout (symlink or git submodule)
4. Claude Code plugin from a marketplace

## Decision Outcome

Chosen option 1, because the parts are then in the repo for every environment, and each update is a reviewable diff.

* **Skills and files:** embedded in the binary. `factory install` copies them into the project's `.agents/skills/` and pins the factory version in `.factory/config.json`. The project commits the copies.
* **Claude Code:** the adapter creates relative symlinks `.claude/skills/<name>` → `../../.agents/skills/<name>`.
* **Binary:** GitHub Releases with binaries for macOS, Linux and Windows, built with `CGO_ENABLED=0 go build -trimpath`, with checksums and GitHub artifact attestations. When a project changes its pinned version, `gh attestation verify` checks the attestation once, and the project pins the checksum. The install step on the engineer's machine, in the cloud setup script and in CI compares the checksum.
* **Finding the binary:** hooks call `factory` from `PATH`. `factory` refuses to run if its version differs from the pin in `.factory/config.json`. A build from source in the factory repo has the version `dev` and skips this check.

### Consequences

* Good, because every environment has the parts without access to the factory repo.
* Good, because each update arrives as a diff in a pull request.
* Good, because the attestation shows where a binary was built, and the pinned checksum keeps installs offline.
* Bad, because the skills exist twice, in the factory and in each project.
* Bad, because releases need a pipeline, and every environment needs an install step.
* Bad, because Windows clones need `core.symlinks`, and Cursor may list each skill twice.
* Bad, because project cloud threads do not load committed `.claude/skills/`. Other cloud sessions do.

### Confirmation

* `factory install --check` runs in the CI of every project and finds edited copies.
* Every release carries checksums and an attestation, and the install step refuses a binary whose checksum differs from the pin.
* An end-to-end test installs into a fixture project and checks that Claude Code lists the skills through the symlinks.

## Pros and Cons of the Options

### Copy and commit

* Good, because the parts are in the repo for every environment.
* Good, because an update is a normal diff.
* Bad, because a project can edit its copies. `factory install --check` finds that.

### Write at start without committing

* Good, because the project repo holds no copies.
* Bad, because every environment needs the binary before the agent starts.
* Bad, because an update does not show as a diff in the project.

### Links to a factory checkout

* Good, because changes show up in the project at once.
* Bad, because a symlink to a checkout outside the repo does not exist in cloud sessions and CI.
* Bad, because a submodule needs an extra fetch step in every environment and brings the whole factory repo into the project.

### Claude Code plugin

* Good, because it installs with one command.
* Bad, because it works only in Claude Code, and cloud sessions do not load local plugins.

## More Information

* **Confidence:** medium. The symlink behaviour comes from the Claude Code docs. The release pipeline is built in increment 03.
* **Revisit when** Claude Code reads `.agents/skills/`, or keeping releases running costs more than it saves.
* **Sources:** agentskills.io specification; vercel-labs/skills agent paths; Cursor skills docs; Claude Code docs on skills, memory and plugins; Anthropic postmortem 2026-04-23; Swamp source repository; GitHub docs on artifact attestations; Go blog on reproducible builds. Details are in the [research notes](../research/2026-10-factory-basis.md).
