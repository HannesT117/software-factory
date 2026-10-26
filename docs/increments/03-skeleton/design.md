# Increment 03: skeleton

Status: draft. It is rewritten with the `designing-increments` skill after increment 02 (scenario 02.15), and only then goes to the engineer for approval.

## Goal

- A project installs the factory and updates it with one command.
- The project detects edits to factory files, also in CI without a factory checkout.
- The factory repo has CI, and a ruleset template protects the default branch.

## Out of scope

- the skill check, the design note check and the test lock: increment 02
- the full lock, receipts, workflow gate, private-paths check: increment 04
- adapters other than Claude Code
- `install.mjs --watch` and `--no-symlink`

## Scenarios

### Install

- **scenario 03.1** (automated) **Fresh install.** Given an empty project folder, when `install.mjs` runs, then the exit code is 0, and the project contains:
  - each factory skill under `.agents/skills/<name>/`
  - a link `.claude/skills/<name>` → `../../.agents/skills/<name>`
  - the factory scripts under `.factory/scripts/`
  - `.factory/config.json` with empty `lockedPaths` and `privatePaths`
  - `AGENTS.md` with the factory block
  - `CLAUDE.md` with the line `@AGENTS.md`
  - `.factory/manifest.json`
- **scenario 03.2** (automated) **Second run.** Given an installed project, when `install.mjs` runs again with the same factory, then no file changes, and the output reports 0 changes.
- **scenario 03.3** (automated) **Dry run.** Given an empty project, when `install.mjs --dry-run` runs, then the project stays empty, and the output lists every planned change.
- **scenario 03.4** (automated) **Check passes.** Given an installed project, when `install.mjs --check` runs, then the exit code is 0.
- **scenario 03.5** (automated) **Check finds an edited file.** Given an installed project with an edited skill file, when `install.mjs --check` runs, and when `.factory/scripts/check-install.mjs` runs in the project without the factory, then both exit with 1 and name the file.
- **scenario 03.6** (automated) **Check finds a missing link and an edited block.** Given an installed project without one skill link and with an edited factory block, when `install.mjs --check` runs, then it exits with 1 and names both.
- **scenario 03.7** (automated) **No overwrite of local edits.** Given an installed project with an edited skill file, and a factory with a newer version of that file, when `install.mjs` runs, then it exits with 1, names the file, and no file in the project changes.
- **scenario 03.8** (automated) **No overwrite of foreign files.** Given a project with its own file at a path where the factory writes, when `install.mjs` runs, then it exits with 1, names the file, and the file keeps its content.
- **scenario 03.9** (automated) **Existing AGENTS.md.** Given a project with its own `AGENTS.md`, when `install.mjs` runs twice, then the project's own text stays at the top, and the factory block appears exactly once.
- **scenario 03.10** (automated) **Existing CLAUDE.md.** Given a project with its own `CLAUDE.md` without `@AGENTS.md`, when `install.mjs` runs, then `CLAUDE.md` stays unchanged, and the output warns about the missing line.
- **scenario 03.11** (automated) **Skill removed from the factory.** Given an installed project and a factory without one of its skills, when `install.mjs` runs, then that skill's folder and link are gone, the other skills stay, and `--check` passes.
- **scenario 03.12** (automated) **No adapter.** Given an empty project, when `install.mjs --adapter none` runs, then there is no `.claude/` folder and no `CLAUDE.md`, and the skills are installed.
- **scenario 03.13** (automated) **Unknown adapter.** Given an empty project, when `install.mjs --adapter nope` runs, then it exits with 2 and names the available adapters.

### CI and ruleset

- **scenario 03.14** (manual) **CI runs.** Given a pull request in the factory repo, then GitHub shows the check `factory-checks`, and it runs the acceptance tests and all checks.
- **scenario 03.15** (manual) **The ruleset blocks.** Given the public repo with the imported ruleset, when a pull request has a failing `factory-checks`, then the merge button stays disabled, also for the repo owner.

## Evidence

- **Approved scenarios:** 03.1 to 03.13, as tests that a separate agent wrote and the engineer reviewed.
- **Manual:** 03.14 and 03.15.
- **ADR confirmation:** no file under `core/` names an agent tool.
- **Not evidence:** any unit tests the implementing agent writes.

## Open questions

1. **Exit codes:** 2 for wrong usage (unknown option or adapter), 1 for a failed check or a refused install. Is that right?
2. **Scripts in the project:** the install copies every script from `core/scripts/` into `.factory/scripts/`. Is that right?
