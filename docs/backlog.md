# Backlog

v0 leaves these parts out on purpose. Each part is built at the point where `health` needs it. Numbers such as 1a refer to the steps of the `health` plan.

| Part | Built when |
|---|---|
| A skill for coding tasks that the engineer implements by hand to learn a pattern | `health` 1a, first new pattern |
| `proving-with-media` skill: video and screenshots as evidence for UI changes | `health` 1b, first UI |
| 2 am support: each increment names the log event that shows it failing; a fix line for each error code; timed debugging drills; an `investigating-incidents` skill | `health` 2, first cron and queue jobs |
| Receipts (rule 14), provenance log, holdout scenarios (rule 7), private-paths check (rule 26) | an observed failure shows the need |
| `run-checks` with computational checks before inferential ones (rule 5) | the first project has more than one kind of check |
| Property statements by the engineer, generators by the agent (rule 8) | the first logic with clear invariants |
| Chaos: sabotage tests and a `running-chaos-exercises` skill; product chaos with a check that the target is not production (rule 23) | the harness or `health` has something worth breaking |
| `reviewing-architecture` and `writing-docs` skills, PR record template | the second `health` increment |
| `factory status` (what waits for the engineer) and `factory next` (starts the session that fits the current step) | the engineer's steps per increment cost more time than the work |
| Decision-model check routing through Opper | `health` has enough changes to measure misses |
| `factory install --watch`: copies factory changes into the project at once | the factory changes often during `health` work |
| Architecture model as code, checked against the import graph | `health` has at least two packages |
