# Increment 05: tests per increment

## Goal

- `factory check-tests --base <ref>` asks for the missing tests of an increment only in the pull requests of that increment. While the tests of one increment are still being written, the code of another increment can merge.
- The output names each increment whose tests are still missing on `<ref>`, so that a reader sees what waits.
- Every other mistake in a test file still counts in every pull request.
- `factory cross-review --base <ref>` sends a code pull request to the other model with the design notes of the increments that it belongs to, by the same rule. A pull request that belongs to no increment, such as a change to the docs, goes without a design note, so that the other model does not measure it against the note of an increment that it does not touch.

## Out of scope

- When the workflow `checks` runs `check-tests`: that is CI config, changed in its own pull request.
- Tying code to an increment by its paths: the factory does not know which code belongs to which increment. A pull request belongs to an increment through the folders that name it.

## Domain

- Terms: none; the factory has no glossary yet.
- Contexts: none.
- Dependencies: none.

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Linking a pull request to its increment | path filters of GitHub Actions (`on.pull_request.paths`); `.github/CODEOWNERS`; the folders that the factory already names per increment | the folders per increment | Path filters decide whether a workflow runs, not which tests a check asks for. Every pull request with code already holds a review result in `evidence/<nn>-<slug>/`, because the stop hook asks for one. |
| The design note for cross-review | the note of the latest increment, as `factory cross-review` sends today when no increment folder changes; the other model choosing the note itself; the folders per increment | the folders per increment | The latest note makes the other model measure an unrelated pull request against it: a docs pull request got blocking findings about increment 04. A model that chooses the note gets every note and can choose wrong. The folders give the same answer as `check-tests`. |

## Interface

`factory check-tests [--base <ref>]`, as in increment 02, with this change.

**With `--base <ref>`.** The pull request is `HEAD` against its merge base with `<ref>`. It belongs to an increment `<nn>-<slug>` when it adds, changes or deletes a path under `test/acceptance/<nn>-<slug>/` or `evidence/<nn>-<slug>/`; it can belong to more than one. A change to a design note alone does not make it belong. The error `missing-test` is reported only for the automated scenarios of increments that the pull request belongs to. For each other increment whose design note on `<ref>` has automated scenarios in force without a test file, the command prints `check-tests: waiting <nn>-<slug> <k> scenarios without tests`, where `<k>` is their number; such a line is not an error. The errors `unknown-id`, `title-mismatch`, `forbidden-command` and `replaced-test` are reported for every test file in the working tree, as before, also when it is already on `<ref>`. The last line and the exit codes are unchanged.

**Without `--base`.** Unchanged: every automated scenario in force without a test file is `missing-test`, and no `waiting` line is printed.

`factory cross-review --base <ref>`, as in increment 02, with this change for a pull request that changes no design note. The request holds the design notes of the increments that the pull request belongs to, by the rule above, the titles of their scenarios in force, and the diff. For a pull request that belongs to no increment, it holds the diff, with no design note and no scenario title, and its instructions ask the other model to judge the change by itself. The model, the output lines and the exit codes are unchanged. A pull request that changes a design note is sent as before.

## Scenarios

In each scenario, the base branch `main` holds the design note of increment 06, with all its test files, and the design note of increment 07, with two automated scenarios. A pull request is a branch from `main`, and `factory check-tests --base main` runs on it; in scenario 05.6, `factory cross-review --base main` runs on it, with a stand-in for the other model that keeps each request. A pull request belongs to an increment when it changes a file in that increment's folder of test files or of evidence, such as `evidence/07-later/`.

- **scenario 05.1** (automated) **Another increment's code is not blocked.** Increment 07 has no test files. A pull request that changes a code file and adds `evidence/06-sample/review.md` belongs only to increment 06 and passes. The output names increment 07 as waiting, with two scenarios.
- **scenario 05.2** (automated) **An increment's own pull request still waits for its tests.** Increment 07 has no test files. Each of these pull requests is rejected, with one missing test for each scenario of 07 without a test file: one that adds `evidence/07-later/review.md`; one that adds the test file of only the first scenario of 07, which leaves the second missing; and one that adds `evidence/06-sample/review.md` and `evidence/07-later/review.md`. A pull request that deletes a test file of increment 06 is rejected with that missing test.
- **scenario 05.3** (automated) **A pull request without increment folders is not blocked.** Increment 07 has no test files. A pull request that changes only `README.md` and a code file belongs to no increment and passes, and names increment 07 as waiting. So does a pull request that changes only the design note of 07.
- **scenario 05.4** (automated) **Mistakes in test files count everywhere.** On `main`, increment 07 has a test file for each scenario, but the first line of one of them names a title that differs from its scenario. A pull request that adds `evidence/06-sample/review.md` belongs only to increment 06, and is rejected for that file of increment 07.
- **scenario 05.5** (automated) **Without a base, every missing test counts.** Increment 07 has no test files. On the branch of scenario 05.1, `factory check-tests` without `--base` rejects both missing tests of increment 07 and prints no line about waiting increments.
- **scenario 05.6** (automated) **The other model gets only the notes of the pull request's own increments.** Increment 07 has no test files. A pull request that changes a code file and adds `evidence/06-sample/review.md` is sent with its diff and the design note of 06, and with neither the note of 07 nor the title of a scenario of 07. One that adds `evidence/06-sample/review.md` and `evidence/07-later/review.md` is sent with both notes. One that changes only `README.md` belongs to no increment and is sent with its diff, but with neither note and no scenario title.

## Replaced scenarios

- 02.7: a pull request is no longer rejected for the missing tests of an increment that it does not belong to. Scenarios 05.1 to 05.3 replace it.

## Evidence

- **Approved scenarios:** 05.1 to 05.6.
- **Not evidence:** unit tests from the implementing agent.

## Rollout

- **Off switch:** none needed; a later increment changes the rule again.
- **Failure signal:** the code of an increment reaches `main` before its tests; or a pull request is rejected for the missing tests of an increment that it does not touch. The `waiting` lines in the output of `checks` show which increments still wait for tests. For cross-review: its comment on a pull request names findings about an increment that the pull request does not touch.
