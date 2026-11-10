# Evals of reviewing-increments

The skill `reviewing-increments` changed in increment 06: each finding is blocking or a note, the Result line counts both, and the reviewer reports nothing about the wording of retros or the README. The graders of three cases now read the new Result line.

Run with `claude plugin eval` on the skill folder, `--runs 1 --ablation none --allow-tools Write Edit Bash Read Glob Grep`, the run model `claude-opus-5-5`, which the sessions use, and the judge model `haiku`. The old version is the skill on `main` before this increment, with the new graders. Each score is the share of a case's graders that passed. Cost: 1.45 USD for the old version, 1.46 USD for the new one.

| Case | Old | New | Failed graders |
|---|---|---|---|
| clean-change-has-no-deviations | 0.71 | 0.86 | both: 02, the final answer says there are no deviations (judge); old: 05, the Result line in the new form |
| explains-deviations | 0.80 | 1.00 | old: 05, no sketch (judge) |
| finds-hard-coded-test-case | 0.67 | 1.00 | old: 01, the result names the special case (judge) |
| sketches-changed-flow | 1.00 | 1.00 | none |
| test-misses-its-scenario | 0.50 | 1.00 | old: 02, the Result line counts a blocking finding (judge, 1 of 3 votes) |
| Overall | 0.74 | 0.97 | |

- The old version fails grader 05 of clean-change-has-no-deviations by design: it writes `Result: no findings`.
- The other failures are verdicts of the judge on the trace, of which it sees only the start and the end (retro of increment 04). One run per case and version does not separate them from noise; they do not show a loss of the new version.
