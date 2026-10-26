# Working on the software factory

Read `README.md` first. It is the entry point and shows the status. The factory follows its own rules from `docs/principles.md`.

## Workflow: every increment follows these steps

An increment is one unit of work. Each increment has a folder `docs/increments/<nn>-<slug>/`.

1. **Design note.** Write `design.md`: the goal, what is out of scope, a sketch only if it answers a question, and the scenarios with ids `scenario <nn>.<n>`. Each scenario names its verification kind. Stop. The engineer sets `Status: approved` or asks for changes.
2. **Acceptance tests.** A separate agent with a fresh context writes the tests for the automated scenarios, in `test/acceptance/<nn>-<slug>.test.mjs`. Each test title starts with `scenario <nn>.<n>:`. This agent sees only the approved scenarios and the public interface, never the implementation. Stop. The engineer reviews the tests and locks them.
3. **Implementation.** The implementing agent writes the code until the tests pass. Stop and ask if a test looks wrong.
4. **Demonstration.** A separate agent runs the demonstrated scenarios like a user and records them. The engineer views the recordings and runs the manual scenarios.
5. **Summary.** At most 10 lines. Wait for approval. Then push.

Do not write implementation code before steps 1 and 2 are approved.

## Verification kinds

- **automated:** a black-box test, written in step 2 by a separate agent. It runs a script as its own process and checks only files, exit code and output.
- **demonstrated:** a separate agent runs the scenario and records it. The engineer views the recording.
- **manual:** the engineer runs the scenario.

## The implementing agent never

- changes an acceptance test, a `tests.lock` or an approved scenario
- writes code that detects a test run, or that hard-codes test cases
- counts its own unit tests as evidence. They are scaffolding.

## Reports

Every report names the kind of each piece of evidence: approved scenario, demonstrated, manual, ADR confirmation, agent test (scaffolding), or not verified.

## Source of truth and push

- The engineer's clone at `~/dev/software-factory` is the source of truth. Before each change, compare GitHub with the local clone. If the engineer works locally without pushing, ask for a push or a copy and build on that state.
- Commit locally. Push only after the engineer has approved everything that is pending. Push only fast-forward, never with force.

## Docs

- Docs state general rules. They contain no "you" and no "I". The roles are engineer and agent.
- ADRs are only for architecture decisions and use MADR. Context and decision drivers are value-neutral facts. Only the engineer sets `accepted`. An accepted ADR changes only editorially, and only with the engineer's approval.
- Process and policy are rules in `docs/principles.md`. Each rule names its source, what it depends on, how it is enforced, and when to re-check it.
- `README.md` stays one page. Update its status in every increment.

## Code

- `core/` is tool-neutral. No file there names a specific agent tool. Tool-specific parts live in `adapters/<tool>/`.
- Scripts are Node `.mjs` files (Node 22 or later) without dependencies.
- Script output is short and grep-able. Each line starts with the script's prefix, for example `factory:`.
- Skill names use the -ing form, for example `designing-increments`.
- Run `node --test` before every commit.
