# Increment 05: tests per increment

## Goal

- `factory check-tests --base <ref>` asks for the missing tests of an increment only in the pull requests of that increment. While the tests of one increment are still being written, the code of another increment can merge.
- The output names each increment whose tests are still missing on `<ref>`, so that a reader sees what waits.
- Every other mistake in a test file still counts in every pull request.

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

## Interface

`factory check-tests [--base <ref>]`, as in increment 02, with this change.

**With `--base <ref>`.** The pull request is `HEAD` against its merge base with `<ref>`. It belongs to an increment `<nn>-<slug>` when it adds, changes or deletes a path under `test/acceptance/<nn>-<slug>/`, `evidence/<nn>-<slug>/` or `docs/increments/<nn>-<slug>/`. The error `missing-test` is reported only for the automated scenarios of increments that the pull request belongs to. For each other increment whose design note on `<ref>` has automated scenarios in force without a test file, the command prints `check-tests: waiting <nn>-<slug> <k> scenarios without tests`, where `<k>` is their number; such a line is not an error. The errors `unknown-id`, `title-mismatch`, `forbidden-command` and `replaced-test` are reported for every test file, as before. The last line and the exit codes are unchanged.

**Without `--base`.** Unchanged: every automated scenario in force without a test file is `missing-test`.

## Scenarios

- **scenario 05.1** (automated) **Another increment's code is not blocked.** On the base branch, the design note of increment 07 has two automated scenarios and no test files, and increment 06 has a note with all its test files. A pull request that changes a code file and adds `evidence/06-sample/review.md` passes `factory check-tests --base main`. The output names increment 07 as waiting, with two scenarios.
- **scenario 05.2** (automated) **An increment's own pull request still waits for its tests.** With the same base branch, a pull request that adds `evidence/07-later/review.md`, or a test file for only one of the two scenarios of increment 07, is rejected, with the missing test of each scenario that has none. A pull request that deletes a test file of increment 06 is rejected with that missing test.
- **scenario 05.3** (automated) **A pull request without increment folders is not blocked.** With the same base branch, a pull request that changes only `README.md` and a code file passes, and names increment 07 as waiting.
- **scenario 05.4** (automated) **Mistakes in test files count everywhere.** With the same base branch, a pull request of increment 06 that leaves a test file of increment 07 with a title that differs from its scenario is rejected for that file, although it does not belong to increment 07.
- **scenario 05.5** (automated) **Without a base, every missing test counts.** On the same repo, `factory check-tests` without `--base` rejects the two missing tests of increment 07.

## Replaced scenarios

- 02.7: a pull request is no longer rejected for the missing tests of an increment that it does not belong to. Scenarios 05.1 to 05.3 replace it.

## Evidence

- **Approved scenarios:** 05.1 to 05.5.
- **Not evidence:** unit tests from the implementing agent.

## Rollout

- **Off switch:** none needed; a later increment changes the rule again.
- **Failure signal:** the code of an increment reaches `main` before its tests; or a pull request is rejected for the missing tests of an increment that it does not touch. The `waiting` lines in the output of `checks` show which increments still wait for tests.
