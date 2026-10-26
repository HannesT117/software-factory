# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. Projects install it with `install.mjs`. The first project is `health`.

## Status

v0, phase 1 of 6 (basis documents). Branch `claude/factory-v0`.

| Phase | Content | State |
|---|---|---|
| 1 Basis | this page, principles, research notes, ADRs 0001–0003 | done, waiting for approval |
| 2 Skeleton | README, AGENTS.md, `install.mjs`, skill validation, CI, ruleset | open |
| 3 Pipeline | checks with receipts, lock, private paths, slice status, hooks, adapter | open |
| 4 Chaos | sabotage tests, `run-chaos` skill, one free exercise | open |
| 5 Skills | `write-acceptance`, `plan-slice`, `architecture-review`, `writing` | open |
| 6 Handover | end-to-end test, `adopt.md`, note for `health` | open |

## Decisions

| ADR | Decision | Status |
|---|---|---|
| [0001](decisions/0001-layers-and-install-by-copy.md) | Vendor-neutral core plus thin adapters. Projects get a copy, not a link. | proposed |
| [0002](decisions/0002-evidence-policy.md) | Evidence is locked scenarios, human-stated properties and types. Agent unit tests are not evidence. | proposed |
| [0003](decisions/0003-merge-gate.md) | Public repos. A GitHub ruleset is the merge gate. | proposed |

## Open

- Your approval of phase 1 and of the three ADRs.
- Phase 2: you import the ruleset and make the repos public.
- Phase 5: your rules for code comments and explanations, and your pick of pstack parts.

## Backlog after `health` slice 1

1. 2 am support: "How will we see it fail?" in `plan-slice`, timed drills, `investigate-incident`.
2. Product chaos with a D1 isolation pre-flight.
3. Decision-model check routing through Opper.
4. `install.mjs --watch`.
5. Architecture model as code with a model-vs-code check.
6. `coding-task`, `prove-with-media`.

## Where things are

- [principles.md](principles.md): every rule, its source, what it depends on, and when to re-check it.
- [research/2026-10-factory-basis.md](research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [decisions/](decisions/): the ADRs.
