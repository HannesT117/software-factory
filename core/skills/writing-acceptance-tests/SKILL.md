---
name: writing-acceptance-tests
description: Writes one testscript file per automated scenario of an approved design note, and the shared runner, so that each test fails before the code exists. The agent uses it in a test-writer session when a design note on main has automated scenarios without tests, when a later note replaces a scenario, or when check-tests reports a missing or mismatched test.
---

# Writing acceptance tests

The acceptance tests are the evidence that an increment works. The person approves the scenarios in the design note, not the test files; a reviewer later checks that each test does what its scenario says. The tests are done when:

- every automated scenario in force has exactly one test file, and a replaced scenario has none,
- each test checks every case that its scenario names, and matches output only where the Interface of the note fixes it,
- `factory check-tests --base main` passes,
- and each new test fails because the code does not exist yet, while every other test passes. A new test that passes tests nothing new.

## Independence

The agent writes the tests from the design note on `main` and the public interface, never from the implementation, its unit tests or a branch with code. A test that follows the code proves only that the code does what it does. When a scenario is unclear, or the Interface leaves open what a test must match, the agent asks the person instead of guessing: a guess becomes the spec.

## Format

`assets/scenario.txtar` shows a test file. `assets/acceptance_test.go.tmpl` shows the shared runner: it registers the command under test, adds the command `exitcode` and sets `REPO` to the repo root. `factory check-tests` is the authority on file names, first lines and forbidden commands, and names each mistake. The `testscript` language is documented at https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript.

Every test runs everywhere, every time: no skips and no conditions. Given files stay as small as the case allows.

## Boundaries

The agent's commits change only the tests group and come first on the branch; `factory check-pr --base main` reports any other path. Once code commits follow, the tests do not change in the same pull request. The agent writes no code, and does not approve or merge.
