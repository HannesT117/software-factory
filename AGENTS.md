# Working on the software factory

Read `README.md` first. It is the entry point and shows the status. The factory follows its own rules from `docs/principles.md`.

## Workflow: every phase is a slice

1. **Design note.** Write `docs/slices/<nn>-<slug>/design.md`: the goal, what is out of scope, a small sketch if it helps, and the acceptance scenarios with ids `AC-<nn>-<n>`. Stop. The engineer approves the note or changes it.
2. **Acceptance tests.** Write one black-box test per approved scenario under `test/acceptance/`, titled with its id. A black-box test runs a script as its own process and checks only what is visible from outside: files, exit code, output. Stop. The engineer reviews the tests. After approval they are locked.
3. **Implementation.** Write the code until the acceptance tests pass. Do not change an approved scenario or its test. If one looks wrong, stop and ask.
4. **Summary.** At most 10 lines. Wait for approval. Then push.

Do not write implementation code before step 1 and step 2 are approved.

## Evidence

- Evidence is: approved acceptance scenarios with their reviewed tests, and the confirmation checks of accepted ADRs.
- Agent-written unit tests are scaffolding, not evidence. Write them only when they help to build something. Never count them in a report.
- Every report names the kind of each piece of evidence: approved scenario, ADR confirmation, agent test (scaffolding), or not verified.

## Source of truth and push

- The engineer's clone at `~/dev/software-factory` is the source of truth. Before each change, compare GitHub with the local clone. If the engineer works locally without pushing, ask for a push or a copy and build on that state.
- Commit locally. Push only after the engineer has approved everything that is pending. Push only fast-forward, never with force.

## Docs

- Docs state general rules. They contain no "you" and no "I". The roles are engineer and agent.
- ADRs are only for architecture decisions and use MADR. Context and decision drivers are value-neutral facts. Only the engineer sets `accepted`. An accepted ADR changes only editorially, and only with the engineer's approval.
- Process and policy are rules in `docs/principles.md`. Each rule names its source, what it depends on, how it is enforced, and when to re-check it.
- `README.md` stays one page. Update its status in every phase.

## Code

- `core/` is tool-neutral. No file there names a specific agent tool. Tool-specific parts live in `adapters/<tool>/`.
- Scripts are Node `.mjs` files (Node 22 or later) without dependencies.
- Script output is short and grep-able. Each line starts with the script's prefix, for example `factory:`.
- Run `node --test` before every commit.
