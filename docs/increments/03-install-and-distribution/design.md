# Increment 03: install and distribution

Status: proposed

This is a draft. It is rewritten with the `designing-increments` skill after increment 02, and only then goes to the engineer.

## Goal

- `factory install` brings the skills, the `AGENTS.md` block and the adapter files into a project, and updates them.
- `factory install --check` finds edits to copied files, also in CI.
- Releases publish binaries with checksums. The engineer's machine, cloud sessions and CI install a pinned version.
- A ruleset template protects the default branch.

## Out of scope

- the provenance log and the workflow gate: increment 04
- adapters other than Claude Code
- `factory install --watch`

## Interface

To be written with `designing-increments`.

## Scenarios

- **scenario 03.1** (automated) **Fresh install.** Given an empty project, when `factory install` runs, then the project has the skills under `.agents/skills/`, the links under `.claude/skills/`, `.factory/config.json` with the pinned version, the `AGENTS.md` block and `CLAUDE.md` with `@AGENTS.md`.
- **scenario 03.2** (automated) **Second run.** Given an installed project, when `factory install` runs again, then no file changes.
- **scenario 03.3** (automated) **Dry run.** Given an empty project, when `factory install --dry-run` runs, then the project stays empty, and the output lists every planned change.
- **scenario 03.4** (automated) **Check finds edits.** Given an installed project with an edited skill file, when `factory install --check` runs, then it exits with 1 and names the file.
- **scenario 03.5** (automated) **No overwrite.** Given an edited or foreign file at a path where the factory writes, when `factory install` runs, then it exits with 1, names the file, and changes nothing.
- **scenario 03.6** (automated) **Existing AGENTS.md and CLAUDE.md.** Given a project with its own files, when `factory install` runs twice, then the own text stays, the block appears once, and a `CLAUDE.md` without `@AGENTS.md` only gets a warning.
- **scenario 03.7** (automated) **Removed skill.** Given a factory without one installed skill, when `factory install` runs, then that skill and its link are gone.
- **scenario 03.8** (manual) **Release.** A tag builds binaries for macOS, Linux and Windows with checksums.
- **scenario 03.9** (manual) **Install everywhere.** The install script works on the Mac, the cloud setup script installs the pinned version, and the CI step installs it.
- **scenario 03.10** (manual) **The ruleset blocks.** In the public repo with the imported ruleset, a pull request with a failing check cannot be merged, also not by the owner.

## Evidence

To be written with `designing-increments`. This increment also holds the rating of the skills in real use (scenario 02.16).

## Open questions

To be written with `designing-increments`.
