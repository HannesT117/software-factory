# Review of increment 04: design questions

Reviewed commit: 4054041fe6a48629feb3905225adcce4627f08f7
Pull request: none open yet for `claude/04-design-questions`; design pull request https://github.com/HannesT117/software-factory/pull/43
Result: 2 findings

## Findings

- **Special cases:** none.
- **Test detection:** none.
- **Work outside the scenarios:** none.
- **Tests that miss their scenario:** none.
- **Docs out of format:** none.
- **Protected paths:** none; the branch changes no test file, no CI or factory config and no design note.
- **Failing checks**, `evidence/04-design-questions/evals.md:30`: the eval `explains-gaps-after-answers` of scenario 04.1 fails its grader `02-one-paragraph-per-question` in both runs after the fixes. The explanation of the gaps comes as lists, while the Interface of the note fixes one paragraph per question, and `core/skills/designing-increments/SKILL.md:40` asks for it. The engineer reads this explanation on every design pull request, so the form that the note chose does not reach the engineer. Doubtful because the eval is a model run judged by `haiku`, not part of `go test`, the sentence of scenario 04.1 itself names no form, and the content of the explanation is right in both runs.
- **Failing checks**, `docs/increments/03-setup-script/design.md`: `factory check-tests --base origin/main` reports `missing-test` for scenarios 03.1 to 03.8. The job `checks` in CI runs `check-tests` on this pull request, because it changes paths under `core/` and `adapters/`, so the required check fails and the pull request cannot merge. Doubtful as a finding of this increment: the cause is on `main`, where the note of increment 03 merged without its tests (they are on `claude/03-setup-script`, pull request #54); `check-tests` on `origin/main` alone reports the same 8 errors, and this branch adds no automated scenario.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/04-design-questions`: `check-design: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 96 checked`
- `factory check-tests --base origin/main`: `check-tests: 8 errors`, each `missing-test` for a scenario of `docs/increments/03-setup-script/design.md`; the same on `origin/main`
- `factory hook stop` on this branch before this result: blocks and names 4054041fe6a48629feb3905225adcce4627f08f7

## Demonstrated scenarios

All three are recorded in `evidence/04-design-questions/evals.md`, from runs of `claude plugin eval` before this review; the review did not run them again.

- **scenario 04.1**: `asks-comprehension-questions` scores 0.0 with the old and 0.8 with the new version, then 1.0 in one run after the fixes: the questions come numbered under `Comprehension questions:`, one asks about a failure, none is answered, and no file holds them. `explains-gaps-after-answers` scores 0.67 with both versions, then 0.8 in two runs: the content of the explanation is right, its form is not (finding above). The old version scores the same, so the case shows the explanation but not that the skill brings it about. That the explanation comes from a context without the design conversation is not shown: its grader was dropped, because the eval run has no design conversation; scenario 04.4 covers it in use.
- **scenario 04.2**: `explains-deviations` scores 0.4 with the old and 0.6 with the new version, then 1.0 in two runs after the fixes: both places listed with file and line, the wording also a finding in the result, no questions and no sketch. `clean-change-has-no-deviations` scores 0.4 and 0.6: the new version answers `none` without questions or sketch; its two graders on the review result fail for both versions because the runner's sandbox could not write. It did not run again after the fixes to the reviewer skill.
- **scenario 04.3**: `sketches-changed-flow` scores 0.0 with the old and 0.8 with the new version, then 1.0 in two runs after the fixes: the retry is listed with a sketch of the push flow, small, with a question as its title. No sketch for a detail: grader 05 of `explains-deviations` and grader 04 of `clean-change-has-no-deviations`.
