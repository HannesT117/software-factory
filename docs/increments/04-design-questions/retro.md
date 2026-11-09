# Retro of increment 04

## Metrics

- Pull requests: the design note (#43) and the skills with their evals.
- Engineer decisions in the conversation: four, on who asks, the debugging question, the reviewer's part after the code, and the sketches.
- Eval runs: 10 runs old against new, then 9 runs after 5 fixes to the skills or the cases; about 2.50 USD.
- Review: 2 findings. One was an eval case that never loaded the skill, fixed; the other is the missing tests of increment 03 on `main`, which increment 05 addresses.
- cross-review, round 1: 1 blocking finding and 1 note on the evals of `designing-increments`, both fixed; the note showed that the design agent explained the gaps itself, so the skill now has it start a fresh session. 6 eval runs, and 2 more in the third review.
- Third review: 5 doubtful findings. Fixed at once: the status in `README.md` and this retro.
- cross-review, round 2: 1 blocking finding, the same as one of the third review: the fresh session must read the whole note. Fixed in the skill, and the prompt of asks-comprehension-questions got a detail of a design conversation for grader 01. 8 eval runs: 4 recorded, 4 dropped because their prompt led the design agent to change the note.
- Fourth review: 3 doubtful findings. Fixed: this retro, and the trace check is now a script in the evidence. Not fixed: the failing `checks`, the same as on `main`.
- cross-review, rounds 3 and 4: the subagent graders relied on a script outside the cases, and the script counted the design agent's own Read. Fixed with `tool_used` graders on the Agent and Read calls, and the script takes only the subagent's Reads. 4 eval runs.
- Fifth review: 3 findings. The engineer decided to fix the graders in this pull request: every grader of the five cases on something in the middle of a trace now matches tool calls. Also fixed: this retro and a merge conflict with `main` in `README.md`.

## Agent version

- Claude Code 2.1.317, from the cross-review findings on 2.1.319; orchestrating model `claude-opus-5-5`; eval runs `claude-sonnet-5-5`, judged by `haiku`.

## Corrections

- `claude plugin eval` refused to run until the cloud container had `bubblewrap` and `socat`; in its sandbox, shell commands could not write.
- The first runs showed three unclear places in `reviewing-increments`: where the question of a Mermaid sketch goes, that a finding also gets a deviation entry, and that a detail inside a step is no new step.

## Points for the factory

- An eval case can pass or fail without the skill ever loading. A grader that checks the skill call, or a look at the trace, shows it.
- `factory check-tests` rejects the code of every increment while one design note on `main` waits for its tests: increment 05.
- The cloud proxy adds a session link to each new pull request description; the agent removes it at once.
- The judge of `claude plugin eval` sees only the first and last 12 lines of a trace, and every streamed chunk is a line, so it sees the skill call and the final answer, never the subagent calls, Writes or commits in between; its verdict on them is a guess, pass or fail. `tool_used` graders on the calls' input check them; a script on the full trace, kept with `--keep-temp`, checks what a regex cannot, such as whose Read a call is. `authoring-skills` says the opposite and needs a fix; the graders of the other skills need the same change.
- In the cloud container the runner runs as root, and its sandbox makes every path read-only for shell commands; only the Write tool can write. A run that tries a shell write first may conclude that it cannot write at all.
- A detail in an eval prompt can change what the agent does: a remark on a later feature made the design agent change the note instead of asking questions.
- cross-review measures a pull request outside any increment against the latest note: increment 05, scenario 05.6.

## Summary

The engineer now answers comprehension questions on the design pull request, from a fresh session that knows only the note, with at least one question on debugging. After the code, the reviewer lists where the code deviates from the note and draws a small sketch only for a changed flow. Evals show the new behaviour for both skills; the old versions did not ask or list. Five fixes after the first eval runs sharpened the skills and one case. Open: the manual scenarios 04.4 and 04.5.
