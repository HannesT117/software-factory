# Increment 03: install and distribution

This is a draft. It is rewritten with the `designing-increments` skill after increment 02. Only then does it go to the engineer.

## Goal

- `factory install` brings the skills, the `AGENTS.md` block, the adapter files and the pull request template into a project, and updates them.
- `factory install --check` finds edits to copied files, also in CI.
- Releases publish binaries for macOS, Linux and Windows, each with a checksum and a GitHub artifact attestation.
- A release starts only by `workflow_dispatch`, which the agent token cannot trigger. The release workflow checks that the commit is on `main`. Releases are immutable.
- When a project changes its pinned version, the engineer verifies the attestation once with `gh attestation verify --signer-workflow`, and the project pins the checksum. The install step on the Mac, in cloud sessions and in CI compares only the checksum, so it needs no network access to Sigstore.
- The ruleset and the approval gate from increment 02 go into projects through `docs/setup/`.

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
- **scenario 03.7** (manual) **Release with attestation.** A `workflow_dispatch` run on `main` builds the binaries, and `gh attestation verify --signer-workflow` passes for each of them. A run on a commit that is not on `main` fails.
- **scenario 03.8** (manual) **Install everywhere.** The install step works on the Mac, in the cloud setup script and in CI, and it refuses a binary whose checksum differs from the pinned one.
- **scenario 03.9** (manual) **The gate works in a project.** In `health`, after the setup steps, a pull request cannot be merged without the engineer's approval, also not by the owner.

## Evidence

To be written with `designing-increments`. This increment also collects the first counts for scenario 02.16.

## Open questions

To be written with `designing-increments`.
