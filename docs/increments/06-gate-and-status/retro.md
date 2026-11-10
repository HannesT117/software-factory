# Retro of increment 06

Increment 06 and the code of increment 05 came in one pull request, after a review of the factory against current sources (`docs/research/2026-11-factory-review.md`). The engineer asked for two pull requests, with the design to read and the implementation free.

## Metrics

- Pull requests: the design note (#68) and the tests and code of increments 05 and 06. The note of 05 changed in #67 along the way.
- Engineer decisions in the conversation: one, the task itself; the agent decided the rest and names it in the report.
- Sources checked: 8 sources, every claim of the review confirmed, four paraphrases corrected.
- Tests: 14 new test files from a separate session (05.7, 06.1 to 06.12, 06.15), two removed (02.26, 05.6), and the runner reports coverage. Two tests checked the output of `sh` instead of `factory` and went back to the test writer before the first code commit; one test (06.11) became stricter after the first code commit, on a finding of the review.
- Eval runs: 10 runs of `reviewing-increments`, old against new, on the session model, about 2.90 USD; 10 more without the Write tool were void.
- Coverage of the acceptance run: 85.4 percent of statements; 3 functions unreached.

## Review sources

| Source | Findings | Findings that changed code or a note | Cost |
|---|---|---|---|
| Test writer, while writing | 3 | 3: scenario 06.15 listed as unchanged behaviour, `-v` in the test job, the order of the notes | one session |
| Acceptance tests on the first code | 3 failing tests | 1 code fix (line number in an error), 2 test fixes | none |
| Fresh-context review, round 1 | 2 blocking, 3 notes | 5: test 06.11 checks the instructions, check-review reads every result, `/.claude/` owned, rule 27, `--no-renames` in the patch; also the sketch of the note | one subagent |
| Fresh-context review, round 2 | 1 blocking, 1 note | 1: the note says when cross-review comments; the blocking finding, the change of test 06.11 after the code, goes to the engineer | the same subagent |
| Other model (cross-review) | runs after the merge of the workflow patch | | |

## Agent version

- Claude Code 2.1.322; model `claude-opus-5-5` for the orchestrating session, the test writer and the eval runs; judge `haiku`.

## Corrections

- The token of the cloud session cannot push changes to `.github/workflows/`: it lacks the scope `workflow`. The workflow change goes to the engineer as a patch.
- The test writer first refused to write tests from notes that were not on `main`, as its definition says, until the engineer's words for this task reached it.
- `tests-start-red` read the list of unchanged behaviour from the note at the end of the test phase, which is the version that the branch merged in before its tests; it now reads the approved note on the base branch.
- `check-review` read only the first review result that named the commit; with two increments in one pull request, the blocking finding was in the second.
- An old test (02.26) sent code to the other model without a review result; scenario 06.15 replaces it.

## Points for the factory

- Applying the clock: rule 15 names `internal/report` as its check; rules 1, 2, 3, 9, 10, 18, 19, 22, 25, 37 and 38 are principles; rules 4, 26 and 41 already name a planned check.
- A new command that a workflow calls can only run after its pull request is on `main`, because the workflows build `factory` from `main`. A pull request that adds a check and the workflow step for it is red on that step until the merge.
- `claude plugin eval` withholds a tool that the run does not grant, so a `tool_used` grader on Write can never pass without `--allow-tools Write`; the runner warns, but the run still ends with a score.
- A testscript line `stdout` after `exec sh -c` reads the output of `sh`. A test that checks a command's output with a shell helper keeps the command's checks before the helper.

## Summary

The factory now shows the state of each increment with `factory status` and records manual results in a file the engineer owns. A blocking finding of the fresh-context review holds the merge through `factory check-review`, and after two rounds the engineer decides. The commands decide which check applies; the workflows only call them, once the engineer pushes the patch. The engineer owns the CLI, the adapter, the skills and the test runner. The other model reviews the tests before the code. The test run reports coverage. Open: the workflow patch, the ruleset (`factory setup`), the manual scenarios, and new evidence for the other skills.
