# Overview

This repo is a software factory: the system around coding agents that decides what counts as done, who checks it, and how a failure turns into a better check. It is a Go CLI, `factory`. Projects install it with `factory install`.

- **Engineer:** the human who owns a project.
- **Agent:** a coding agent that works in it, with its own restricted token.
- **Increment:** one unit of work in two pull requests: the design note, then the tests from a separate session with their code, a review from a separate session, demonstration and retro. The engineer approves owned paths, such as design notes, on GitHub, and agents cannot. [AGENTS.md](AGENTS.md) describes the steps.

## Status

v0, increment 06 of 09. `factory status` shows where each increment stands: missing tests, the review result, open manual and demonstrated scenarios, the retro, and the code commit that waits for a review.

| Increment | Content |
|---|---|
| 01 Basis | this page, principles, research notes, ADRs 0001 and 0002 |
| 02 Trust core | Go CLI with checks for skills, design notes, ADRs, tests and pull request paths; an agent account, with the engineer approving owned paths as code owner; a reviewer that a hook starts, and a review by a model of another vendor; skills with evals; guard hook |
| 03 Setup command | `factory setup`: the admin part of the setup of any repo, run by the engineer with their own `gh`: merge settings, agent account, environment `cross-review` with its secret, rulesets with required checks; a dry run |
| 04 Design questions | comprehension questions on the design pull request, from a fresh session, before the engineer approves; after the code, the reviewer explains where the code deviates from the note |
| 05 Tests per increment | `factory check-tests` blocks only the pull requests of the increment whose tests are missing; a note lists the scenarios whose tests pass before the code |
| 06 Gate and status | `factory status`; an owned results file for manual and demonstrated scenarios; review findings with a severity that `factory check-review` reads, and two rounds before the engineer decides; owned checks, adapter and skills; the commands decide which check applies; the other model reviews the tests before the code; coverage of the acceptance run |
| Next | one increment of the project `health`, with a hand copy of the skills and workflows; until then no new rule without an observed failure |
| 07 Install and distribution | `factory install`, which copies into a project everything it needs to work with the factory: the skills, the adapter, workflow templates, a `CODEOWNERS` template, `.factory/` defaults, and the sandbox host list from `.factory/allowed-hosts`; the binary through releases with attestations or through `go install …@vX.Y.Z` (open); install in cloud and CI; for the project's language: an architecture check, a linter with the project's own rules, each with test cases and a message that names the fix (for TypeScript: oxlint, ESLint or ast-grep for own rules, and dependency-cruiser, open), a ratchet on lint warnings and complexity, and acceptance tests in the project's own format; lint configuration and rules, manifests and lockfiles owned by the engineer; a block in the project's `AGENTS.md` that makes `authoring-skills` the skill for skills and lets agents start subagents and sessions on their own |
| 08 Docs for people | what a project gets to show its whole system to people, apart from the design notes for agents, including a generated view of the interface and the scenarios in force; scope open |
| 09 Handover and feedback loop | a coordinator session that starts the test writer after a design note merges and the implementing session after the tests are pushed; retro per increment, which turns a correction seen twice into a lint rule; `factory check-retro`, `factory feedback` (prefilled issue links for the factory repo), version updates, end-to-end test |

## Decisions

| Record | Decision |
|---|---|
| [architecture 0001](docs/decisions/architecture/0001-go-cli-with-thin-adapters.md) | The factory tooling is a Go CLI with thin adapters per agent tool; the rules for skills live only in the validator of the skill `authoring-skills`, which the CLI runs. |
| [architecture 0002](docs/decisions/architecture/0002-delivery-as-pinned-copies.md) | Projects get committed copies of the factory's parts and pin the version. |
| [process 0002](docs/decisions/process/0002-approval-through-code-owners-and-an-agent-account.md) | Agents act as their own GitHub account; the engineer approves owned paths as code owner. |

Decision records use the [MADR](https://adr.github.io/madr/) format, with one folder per category: `architecture` for architecture decisions, `process` for process decisions that had real alternatives. A record on `main` is in force: the engineer accepts it by merging its pull request. The table lists only records in force. Superseded and deprecated records stay in their folder.

Process and policy are rules in [principles.md](docs/principles.md):
- what counts as evidence: rules 6–13
- the merge gate and approvals: rules 16 and 31
- no private data in public repos: rule 26
- order of work, test writer, protected paths, agent token, checks from `main`, feedback to the factory: rules 27–35
- results of manual and demonstrated scenarios, and the state of each increment: rule 43
- docs show the current state and follow their format: rule 36
- pull request descriptions: rule 37
- existing tools, the domain model and the scope of `AGENTS.md`: rules 39–41

## Open

- Code owner review and the required checks in the ruleset `main`: `factory setup --dry-run` shows the difference.
- Results of the manual scenarios 02.10, 02.11 and 02.12, and of the demonstrated scenario 05.8, which replaces 02.20, before increment 07 hands these guards to a project.
- New evidence for the skills on the model of the sessions: about 100 eval runs.

## Where things are

- [docs/principles.md](docs/principles.md): every rule, its source, what it depends on, how it is enforced, and when to re-check it.
- [docs/research/2026-10-factory-basis.md](docs/research/2026-10-factory-basis.md): the sources behind the rules. Reference only, not required reading.
- [docs/research/2026-11-factory-review.md](docs/research/2026-11-factory-review.md): the factory held against current sources on agent factories and skills, with the findings that increment 06 answers. Reference only, not required reading.
- [docs/research/2026-10-domain-design-with-agents.md](docs/research/2026-10-domain-design-with-agents.md): the sources behind the domain model and the architecture checks. Reference only, not required reading.
- [docs/decisions/](docs/decisions/): the ADRs.
- [docs/increments/](docs/increments/): one folder per increment, with its design note, review, evidence and retro.
- [docs/backlog.md](docs/backlog.md): parts that are built when a project needs them.
- [docs/setup/](docs/setup/README.md): the engineer's one-time setup per repo.
- [AGENTS.md](AGENTS.md): the workflow and rules for agents that work on the factory.
