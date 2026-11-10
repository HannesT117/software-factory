# Increment 06: gate and status

Task of the engineer, 2026-11-10: work on the points of `docs/research/2026-11-factory-review.md` (uploaded review), check its best-practice claims against the original pages, ignore existing rules where needed as long as design stays separate from implementation, two pull requests, keep going alone, then a summary with visuals (format of that summary to discuss afterwards, maybe as a factory skill).

## State

- #68 design (note 06 + research note), waits for the engineer; comprehension questions shown in the final report.
- #69 draft, tests and code of 05 and 06 on `claude/06-tests-and-code`; workflow change as patch in its first comment (token lacks scope `workflow`); local `simbase` = main + #67 + #68 for checks.
- #67 note 05 changed: tests-start-red reads unchanged behaviour from the base branch.
- Tests by session_01PcDxYHzSki2JvX18xq1jLD from the unapproved notes, on the engineer's words.
- Eval reviewing-increments old 0.74, new 0.97, session model.
- Summary page: scratchpad summary/review-umsetzung.html, published as artifact.

## Open

- Engineer: merge #67, #68, then main into #69; apply the workflow patch; merge #69 (break glass); `factory setup`; manual scenarios; skill evidence (about 100 eval runs).
- Agent: reviewer result on #69; after the merge, demonstrate 05.8 on CI; cross-review comments after the patch.

## Checkpoints

- 2026-11-10T20:40Z #68 and #69 open; reviewer running; next: review result, artifact, report.
- 2026-11-10T22:50Z review round 2 on 16e4808: 05 no findings; 06 one blocking (test 06.11 changed after the code) for the engineer, one note fixed in #68 (cross-review comments only for a sent request). #69 at e11a92d. Summary page https://claude.ai/artifact/6jGnmnCv8s6a2zvuxgL7GR. Branch claude/06-code-wip deleted (content in #69). Next: the engineer's answers to the comprehension questions of #68 and the decision on the blocking finding; then the merges and the workflow patch.
