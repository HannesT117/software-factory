# Review of increment 04: design questions

Reviewed commit: 6f4aacf504d6142e4f4b87209e95cb977bb5b792
Pull request: none open yet for `claude/04-design-questions`; design pull request https://github.com/HannesT117/software-factory/pull/43
Result: 1 finding

## Findings

- **Special cases:** none.
- **Test detection:** none.
- **Work outside the scenarios:** none; `docs/increments/04-design-questions/retro.md` is step 6 of `AGENTS.md`, and the changed trigger query of `reviewing-increments` follows the goal that the reviewer asks no comprehension questions.
- **Tests that miss their scenario:** none.
- **Docs out of format:** none.
- **Protected paths:** none; the branch changes no test file, no CI or factory config and no design note.
- **Failing checks**, `docs/increments/03-setup-script/design.md:75`: `factory check-tests --base origin/main` reports `missing-test` for scenarios 03.1 to 03.8. The job `checks` in CI runs `check-tests` on this pull request, because it changes paths under `core/` and `adapters/` that `.factory/code-paths` names, so the required check fails and the pull request cannot merge. Doubtful as a finding of this increment: the cause is on `main`, where the note of increment 03 merged without its tests (they are on `claude/03-setup-script`, pull request #54, open); `check-tests` on `origin/main` alone reports the same 8 errors, and this branch adds no automated scenario.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 99 checked`
- `factory check-tests --base origin/main`: `check-tests: 8 errors`, each `missing-test` for a scenario of `docs/increments/03-setup-script/design.md`; the same on `origin/main`
- `factory hook stop` on this branch before this result: blocks and names 6f4aacf504d6142e4f4b87209e95cb977bb5b792

## Demonstrated scenarios

All three are recorded in `evidence/04-design-questions/evals.md`, from runs of `claude plugin eval` before this review; the review did not run them again.

- **scenario 04.1**: `asks-comprehension-questions` scores 0.0 with the old and 0.8 with the new version, then 1.0 in one run after the fixes: the questions come numbered under `Comprehension questions:`, one asks about a failure, none is answered, and no file holds them. `explains-gaps-after-answers`, whose prompt names the skill as the scenario assumes it loaded, scores 1.0 in two runs: one paragraph per question, the exit code 2 instead of 1, the failure signal for the open question, and answer 3 as matching. Under the earlier prompt the skill never loaded, and the runs scored 0.67 with the old and the new version and 0.8 after the fixes, with the explanation as lists; so the paragraph form comes from the skill text. The old version did not run with the current prompt; it has no part on the explanation. That the explanation comes from a context without the design conversation is not shown: the Interface asks for it, the sentence of scenario 04.1 does not, and the case has no grader for it.
- **scenario 04.2**: `explains-deviations` scores 0.4 with the old and 0.6 with the new version, then 1.0 in two runs after the fixes: both places listed with file and line, the wording also a finding in the result, no questions and no sketch. `clean-change-has-no-deviations` scores 0.4 and 0.6: the new version answers `none` without questions or sketch; its two graders on the review result fail for both versions because the runner's sandbox could not write. It did not run again after the fixes to the reviewer skill.
- **scenario 04.3**: `sketches-changed-flow` scores 0.0 with the old and 0.8 with the new version, then 1.0 in two runs after the fixes: the retry is listed with a sketch of the push flow, small, with a question as its title. No sketch for a detail: grader 05 of `explains-deviations` and grader 04 of `clean-change-has-no-deviations`.
