# Increment 05: tests per increment

## Goal

- `factory check-tests --base <ref>` asks for the missing tests of an increment only in the pull requests of that increment. While the tests of one increment are still being written, the code of another increment can merge.
- The output names each increment whose tests are still missing on `<ref>`, so that a reader sees what waits.
- Every other mistake in a test file still counts in every pull request.
- `factory cross-review --base <ref>` sends a code pull request to the other model with the design notes of the increments that it belongs to, by the same rule. A pull request that belongs to no increment, such as a change to the docs, goes without a design note, so that the other model does not measure it against the note of an increment that it does not touch.
- A design note can list the scenarios whose behaviour the code before the increment already has, such as a scenario that fixes what a narrower rule keeps. The job `tests-start-red` expects the tests of these scenarios to pass before the code, and every other added test to fail. So such a scenario needs no bypass of the merge.

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
| Tests of unchanged behaviour | a bypass of the merge for each pull request with such a test; no scenario for unchanged behaviour; unchanged behaviour as a step of a scenario with new behaviour; a list in the note, like `## Replaced scenarios` | a list in the note | A bypass is needed again in each increment that makes a rule narrower or wider. Without such scenarios, code that goes too far passes, such as code that also skips the mistakes in the test files of other increments. As a step of a larger scenario, a failure names only the larger scenario. The list shows the engineer at approval which tests start green, and `check-design` on `main` already allows other headings, so a note can use the list before its check exists. |
| The design note for cross-review | the note of the latest increment, as `factory cross-review` sends today when no increment folder changes; the other model choosing the note itself; the folders per increment | the folders per increment | The latest note makes the other model measure an unrelated pull request against it: a docs pull request got blocking findings about increment 04. A model that chooses the note gets every note and can choose wrong. The folders give the same answer as `check-tests`. |

## Interface

`factory check-tests [--base <ref>]`, as in increment 02, with this change.

**With `--base <ref>`.** The pull request is `HEAD` against its merge base with `<ref>`. It belongs to an increment `<nn>-<slug>` when it adds, changes or deletes a path under `test/acceptance/<nn>-<slug>/` or `evidence/<nn>-<slug>/`; it can belong to more than one. A change to a design note alone does not make it belong. The error `missing-test` is reported only for the automated scenarios of increments that the pull request belongs to. For each other increment whose design note on `<ref>` has automated scenarios in force without a test file, the command prints `check-tests: waiting <nn>-<slug> <k> scenarios without tests`, where `<k>` is their number and the word stays `scenarios` also for 1; the lines come in the order of the increment numbers, and such a line is not an error. The errors `unknown-id`, `title-mismatch`, `forbidden-command` and `replaced-test` are reported for every test file in the working tree, as before, also when it is already on `<ref>`. The last line and the exit codes are unchanged.

**Without `--base`.** Unchanged: every automated scenario in force without a test file is `missing-test`, and no `waiting` line is printed.

`factory cross-review --base <ref>`, as in increment 02, with this change for a pull request that changes no design note. The request holds the design notes of the increments that the pull request belongs to, by the rule above, the titles of their scenarios in force, and the diff. For a pull request that belongs to no increment, it holds the diff, with no design note and no scenario title, and its instructions ask the other model to judge the change by itself. The model, the output lines and the exit codes are unchanged. A pull request that changes a design note is sent as before.

**Design note**, as in increment 02, with an optional heading `## Unchanged behaviour`. Each line under it matches `^- (\d{2})\.([1-9]\d*): (.+)$`: the id of an automated scenario of the same note, and why the code before this increment already behaves so. `factory check-design` reports the error `bad-unchanged-line` for a line in another form, and `unknown-scenario` for an id that is not an automated scenario of the note.

**Job `tests-start-red`** of the workflow `test`, as in increment 02, with this change: on the commit where the test phase ends, the added test of each scenario listed under `## Unchanged behaviour` in its note at that commit passes, each other added test fails, and every other test passes. The output lines are unchanged: `tests-start-red: error <test>: <result>, want <result>` for each test with another result, then `tests-start-red: ok <n> checked` or `tests-start-red: <k> errors`. The change is in a pull request of the group CI and config, which merges before the tests of this increment.

## Scenarios

