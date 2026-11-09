# Plan: increment 05, tests per increment

- **Folder:** `docs/increments/05-tests-per-increment/`
- **Target pull requests:**
  1. Design pull request #52, from `claude/05-tests-per-increment-design`.
  2. Tests and code on one branch, `claude/05-tests-per-increment`, after the design note is on `main`.
- **Runs as:** the orchestrating session for the design; then a test-writer session and an implementing session, as `AGENTS.md` describes.

## Context

- `factory check-tests --base main` rejects every pull request while any design note on `main` has automated scenarios without test files (scenario 02.7). The tests of an increment come in the same pull request as its code, so while one increment is in work, the code of every other increment fails the check.
- The workflow fix #51 already keeps docs-only pull requests from running the check.

## Decisions of the engineer, 2026-11-03

- Increments must not block each other: option B, narrow the check in a small increment. It comes before install and distribution (06).

## Steps

1. Design note, fresh-context check, design pull request: done (#52).
2. After the merge: a test-writer session writes 05.1 to 05.5 and removes the test file of the replaced scenario 02.7.
3. An implementing session changes `internal/testcheck` and `cmd/factory` until the tests pass; review; retro.

## Checkpoints
- 2026-11-03T02:10Z step 1 done; claude/05-tests-per-increment-design pushed; next: step 2 after the engineer approves #52; open: none
- 2026-11-09T00:14Z design merged as #60; next: test writer for 05.1 to 05.5 after 03 and 04 merge; open: none
- 2026-11-09T00:41Z engineer's decision: the fix of cross-review for pull requests without an increment belongs to 05; design pull request #64 (https://github.com/HannesT117/software-factory/pull/64) adds it to the note with scenario 05.6; claude/05-cross-review-per-increment@525390d; comprehension questions shown; next: after #64 merges, the test writer writes 05.1 to 05.6 and removes the test of 02.7; open: none
- 2026-11-09T01:49Z note of cross-review on #64 fixed (form and order of waiting lines); claude/05-cross-review-per-increment@02e3a3b; the engineer answered the comprehension questions and got the explanation; next: the engineer merges #64, then the test writer; open: none
