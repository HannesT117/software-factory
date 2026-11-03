# Evals of increment 04

Demonstrated scenarios 04.1 to 04.3, run with `claude plugin eval` on the skill folders, the run model `claude-sonnet-5-5` and the judge model `haiku`. Each score is the share of a case's graders that passed.

## Old against new version

The old version is the skill on `main` before this increment, with the new cases. One run per case and version.

| Skill | Case | Scenario | Old | New |
|---|---|---|---|---|
| reviewing-increments | clean-change-has-no-deviations | 04.2 | 0.4 | 0.6 |
| reviewing-increments | explains-deviations | 04.2 | 0.4 | 0.6 |
| reviewing-increments | sketches-changed-flow | 04.3 | 0.0 | 0.8 |
| designing-increments | asks-comprehension-questions | 04.1 | 0.0 | 0.8 |
| designing-increments | explains-gaps-after-answers | 04.1 | 0.67 | 0.67 |

The old reviewer asked comprehension questions and listed no deviations; the old design skill asked no questions.

## After the fixes

The skills changed after the first runs: the question of a sketch goes into its Mermaid title; a deviation that is also a finding gets an entry; a detail inside a step gets no sketch; the explanation of the gaps is prose, one paragraph per question.

| Case | Runs | Scores |
|---|---|---|
| explains-deviations | 2 | 1.0, 1.0 |
| sketches-changed-flow | 2 | 1.0, 1.0 |
| asks-comprehension-questions | 1 | 1.0 |
| explains-gaps-after-answers | 2 | 1.0, 1.0 |

explains-gaps-after-answers loaded the skill only after its prompt named it; before, the runs scored 0.8 with the explanation as lists. In the real flow, the design session has the skill loaded already.

## Limits

- In the runner's sandbox, shell commands could not write, so no run built or tested the fixture code; the graders on the review result in clean-change-has-no-deviations fail for both versions for this reason.
- Not run again: the unchanged cases finds-hard-coded-test-case and test-misses-its-scenario of reviewing-increments, and asks-before-deciding-terms, replaces-earlier-scenario and writes-note-from-answers of designing-increments.
