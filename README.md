# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. It is a Go CLI, `factory`. Projects install it with `factory install`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it, with its own restricted token.
- **Increment:** one unit of work in two pull requests: the design note, and the tests from a separate session with their code. The engineer approves the design note with its scenarios on GitHub, and agents cannot. The second pull request merges without the engineer only when its tests came first, the code changed only approval-free paths, all checks pass and a reviewer in a fresh context found nothing. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 02 of 04.

| Increment | Content | State |
|---|---|---|
| 01 Basis | this page, principles, research notes, ADRs 0001 and 0002 | done |
| 02 Trust core | Go CLI with checks for skills, design notes, ADRs, tests and pull request paths; approval gate on GitHub; skills for the workflow, pull requests and commit messages; guard hook and stop hook | [design note](docs/increments/02-trust-core/design.md) approved; tests and code in review; CI pull request next |
| 03 Install and distribution | `factory install`, releases with attestations, install in cloud and CI | [draft](docs/increments/03-install-and-distribution/draft.md) |
| 04 Handover and feedback loop | retro per increment, `factory check-retro`, `factory feedback` (prefilled issue links for the factory repo), version updates, end-to-end test | [draft](docs/increments/04-handover-and-feedback/draft.md) |

## Decisions

| Record | Decision |
|---|---|
| [architecture 0001](docs/decisions/architecture/0001-go-cli-with-thin-adapters.md) | The factory tooling is a Go CLI with thin adapters per agent tool. |
| [architecture 0002](docs/decisions/architecture/0002-delivery-as-pinned-copies.md) | Projects get committed copies of the factory's parts and pin the version. |
| [process 0001](docs/decisions/process/0001-approval-through-a-protected-environment.md) | A person approves each pull request on a protected GitHub environment. |

Decision records use the [MADR](https://adr.github.io/madr/) format, with one folder per category: `architecture` for architecture decisions, `process` for process decisions that had real alternatives. A record on `main` is in force: the engineer accepts it by merging its pull request. The table lists only records in force. Superseded and deprecated records stay in their folder.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate and approvals: rules 16 and 31
- no private data in public repos: rule 26
- order of work, test writer, protected paths, agent token, checks from `main`, feedback to the factory: rules 27–35
- existing tools before building, and the engineer's ownership of the domain model: rules 39 and 40
- docs show the current state and follow their format: rule 36
- pull request descriptions: rule 37

## Open

- Review of ADRs 0001 and 0002
- Setup before the first pull request: agent token, environment `approval`, approval workflow, ruleset for `main`. Steps in [docs/setup/](docs/setup/README.md).
- CI pull request after increment 02: workflows `gate` and `test`, `.factory/approval-free-paths`, required checks `checks` and `test`

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note, review, evidence and retro.
- [docs/backlog.md](docs/backlog.md): parts that are built when a project needs them.
- [docs/setup/](docs/setup/README.md): the engineer's one-time setup per repo.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
- `cmd/factory/` and `internal/`: the CLI and its tool-neutral checks. `adapters/claude-code/`: hooks, settings and agent definitions for Claude Code.
- `core/skills/`: the skills of the workflow, with their templates.
- `test/acceptance/`: one folder of `testscript` files per increment.
