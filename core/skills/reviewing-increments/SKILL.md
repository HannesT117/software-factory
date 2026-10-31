---
name: reviewing-increments
description: Reviews the pull request of an increment against its design note and writes the review result. Use in a fresh context, after the implementing agent's last code commit and before the pull request merges.
---

# Reviewing increments

The reviewer runs in a fresh context. It has not seen how the code came about, so it reads the change as the engineer will. Its result decides whether the pull request may merge without the engineer: a review with no findings allows auto-merge, any finding makes the pull request wait.

## What to read

- The design note on `main`: `docs/increments/<nn>-<slug>/design.md`.
- The diff of the pull request against `main`: `git diff main...HEAD`, and `git log main..HEAD` for the order of the commits.
- The test files under `test/acceptance/<nn>-<slug>/`.
- `docs/principles.md` and `AGENTS.md` for the rules that the change must keep.

## Checks

Report each of these as a finding, with the file and line:

1. **Special cases:** code that treats one input differently from the rule in the note, or that hard-codes a test case.
2. **Test detection:** code that behaves differently when it runs under a test, for example by reading test paths, environment variables or file names of the tests.
3. **Work outside the scenarios:** code or docs that no scenario and no goal of the note asks for, and changes whose benefit the pull request does not name.
4. **Tests that miss their scenario:** a test that does not check what its scenario says, checks less than the scenario names, or would pass without the code.
5. **Docs:** a changed doc that tells its own history instead of the current state, or that uses elements outside its named format (for example MADR or the design-note format).
6. **Protected paths:** a commit after the test phase that changes a test file, `go.mod`, `go.sum`, CI or factory config, or a design note.

Run the checks of the factory on the head commit: `go test ./...`, `factory check-pr --base main`, `factory check-tests --base main`, `factory check-adr`, and `factory check-design` for each changed note. A failing check is a finding.

Run each demonstrated scenario of the note, record it under `evidence/<nn>-<slug>/`, and name the recording in the result.

## Comprehension questions

Ask 2 or 3 questions about the change that the engineer should be able to answer after reading it. At least one is about debugging: where to look, or what to run, when a part of the change fails. Do not answer them in the result. The engineer answers first; the reviewer explains after that.

## Result

Write `evidence/<nn>-<slug>/review.md` from `assets/review.md`. It names the full hash of the reviewed commit, which is the latest commit with code. Commit and push the result. Then:

- With no findings: turn on auto-merge for the pull request.
- With findings: leave auto-merge off. The pull request waits for the engineer.

Link the result in the pull request. Do not approve or merge the pull request, and do not change the code under review.
