# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. Projects install it with `install.mjs`. The first project is `health`.

Roles: the **engineer** is the human who owns a project. An **agent** is a coding agent that works in it.

## Status

v0, phase 2 of 6 (slice skills, then skeleton). Branch `claude/factory-v0`. Every phase is a slice: design note, approved scenarios, reviewed acceptance tests, then code (see [AGENTS.md](AGENTS.md)).

| Phase | Content | State |
|---|---|---|
| 1 Basis | this page, principles, research notes, ADR 0001 | done |
| 2a Slice skills | `plan-slice`, `write-acceptance`, skill check, design and test checks | [design note](docs/slices/02a-slice-skills/design.md) waiting for approval |
| 2 Skeleton | `install.mjs`, CI, ruleset | [design note](docs/slices/02-skeleton/design.md) draft, rewritten with `plan-slice` after 2a |
| 3 Pipeline | checks with receipts, lock, private paths, slice status, hooks, adapter | open |
| 4 Chaos | sabotage tests, `run-chaos` skill, one free exercise | open |
| 5 Skills | `write-acceptance`, `plan-slice`, `architecture-review`, `writing` | open |
| 6 Handover | end-to-end test, `adopt.md`, note for `health` | open |

## Decisions

| ADR | Decision | Status |
|---|---|---|
| [0001](docs/decisions/0001-layers-and-install-by-copy.md) | Vendor-neutral core plus thin adapters. Projects get a copy, not a link. | accepted |

ADRs are only for architecture decisions and use the [MADR](https://adr.github.io/madr/) format. Only the engineer sets an ADR to `accepted`.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate: rule 16
- no private data in public repos: rule 26

## Open

- Approval of the slice 02a design note, its scenarios and its open questions
- Phase 2: ruleset import and public repos
- Phase 5: rules for code comments and explanations, choice of pstack parts

## Later: built when `health` needs it

v0 leaves these parts out on purpose. Each part is built at the point where `health` needs it. A slice is one step of the `health` plan.

| Part | Built when |
|---|---|
| `coding-task` skill: a task the engineer implements by hand to learn a pattern | `health` slice 1a, first new pattern |
| `prove-with-media` skill: video and screenshots as evidence for UI changes | `health` slice 1b, first UI |
| 2 am support: each slice names the log event that shows it failing; timed debugging drills; an `investigate-incident` skill | `health` slice 2, first cron and queue jobs |
| Product chaos, for example flooding a test database, with a check that the target is not production | `health` has a preview database |
| Decision-model check routing through Opper | `health` has enough changes to measure misses |
| `install.mjs --watch`: copies factory changes into the project at once | the factory changes often during `health` work |
| Architecture model as code, checked against the import graph | `health` has at least two packages |

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/slices/](docs/slices/): one design note per slice.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
