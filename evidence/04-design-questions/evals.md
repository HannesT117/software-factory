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

## After the findings of cross-review

Grader 02 of asks-comprehension-questions now checks that the shown questions are all those of the fresh session. explains-gaps-after-answers got two details of a design conversation in its prompt and grader 01, which checks that the explanation comes from a session that holds neither. Two runs per version; the old version is the skill before the change of the second paragraph of Comprehension questions.

| Case | Version | Scores |
|---|---|---|
| asks-comprehension-questions | unchanged skill | 1.0, 0.8 |
| explains-gaps-after-answers | old | 0.5, 0.67 |
| explains-gaps-after-answers | new | 0.83, 1.0 |

- The judge of the runner sees a shortened trace without the subagent call, so it failed grader 01 in all runs. For the new version, a script checked the full traces: both runs start a subagent whose prompt holds the note path, the three questions and the three answers, and neither detail of the design conversation; the scores above count that result. The traces of the old version were not kept; its final answers say that the agent explained the gaps itself.
- Grader 02 of explains-gaps-after-answers failed once in each version: after the three paragraphs, the agent adds a line that nothing blocks the approval.
- In asks-comprehension-questions, grader 01 failed once with the unchanged skill.

## Limits

- In the runner's sandbox, shell commands could not write, so no run built or tested the fixture code; the graders on the review result in clean-change-has-no-deviations fail for both versions for this reason.
- Not run again: the unchanged cases finds-hard-coded-test-case and test-misses-its-scenario of reviewing-increments, and asks-before-deciding-terms, replaces-earlier-scenario and writes-note-from-answers of designing-increments.
