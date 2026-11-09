# Review of increment 04: design questions

Reviewed commit: b4255baafd0c3f411c1303448fecb2c18796a4bc
Pull request: https://github.com/HannesT117/software-factory/pull/59
Result: 5 findings

## Findings

- **Special cases:** none.
- **Test detection:** none.
- **Work outside the scenarios:** none; `docs/increments/04-design-questions/retro.md` is step 6 of `AGENTS.md`, the evidence files record scenarios 04.1 to 04.3, and the changed descriptions, the trigger query of `reviewing-increments` and `adapters/claude-code/agents/increment-reviewer.md` follow the goals of the note.
- **Contradicts the note**, `core/skills/designing-increments/SKILL.md:38`: the note decides in Existing tools that the questions come from "a fresh session that reads the whole note", because the scenario check reads only the scenarios and questions need the Interface and the Rollout too. The skill says only that the session gets the path of the note, right after the scenario check of line 30. In the run in `evidence/04-design-questions/review-runs.md`, the design agent told the fresh session to read only the Scenarios section, and the Interface and the Rollout only for the failure question. Without that sentence, a design agent can narrow the questions to the scenarios, which is the gap that the note's decision closes. Doubtful: one run, the session read the whole note anyway, and the Interface of the note also says only that the session gets the path.
- **Tests that miss their scenario**, `core/skills/designing-increments/evals/asks-comprehension-questions/prompt.md:8`: scenario 04.1 says that the fresh session reads only the note, and grader 01 checks that its prompt holds nothing else from the conversation. The prompt of the case holds no detail of a design conversation, so the grader would also pass a design agent that hands its conversation to the fresh session. `explains-gaps-after-answers` got two such details for this reason in b4255ba. Doubtful: the grader still fails a prompt that adds the agent's own notes, and cross-review raised the gap only for the other case.
- **Docs out of format**, `README.md:18`: the status of increment 04 says "design note approved; skills next", which no longer holds once this pull request merges; `AGENTS.md` asks for the status in every increment. The first comment on #59 deferred it until #50 merged; #50 is on `main` as 80fbbb7 and in this branch since df0982a. Doubtful: the branch does not change `README.md`, and the kind names changed docs.
- **Docs out of format**, `docs/increments/04-design-questions/retro.md:7`: the metrics count 19 eval runs, 5 fixes and the 2 findings of the first review. They leave out the round after cross-review: its 2 findings, the fixes in 00cf931 and b4255ba, and the 6 runs in `evidence/04-design-questions/evals.md`. Points for the factory leave out that the runner's judge does not see the subagent call, so the graders on it fail under the runner. Doubtful: the retro may change again with the manual scenarios.
- **Protected paths:** none; the branch changes no test file under `test/`, no CI or factory config and no design note. The notes under `core/skills/*/evals/*/files/` are eval fixtures.
- **Failing checks**, `docs/increments/03-setup-script/design.md` and `docs/increments/05-tests-per-increment/design.md:46`: the required job `checks` fails on b63b5c2 in `check-tests`, with 14 errors: the missing tests of scenarios 03.1 to 03.8 and 05.1 to 05.5, and the test of scenario 02.7, which the note of increment 05 replaces. So the pull request cannot merge without the bypass. Doubtful as a finding of this increment: `origin/main` alone gives the same 14 errors, and the branch adds no automated scenario.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 100 checked`
- `factory check-tests --base origin/main`: `check-tests: 14 errors`; the same on `origin/main` alone
- `factory hook stop` on this branch before this result: blocks and names b4255baafd0c3f411c1303448fecb2c18796a4bc
- CI on b63b5c2: `test` and `tests-start-red` pass; `checks` fails in `check-tests`; `cross-review` waits for a review result for the latest code commit

Each `factory` command ran with a binary built from the branch and with one built from `origin/main` (1b48461), with the same output.

## Demonstrated scenarios

- **scenario 04.1**: `evidence/04-design-questions/evals.md` and `evidence/04-design-questions/review-runs.md`. In `explains-gaps-after-answers`, the design agent starts a subagent whose prompt holds the note path, the three questions and the three answers, and neither `--loud` nor the greeting in German; it shows the three returned paragraphs unchanged, and then adds a paragraph of its own and a closing line. In `asks-comprehension-questions`, the shown questions are the three of the fresh session, word for word, one on the failure signal, none answered. The judge scores the runs 0.67 and 0.6; the graders on the subagent call fail under the judge although the traces hold the call. The old version did not run again; `evals.md` records it.
- **scenario 04.2**: `evidence/04-design-questions/evals.md`: `explains-deviations` and `clean-change-has-no-deviations`, old against new and after the fixes. `reviewing-increments` has not changed since the earlier round, so the review did not run them again.
- **scenario 04.3**: `evidence/04-design-questions/evals.md`: `sketches-changed-flow`, old against new and after the fixes; no sketch for a detail in grader 05 of `explains-deviations` and grader 04 of `clean-change-has-no-deviations`. Not run again, for the same reason.