In scenarios 05.1 to 05.6, the base branch `main` holds the design note of increment 06, with all its test files, and the design note of increment 07, with two automated scenarios. A pull request is a branch from `main`, and `factory check-tests --base main` runs on it; in scenario 05.6, `factory cross-review --base main` runs on it, with a stand-in for the other model that keeps each request. A pull request belongs to an increment when it changes a file in that increment's folder of test files or of evidence, such as `evidence/07-later/`.

- **scenario 05.1** (automated) **Another increment's code is not blocked.** Increment 07 has no test files. A pull request that changes a code file and adds `evidence/06-sample/review.md` belongs only to increment 06 and passes. The output names increment 07 as waiting, with two scenarios.
- **scenario 05.2** (automated) **An increment's own pull request still waits for its tests.** Increment 07 has no test files. Each of these pull requests is rejected, with one missing test for each scenario of 07 without a test file: one that adds `evidence/07-later/review.md`; one that adds the test file of only the first scenario of 07, which leaves the second missing; and one that adds `evidence/06-sample/review.md` and `evidence/07-later/review.md`. A pull request that deletes a test file of increment 06 is rejected with that missing test.
- **scenario 05.3** (automated) **A pull request without increment folders is not blocked.** Increment 07 has no test files. A pull request that changes only `README.md` and a code file belongs to no increment and passes, and names increment 07 as waiting. So does a pull request that changes only the design note of 07.
- **scenario 05.4** (automated) **Mistakes in test files count everywhere.** On `main`, increment 07 has a test file for each scenario, but the first line of one of them names a title that differs from its scenario. A pull request that adds `evidence/06-sample/review.md` belongs only to increment 06, and is rejected for that file of increment 07.
- **scenario 05.5** (automated) **Without a base, every missing test counts.** Increment 07 has no test files. On the branch of scenario 05.1, `factory check-tests` without `--base` rejects both missing tests of increment 07 and prints no line about waiting increments.
- **scenario 05.6** (automated) **The other model gets only the notes of the pull request's own increments.** Increment 07 has no test files. A pull request that changes a code file and adds `evidence/06-sample/review.md` is sent with its diff and the design note of 06, and with neither the note of 07 nor the title of a scenario of 07. One that adds `evidence/06-sample/review.md` and `evidence/07-later/review.md` is sent with both notes. One that changes only `README.md` belongs to no increment and is sent with its diff, but with neither note and no scenario title.
- **scenario 05.7** (automated) **A note can list the scenarios whose behaviour stays.** A design note of increment 07 with an automated and a manual scenario lists the automated one as unchanged behaviour, with a reason, and `factory check-design` accepts the note. The check rejects the note when the list names the manual scenario, a scenario that the note does not have, or a scenario without a reason.
- **scenario 05.8** (demonstrated) **A test of unchanged behaviour may start green.** On a pull request that changes test files, `tests-start-red` passes only when, on the last commit before the code, the added test of each scenario listed as unchanged behaviour passed, every other added test failed, and every other test passed. It fails when the test of a listed scenario failed there.

## Replaced scenarios

- 02.7: a pull request is no longer rejected for the missing tests of an increment that it does not belong to. Scenarios 05.1 to 05.3 replace it.
- 02.20: the added test of a scenario listed as unchanged behaviour passes before the code. Scenario 05.8 replaces it.

## Unchanged behaviour

- 05.4: `check-tests` already reports a wrong title in every test file of the working tree.
- 05.5: without `--base`, `check-tests` already reports every missing test, and it prints no line about waiting increments.

## Evidence

- **Approved scenarios:** 05.1 to 05.7.
- **Demonstrated:** 05.8, by the review session: the output of `tests-start-red` on the pull request of this increment, where the tests of 05.4 and 05.5 pass before the code, and a run of the same check on those results with the test of 05.4 set to failed.
- **Not evidence:** unit tests from the implementing agent.

## Rollout

- **Off switch:** none needed; a later increment changes the rule again. A note without `## Unchanged behaviour` keeps the rule of increment 02 for `tests-start-red`.
- **Order:** the pull request that changes `tests-start-red` merges before the tests of this increment; until then, the tests of 05.4 and 05.5 get `pass, want fail`.
- **Failure signal:** the code of an increment reaches `main` before its tests; or a pull request is rejected for the missing tests of an increment that it does not touch. The `waiting` lines in the output of `checks` show which increments still wait for tests. For cross-review: its comment on a pull request names findings about an increment that the pull request does not touch. For unchanged behaviour: `tests-start-red` reports `fail, want pass` for a listed test, so the note or the test is wrong.
