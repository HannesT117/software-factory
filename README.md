# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. It is a Go CLI, `factory`. Projects install it with `factory install`. The first project is `health`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it, with its own restricted token.
- **Increment:** one unit of work: design note, tests from a separate session, implementation, review from a separate session, demonstration, squash merge. The engineer approves the note and the tests with commits signed by Touch ID. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 02 of 04. Branch `claude/factory-v0`. `health` starts after increment 03.

| Increment | Content | State |
|---|---|---|
| 01 Basis | this page, principles, research notes, ADR 0001 | done |
| 02 Trust core | Go CLI with checks for skills, design notes and tests; signed approvals with `factory approve` and `factory verify`; skills `designing-increments`, `writing-acceptance-tests`, `reviewing-increments`; guard hook | [design note](docs/increments/02-trust-core/design.md) waiting for approval |
| 03 Install and distribution | `factory install`, releases with attestations, install in cloud and CI, ruleset | [design note](docs/increments/03-install-and-distribution/design.md) draft |
| 04 Handover | end-to-end test, `adopt.md`, note for `health` | open |

## Decisions

| ADR | Decision | Status |
|---|---|---|
| [0001](docs/decisions/0001-go-cli-and-install-by-copy.md) | Go CLI with thin adapters per agent tool. Projects get a copy of the skills. | proposed (reopened 2026-10-26) |

ADRs are only for architecture decisions and use the [MADR](https://adr.github.io/madr/) format. Only the engineer sets an ADR to `accepted`.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate: rule 16
- no private data in public repos: rule 26
- order of work, test writer, approvals, agent token, verifier source: rules 27–34

## Open

- Approval of ADR 0001 and of the increment 02 design note with its open questions
- Setup before the first approval: an agent token, a Touch ID key in Secretive registered on GitHub as a signing key, and the engineer's login in `.factory/config.json`
- Increment 03: ruleset import and public repos

## Later: built when `health` needs it

v0 leaves these parts out on purpose. Each part is built at the point where `health` needs it. Numbers such as 1a refer to the steps of the `health` plan.

| Part | Built when |
|---|---|
| A skill for coding tasks that the engineer implements by hand to learn a pattern | `health` 1a, first new pattern |
| `proving-with-media` skill: video and screenshots as evidence for UI changes | `health` 1b, first UI |
| 2 am support: each increment names the log event that shows it failing; timed debugging drills; an `investigating-incidents` skill | `health` 2, first cron and queue jobs |
| Receipts, provenance log, workflow gate, holdout scenarios, private-paths check | an observed failure shows the need |
| Chaos: sabotage tests and a `running-chaos-exercises` skill; product chaos with a check that the target is not production | the harness or `health` has something worth breaking |
| `reviewing-architecture` and `writing-docs` skills, PR record template | the second `health` increment |
| Decision-model check routing through Opper | `health` has enough changes to measure misses |
| `factory install --watch`: copies factory changes into the project at once | the factory changes often during `health` work |
| Architecture model as code, checked against the import graph | `health` has at least two packages |

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note, approvals, review and evidence.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
