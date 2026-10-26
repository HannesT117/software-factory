# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. Projects install it with `install.mjs`. The first project is `health`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it.
- **Increment:** one unit of work. It runs in five steps: design note, approved scenarios, acceptance tests by a separate agent, implementation, approval. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 02 of 07. Branch `claude/factory-v0`.

| Increment | Content | State |
|---|---|---|
| 01 Basis | this page, principles, research notes, ADR 0001 | done |
| 02 Workflow skills | `designing-increments`, `writing-acceptance-tests`, checks for skills, design notes and tests, test lock, guard hook | [design note](docs/increments/02-workflow-skills/design.md) waiting for approval |
| 03 Skeleton | `install.mjs`, CI, ruleset | [design note](docs/increments/03-skeleton/design.md) draft |
| 04 Pipeline | receipts, full lock with provenance log, workflow gate, private paths, increment status | open |
| 05 Chaos | sabotage tests, `running-chaos-exercises` skill, one free exercise | open |
| 06 Skills | `reviewing-architecture`, `writing-docs`, PR record | open |
| 07 Handover | end-to-end test, `adopt.md`, note for `health` | open |

## Decisions

| ADR | Decision | Status |
|---|---|---|
| [0001](docs/decisions/0001-layers-and-install-by-copy.md) | Vendor-neutral core plus thin adapters. Projects get a copy, not a link. | accepted |

ADRs are only for architecture decisions and use the [MADR](https://adr.github.io/madr/) format. Only the engineer sets an ADR to `accepted`.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate: rule 16
- no private data in public repos: rule 26
- the order of work and who writes the tests: rules 27–30

## Open

- Approval of the increment 02 design note, its scenarios and its open questions
- Increment 03: ruleset import and public repos
- Increment 06: rules for code comments and explanations, choice of pstack parts

## Later: built when `health` needs it

v0 leaves these parts out on purpose. Each part is built at the point where `health` needs it. Numbers such as 1a refer to the steps of the `health` plan.

| Part | Built when |
|---|---|
| `coding-task` skill: a task the engineer implements by hand to learn a pattern | `health` 1a, first new pattern |
| `prove-with-media` skill: video and screenshots as evidence for UI changes | `health` 1b, first UI |
| 2 am support: each increment names the log event that shows it failing; timed debugging drills; an `investigate-incident` skill | `health` 2, first cron and queue jobs |
| Product chaos, for example flooding a test database, with a check that the target is not production | `health` has a preview database |
| Decision-model check routing through Opper | `health` has enough changes to measure misses |
| `install.mjs --watch`: copies factory changes into the project at once | the factory changes often during `health` work |
| Architecture model as code, checked against the import graph | `health` has at least two packages |

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note and evidence.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
