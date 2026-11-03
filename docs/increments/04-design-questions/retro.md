# Retro of increment 04

## Metrics

- Pull requests: the design note (#43) and the skills with their evals.
- Engineer decisions in the conversation: four, on who asks, the debugging question, the reviewer's part after the code, and the sketches.
- Eval runs: 10 runs old against new, then 9 runs after 5 fixes to the skills or the cases; about 2.50 USD.
- Review: 2 findings. One was an eval case that never loaded the skill, fixed; the other is the missing tests of increment 03 on `main`, which increment 05 addresses.

## Agent version

- Claude Code 2.1.317; orchestrating model `claude-opus-5-5`; eval runs `claude-sonnet-5-5`, judged by `haiku`.

## Corrections

- `claude plugin eval` refused to run until the cloud container had `bubblewrap` and `socat`; in its sandbox, shell commands could not write.
- The first runs showed three unclear places in `reviewing-increments`: where the question of a Mermaid sketch goes, that a finding also gets a deviation entry, and that a detail inside a step is no new step.

## Points for the factory

- An eval case can pass or fail without the skill ever loading. A grader that checks the skill call, or a look at the trace, shows it.
- `factory check-tests` rejects the code of every increment while one design note on `main` waits for its tests: increment 05.
- The cloud proxy adds a session link to each new pull request description; the agent removes it at once.

## Summary

The engineer now answers comprehension questions on the design pull request, from a fresh session that knows only the note, with at least one question on debugging. After the code, the reviewer lists where the code deviates from the note and draws a small sketch only for a changed flow. Evals show the new behaviour for both skills; the old versions did not ask or list. Five fixes after the first eval runs sharpened the skills and one case. Open: the manual scenarios 04.4 and 04.5.
