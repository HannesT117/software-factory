# Review of increment 04: design questions

Reviewed commit: 4bd0a2d900cdab502cc8482b5bc18c66ec39cb44
Pull request: https://github.com/HannesT117/software-factory/pull/59
Result: 3 findings

## Findings

- **Special cases:** none; neither skill names a case of the evals (`greet`, `outbox`, `--loud`, Birgit).
- **Test detection:** none.
- **Work outside the scenarios:** none; 86e761d serves the decision of the note in Existing tools that the fresh session reads the whole note, 4bd0a2d serves scenario 04.1, and b6bd67e follows the status and retro rules of `AGENTS.md`.
- **Tests that miss their scenario**, `core/skills/designing-increments/evals/asks-comprehension-questions/graders/01-a-fresh-context-reads-only-the-note.md:7`: the grader checks the prompt of the fresh session and its read of the note, but the judge of `claude plugin eval` gets a shortened trace that holds neither the Agent call nor the events of the subagent. In the last round, the judge passed this grader in 3 of 4 runs and failed it once, without the call in its evidence. The same holds for grader 02 of this case, which compares the shown questions with those of the subagent, and for grader 01 of `explains-gaps-after-answers`. So the case as committed can pass a design agent that hands `--loud` or Birgit to the session, or limits it to sections, which is what 86e761d and 4bd0a2d guard against. The scores in `evidence/04-design-questions/evals.md` count a script on the full traces; that script is neither in the repo nor in the evidence, so a later run, such as on a model change (rule 22), cannot repeat the check. Doubtful: the limit lies in the runner, not in the case, and `evals.md` names the script check, whose result the kept traces confirm.
- **Docs out of format**, `docs/increments/04-design-questions/retro.md:10`: the retro says that the whole-note read and the detail of a design conversation in the prompt of asks-comprehension-questions are not fixed; 86e761d and 4bd0a2d fixed both. Line 9 counts 2 eval runs in the third review and not the 4 recorded and 4 dropped runs after it (`evidence/04-design-questions/evals.md`, After the third review). Line 14 names Claude Code 2.1.317, while the runs since the cross-review round ran on 2.1.319, which rule 21 asks to record. Line 26 says that the judge fails every grader on the prompt of a subagent; in the last round it passed grader 01 in 3 of 4 runs without seeing the call, which is the worse case for a point of the factory. The engineer reads the retro to judge the increment, and from increment 08 on its factory points become issues. Doubtful: the retro may change again with the manual scenarios, and a retro records the course of an increment by design, so the kind fits only loosely.
- **Protected paths:** none; the branch changes no test file under `test/`, no CI or factory config and no design note. The notes and tests under `core/skills/*/evals/*/files/` are eval fixtures.
- **Failing checks**, `docs/increments/03-setup-script/design.md:75` and `docs/increments/05-tests-per-increment/design.md:46` on `main`: the required job `checks` fails on 370ec76 in `check-tests` with 14 errors: the missing tests of scenarios 03.1 to 03.8 and 05.1 to 05.5, and the test of scenario 02.7, which the note of increment 05 replaces. So the pull request cannot merge without the bypass. Doubtful as a finding of this increment: `origin/main` (5cb2af7) alone gives the same 14 errors, the branch adds no automated scenario, and the note of increment 05 plans the fix.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 102 checked`
- `factory check-tests --base origin/main`: `check-tests: 14 errors`; the same on `origin/main` alone
- `factory hook stop` on this branch before this result: blocks and names 4bd0a2d900cdab502cc8482b5bc18c66ec39cb44
- `git merge-tree --write-tree HEAD origin/main`: no conflict with the 3 commits of `main` that the branch lacks (#60, #61, #62)
- CI on 370ec76: `test` and `tests-start-red` pass; `checks` fails in `check-tests` with the 14 errors above; `cross-review` waits for a review result for 4bd0a2d

Each `factory` command ran with a binary built from the branch and with one built from `origin/main` (5cb2af7), with the same output.

## Demonstrated scenarios

- **scenario 04.1**: `evidence/04-design-questions/evals.md` and `evidence/04-design-questions/review-runs.md`. The review read the four kept traces of the round after the third review, two with the skill before 86e761d and two with the skill of this commit; the copies of the skill and of the case in that round equal the branch. In each run, the design agent starts one subagent whose prompt names the note and holds neither `--loud` nor Birgit; the new prompts say "Read the whole design note", the old ones "Read only that file" or "Read only that note". Each subagent reads the note once, without offset or limit. The design agent shows the 3 returned questions under `Comprehension questions:`, unanswered, one on the failure signal; with the new skill it rewords them slightly. No run writes a file. This case does not separate the two versions here, as `evals.md` says. Against the skill before this increment, the first table of `evals.md` shows that the old version asked no questions. `explains-gaps-after-answers` ran last in the cross-review round; the paragraph of the skill that it tests has not changed since.
- **scenario 04.2**: `evidence/04-design-questions/evals.md`: `explains-deviations`, old 0.4 against new 0.6, then 1.0 and 1.0 with the skill of this commit; `reviewing-increments` last changed in c8e3b63, before those runs. `clean-change-has-no-deviations` ran only before aca95e2 and c8e3b63, old 0.4 against new 0.6, and its graders on the review result fail in the sandbox of the runner. Not run again.
- **scenario 04.3**: `evidence/04-design-questions/evals.md`: `sketches-changed-flow`, old 0.0 against new 0.8, then 1.0 and 1.0 with the skill of this commit; no sketch for a detail in grader 05 of `explains-deviations`. Not run again.
