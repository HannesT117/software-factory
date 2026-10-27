# Backlog

v0 leaves these parts out on purpose. Each part is built when a project that uses the factory needs it.

| Part | Built when |
|---|---|
| A skill for coding tasks that the engineer implements by hand to learn a pattern | a project needs its first new pattern |
| `proving-with-media` skill: video and screenshots as evidence for UI changes | the first project with a UI |
| 2 am support: each increment names the log event that shows it failing; a fix line for each error code; timed debugging drills; an `investigating-incidents` skill | the first project with scheduled or queued jobs |
| Receipts (rule 14), provenance log, holdout scenarios (rule 7), private-paths check (rule 26) | an observed failure shows the need |
| `run-checks` with computational checks before inferential ones (rule 5) | the first project has more than one kind of check |
| Property statements by the engineer, generators by the agent (rule 8) | the first logic with clear invariants |
| Chaos: sabotage tests and a `running-chaos-exercises` skill; product chaos with a check that the target is not production (rule 23) | the harness or a project has something worth breaking |
| `reviewing-architecture` and `writing-docs` skills | the second increment of a project |
| `factory status` (what waits for the engineer) and `factory next` (starts the session that fits the current step) | the engineer's steps per increment cost more time than the work |
| Decision-model check routing through Opper | a project has enough changes to measure misses |
| `factory install --watch`: copies factory changes into the project at once | the factory changes often while a project uses it |
| Architecture model as code, checked against the import graph | a project has at least two packages |
| A check for pull request titles: length, capital letter, final period, internal numbers | the same title correction is needed a second time (rule 4) |
