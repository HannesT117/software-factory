---
name: writing-acceptance-tests
description: Writes one testscript file per automated scenario of an approved design note, and the runner. Use in the test-writer session, after the design note is on main and before any code exists.
---

# Writing acceptance tests

The tests are the evidence that an increment works. They come from a session of their own, so that they do not follow the implementation. The engineer approves the scenarios in the design note, not the test files. A reviewer later checks that each test does what its scenario says.

## What to read

- The design note on `main`: `docs/increments/<nn>-<slug>/design.md`. Its `## Scenarios` say what to test; its `## Interface` gives the commands, output, error codes and exit codes.
- The public interface: the usage text of the command and the docs that describe it.
- `assets/scenario.txtar` and `assets/acceptance_test.go.tmpl` in this skill.

Do not read the implementation, its unit tests, or a branch with code. When a scenario is unclear, or the Interface does not say what a test needs, stop and ask the engineer. Do not guess.

## Steps

1. Create the branch of the increment from `main`.
2. Write one file per automated scenario: `test/acceptance/<nn>-<slug>/scenario-<nn>.<n>.txtar`. Start from `assets/scenario.txtar`.
3. Write the runner `test/acceptance/acceptance_test.go` from `assets/acceptance_test.go.tmpl`, if the project does not have it yet. Replace `./cmd/factory` with the main package of the command under test. Add the test-only dependency to `go.mod` and `go.sum` with `go get github.com/rogpeppe/go-internal/testscript`.
4. Check: `factory check-tests --base main` must pass.
5. Run the tests. Each new test must fail, because the code does not exist yet, and every other test must pass. A new test that passes tests nothing new: change it until it fails for the right reason.
6. Commit only paths in the tests group: `test/acceptance/**`, `go.mod`, `go.sum`. Push. Stop. The implementing agent continues on the same branch.

## Format of a test file

- Line 1 is `# scenario <nn>.<n>: <title>`, with the title exactly as in the note: same case, without the final period or question mark.
- Not allowed: `skip`, `stop`, and any condition in square brackets, such as `[net]` or `[darwin]`. Each test runs everywhere, every time.
- Test each case that the scenario names, with given files as small as the case allows.
- Match output by the patterns that the Interface fixes, for example `^check-pr: ok \d+ checked$`. Do not match wording that the Interface leaves open.

How to read one: `exec` runs a command that must succeed. `! exec` runs one that must fail. `exitcode N <command> …` runs the command and checks the exact exit code. `stdout '…'` and `stderr '…'` check for a line that matches the pattern. `exists path` checks that a file exists. Each `-- path --` starts a given file.

## Format of the runner

One file runs all test files. It registers the command under test, adds the command `exitcode` (because `testscript` checks only success or failure, not the exit code), sets `REPO` to the repo root (because each file runs in its own empty folder), and sets `HOME`, `GIT_CONFIG_NOSYSTEM=1` and a fresh temporary folder per file.
