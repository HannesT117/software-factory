# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. It is a Go CLI, `factory`. Projects install it with `factory install`. The first project is `health`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it, with its own restricted token.
- **Increment:** one unit of work in three pull requests: design note, tests from a separate session, implementation with a review from a separate session, demonstration and retro. The engineer approves each pull request on GitHub, and agents cannot. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 02 of 04. `health` starts after increment 04.

| Increment | Content | State |
|---|---|---|
| 01 Basis | this page, principles, research notes, ADRs 0001 and 0002 | done |
| 02 Trust core | Go CLI with checks for skills, design notes, ADRs, tests and pull request paths; approval gate on GitHub; skills for the workflow, pull requests and commit messages; guard hook | [design note](docs/increments/02-trust-core/design.md) waiting for approval |
| 03 Install and distribution | `factory install`, releases with attestations, install in cloud and CI | [draft](docs/increments/03-install-and-distribution/draft.md) |
| 04 Handover and feedback loop | retro per increment, `factory check-retro`, `factory feedback` (prefilled issue links for the factory repo), version updates, end-to-end test, note for `health` | [draft](docs/increments/04-handover-and-feedback/draft.md) |

## Decisions

| ADR | Decision |
|---|---|
| [0001](docs/decisions/0001-go-cli-with-thin-adapters.md) | The factory tooling is a Go CLI with thin adapters per agent tool. |
| [0002](docs/decisions/0002-delivery-as-pinned-copies.md) | Projects get committed copies of the factory's parts and pin the version. |

ADRs are only for architecture decisions and use the [MADR](https://adr.github.io/madr/) format. An ADR on `main` is in force: the engineer accepts it by merging its pull request. The table lists only ADRs in force. Superseded and deprecated ADRs stay in the folder.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate and approvals: rules 16 and 31
- no private data in public repos: rule 26
- order of work, test writer, protected paths, agent token, checks from `main`, feedback to the factory: rules 27–35
- docs show the current state and follow their format: rule 36
- pull request descriptions: rule 37

## Open

- Review of ADRs 0001 and 0002, and of the increment 02 design note with its open questions
- Setup before the first pull request: agent token, environment `approval`, approval workflow, ruleset for `main`. Steps in [docs/setup/](docs/setup/README.md).

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note, review, evidence and retro.
- [docs/backlog.md](docs/backlog.md): parts that are built when `health` needs them.
- [docs/setup/](docs/setup/README.md): the engineer's one-time setup per repo.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
