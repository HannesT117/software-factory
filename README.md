# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. It is a Go CLI, `factory`. Projects install it with `factory install`. The first project is `health`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it.
- **Increment:** one unit of work. It runs in six steps: design note, tests from a separate session, implementation, fresh-context review, demonstration, approval. The engineer approves the note and the tests with commits signed by Touch ID. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 02 of 07. Branch `claude/factory-v0`.

| Increment | Content | State |
|---|---|---|
| 01 Basis | this page, principles, research notes, ADR 0001 | done |
| 02 Trust core | Go CLI with checks for skills, design notes and tests; signed approvals and `factory verify`; the skills `designing-increments` and `writing-acceptance-tests`; guard hook; CI of the factory | [design note](docs/increments/02-trust-core/design.md) waiting for approval |
| 03 Install and distribution | `factory install`, releases with checksums, install in cloud and CI, ruleset | [design note](docs/increments/03-install-and-distribution/design.md) draft |
| 04 Pipeline | receipts, provenance log, workflow gate, holdout, private paths, increment status | open |
| 05 Chaos | sabotage tests, `running-chaos-exercises` skill, one free exercise | open |
| 06 Skills | `reviewing-architecture`, `writing-docs`, PR record | open |
| 07 Handover | end-to-end test, `adopt.md`, note for `health` | open |

## Decisions

| ADR | Decision | Status |
|---|---|---|
| [0001](docs/decisions/0001-go-cli-and-install-by-copy.md) | Go CLI with thin adapters per agent tool. Projects get a copy of the skills. | proposed (reopened 2026-10-26) |

ADRs are only for architecture decisions and use the [MADR](https://adr.github.io/madr/) format. Only the engineer sets an ADR to `accepted`.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate: rule 16
- no private data in public repos: rule 26
- the order of work, who writes the tests, and signed approvals: rules 27–32

## Open

- Approval of ADR 0001 (Go CLI) and of the increment 02 design note with its open questions
- Before the first approval: a Touch ID key in Secretive, and its public key in `.factory/allowed_signers` on `main`
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
| `factory install --watch`: copies factory changes into the project at once | the factory changes often during `health` work |
| Architecture model as code, checked against the import graph | `health` has at least two packages |

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note and evidence.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
