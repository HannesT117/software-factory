# Increment 03: install and distribution

Status: proposed

This is a draft. It is rewritten with the `designing-increments` skill after increment 02. Only then does it go to the engineer.

## Goal

- `factory install` brings the skills, the `AGENTS.md` block and the adapter files into a project, and updates them.
- `factory install --check` finds edits to copied files, also in CI.
- Releases publish binaries for macOS, Linux and Windows, each with a checksum and a GitHub artifact attestation. The install step verifies the attestation on the engineer's machine, in cloud sessions and in CI.
- A ruleset template protects the default branch: pull requests only, squash merges only, required checks, empty bypass list.
- After this increment, `health` starts.

## Out of scope

- the provenance log, the workflow gate, holdout scenarios, chaos and more skills: later, when `health` needs them
- adapters other than Claude Code

## Interface

To be written with `designing-increments`.

## Scenarios

- **scenario 03.1** (automated) **Fresh install.** Given an empty project, when `factory install` runs, then the project has the skills under `.agents/skills/`, the links under `.claude/skills/`, `.factory/config.json` with the pinned version, the `AGENTS.md` block and `CLAUDE.md` with `@AGENTS.md`.
- **scenario 03.2** (automated) **Second run.** Given an installed project, when `factory install` runs again, then no file changes.
- **scenario 03.3** (automated) **Dry run.** Given an empty project, when `factory install --dry-run` runs, then the project stays empty, and the output lists every planned change.
- **scenario 03.4** (automated) **Check finds edits.** Given an installed project with an edited skill file, when `factory install --check` runs, then it exits with 1 and names the file.
- **scenario 03.5** (automated) **No overwrite.** Given an edited or foreign file at a path where the factory writes, when `factory install` runs, then it exits with 1, names the file, and changes nothing.
- **scenario 03.6** (automated) **Version pin.** Given a project that pins another version, when `factory` runs, then it exits with 2 and names both versions.
- **scenario 03.7** (manual) **Release with attestation.** A tag builds the binaries, and `gh attestation verify` passes for each of them.
- **scenario 03.8** (manual) **Install everywhere.** The install step works on the Mac, in the cloud setup script and in CI, and it refuses a binary without a valid attestation.
- **scenario 03.9** (manual) **The ruleset blocks.** In the public repo with the imported ruleset, a pull request with a failing check cannot be merged, also not by the owner.

## Evidence

To be written with `designing-increments`. This increment also collects the first ratings for scenario 02.17.

## Open questions

To be written with `designing-increments`.
