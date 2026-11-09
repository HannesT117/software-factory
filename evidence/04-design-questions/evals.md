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

- The judge of the runner sees a shortened trace without the subagent call, so it failed grader 01 in all runs. For the new version, a script checked the full traces: both runs start a subagent whose prompt holds the note path, the three questions and the three answers, and neither detail of the design conversation; the scores above count that result. In run 1 the design agent shows the explanation word for word; in run 0 it changes "the person" to "you", with the same content. The traces of the old version were not kept; its final answers say that the agent explained the gaps itself.
- Grader 02 of explains-gaps-after-answers failed once in each version: after the three paragraphs, the agent adds a line that nothing blocks the approval.
- In asks-comprehension-questions, grader 01 failed once with the unchanged skill.

## After the third review

The skill now says that the fresh session reads the whole note. The prompt of asks-comprehension-questions carries two details of a design conversation (the former flag name `--loud` and a colleague named Birgit), and grader 01 checks that the session's prompt holds neither, does not limit the session to parts of the note, and that the session reads the whole note. Two runs per version; the old version is the skill before this change.

| Case | Version | Judge scores | Scores after the script checks |
|---|---|---|---|
| asks-comprehension-questions | old | 0.8, 0.8 | 1.0, 1.0 |
| asks-comprehension-questions | new | 0.6, 1.0 | 1.0, 1.0 |

- [check-traces.py](check-traces.py) checks graders 01 and 02 on the full traces of a run with `--keep-temp`; grader 05 was checked by reading the tool calls. In all four runs, the subagent's prompt names the note, holds neither detail of the conversation, and the subagent reads the note without offset or limit; the shown questions have the content of the subagent's questions; no run writes a file. The judge failed some of these graders; where its verdict differs from the trace, the score counts the trace.
- The case does not separate the versions here: the old skill also had the whole note read in both runs. The difference is in the subagent's prompt: the new version says "Read the whole design note", the old one "Read only that note". In the run of the third review, the old skill limited the session to sections.
- With the new skill, the design agent reworded the questions slightly; the skill asks it to show all of them, not to show them unchanged.
- A first try of this round was dropped: its prompt said that a greeting in German may come later, and the design agent then tried to change the note, which the read-only sandbox refused, and asked no questions in two of four runs.

## Graders on the trace

After cross-review, the graders that judge a subagent's prompt are `tool_used` graders: they match the input of the Agent call (the note and the request to read it whole, the answers, no detail of the design conversation) and of each Read of the note (whole, never with offset or limit). [check-traces.py](check-traces.py) adds what these graders cannot tell: that the subagent itself reads the note, and what the subagent returned next to what the design agent showed. Two runs per case with the new skill:

| Case | Scores |
|---|---|
| asks-comprehension-questions | 0.88, 1.0 |
| explains-gaps-after-answers | 0.88, 1.0 |

- Every `tool_used` grader passed in all four runs; the runner counted the Agent call and the subagent's Read. On the kept traces of the round before, grader 01 of asks-comprehension-questions fails for the old skill in both runs, because its prompt to the subagent did not ask for the whole note, and passes for the new one.
- The script passed the subagent's reads and the shown questions in all runs. In one run of explains-gaps-after-answers the design agent put a bold label before each paragraph of the explanation, with the same content.
- The two failed graders are judged by the model: grader 08 of asks-comprehension-questions in a run that wrote no file, and grader 04 of explains-gaps-after-answers, where the agent adds a line on the approval after the three paragraphs.

## Every grader on the middle of a trace on tool calls

After the fifth review, each grader of the five cases that judged something between the skill call and the final answer, which the judge does not see, matches tool calls instead: the Write of the review result and its content, no Write or Edit of code and tests, no Write, Edit or commit where a run must change no file. The scaffold of clean-change-has-no-deviations fixes the commit dates and ignores the host's git config, so that the latest commit has the same hash in every run and a grader names it. Two runs per case with the skills of this branch:

| Case | Scores |
|---|---|
| asks-comprehension-questions | 1.0, 1.0 |
| explains-gaps-after-answers | 1.0, 1.0 |
| clean-change-has-no-deviations | 1.0, 0.71 (after the fix of the hash) |
| explains-deviations | 0.8, 0.8 |
| sketches-changed-flow | 1.0, 1.0 |

- Shell commands cannot write in the runner's sandbox, but the Write tool can. In the failed runs of the graders on the review result, the reviewer concluded from a failed shell write that it could not write at all and did not try the Write tool; in explains-deviations both runs gave the result in the final answer instead, with the finding on the wording "is known". The other runs wrote the result with the Write tool.
- Before the fix of the hash, clean-change-has-no-deviations scored 0.57 and 0.71 with the right hash in both written results: grader 05 failed because the reviewer counted its unrun checks as a finding, and grader 02 failed once because the answer listed a deviation.
- In run 0 of asks-comprehension-questions the final answer goes on after the questions with a remark that answers question 3; the judge passed grader 07. Its text now names text before or after the questions, and check-traces.py checks the subagent's report.
- check-traces.py, in its version after the sixth review, passes every check of the designing cases but one: in run 0 of explains-gaps-after-answers the design agent changed "the person" to "you" in two paragraphs of the explanation, which the skill asks it to show unchanged. The content is the same. Not fixed: the fix would be in the skill, with another round of evals and review, for a change of address only.

## No remark after the questions

The seventh review found that in run 0 of asks-comprehension-questions above, the design agent's remark after the questions answers question 3 and part of question 2; judged again with the new text of grader 07, that run fails it, so the case scored 0.9 and 1.0 there. The skill now says that the agent adds no remark on what the questions ask. Three runs per version:

| Case | Version | Scores |
|---|---|---|
| asks-comprehension-questions | old | 1.0, 1.0, 0.9 |
| asks-comprehension-questions | new | 1.0, 1.0, 1.0 |

- The old version failed grader 07 once, but not through a remark after the questions: in that run the fresh session's own question 2 names one of the two signals it asks for, and the design agent shows it word for word. With the new version, the design agent's prompt to the fresh session asks for questions without hints in all three runs, which no old prompt does; two runs show nothing after the questions and one only the next step.
- check-traces.py: in one run of each version the design agent reworded one question: "the two signals" became "the signals", and "gets no name in the output" became "gives no name". The skill asks for all questions, not for their exact words, so these are no failures of the skill.

## Limits

- In the runner's sandbox, shell commands cannot write, so no run builds or tests the fixture code, and a reviewer that tries a shell write first may not try the Write tool.
- Not run again: the unchanged cases finds-hard-coded-test-case and test-misses-its-scenario of reviewing-increments, and asks-before-deciding-terms, replaces-earlier-scenario and writes-note-from-answers of designing-increments.
