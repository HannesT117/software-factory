# Working on the software factory

Read `README.md` first. It is the entry point and shows the status. The factory follows its own rules from `docs/principles.md`.

## Workflow: every increment follows these steps

An increment is one unit of work. Each increment has a folder `docs/increments/<nn>-<slug>/`. The format of a design note and of an acceptance test is fixed in `docs/increments/02-trust-core/design.md`, section Format.

1. **Design note.** Write `design.md` with goal, scope, interface and scenarios. Each scenario names its verification kind. `factory check-design` must pass. Stop. The engineer approves with a signed commit.
2. **Acceptance tests.** A separate session runs the agent `acceptance-test-writer` with its fixed prompt. It writes one `testscript` file per automated scenario. It sees only the approved note and the public interface, never the implementation. `factory check-tests` must pass. Stop. The engineer approves with a signed commit.
3. **Implementation.** The implementing agent writes the code until the tests pass and `factory verify` passes. Stop and ask if a test looks wrong.
4. **Review.** A fresh-context reviewer checks the diff for special cases, test detection and anything outside the scenarios. Then the agent asks 2 or 3 comprehension questions, at least one about debugging. The engineer answers before the agent explains.
5. **Demonstration.** Demonstrated scenarios run in a separate session, and their recordings go to `evidence/`. The engineer runs the manual scenarios.
6. **Summary.** At most 10 lines. Wait for approval. Then push.

Do not write implementation code before steps 1 and 2 are approved.

## Verification kinds

- **automated:** a `testscript` file from step 2
- **demonstrated:** a separate session runs it and records it, and the engineer views the recording
- **manual:** the engineer runs it

## The implementing agent never

- changes an acceptance test, an `approvals.json`, `.factory/allowed_signers` or an approved design note
- writes code that detects a test run, or that hard-codes test cases
- signs commits. Only the engineer signs, and only approvals.
- counts its own unit tests as evidence. They are scaffolding.

## Reports

Every report names the kind of each piece of evidence: approved scenario, demonstrated, manual, review, ADR confirmation, agent test (scaffolding), or not verified.

## Source of truth and push

- The engineer's clone at `~/dev/software-factory` is the source of truth. Before each change, compare GitHub with the local clone. If the engineer works locally without pushing, ask for a push or a copy and build on that state.
- Commit locally. Push only after the engineer has approved everything that is pending. Push only fast-forward, never with force.

## Docs

- Docs state general rules. They contain no "you" and no "I". The roles are engineer and agent.
- ADRs are only for architecture decisions and use MADR. Context and decision drivers are value-neutral facts. Only the engineer sets `accepted`.
- Process and policy are rules in `docs/principles.md`. Each rule names its source, what it depends on, how it is enforced, and when to re-check it.
- `README.md` stays one page. Update its status in every increment.

## Code

- The CLI is Go. It uses only the standard library at runtime. Test-only dependencies are allowed.
- `core/` and the tool-neutral packages name no specific agent tool. Tool-specific parts live in `adapters/<tool>/`.
- Output is short and grep-able. Each line starts with the command's prefix, for example `verify:`.
- Skill names use the -ing form, for example `designing-increments`.
- Run `go test ./...` before every commit.
