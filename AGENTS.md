# Working on the software factory

Read `README.md` first. It is the entry point and shows the status. The factory follows its own rules from `docs/principles.md`.

## Workflow: every increment follows these steps

An increment is one unit of work in three pull requests. Each increment has a folder `docs/increments/<nn>-<slug>/`. The format of a design note and of an acceptance test, and the path groups, are fixed in `docs/increments/02-trust-core/design.md`, section Format.

1. **Design pull request.** Write `design.md` with goal, scope, interface and scenarios. Each scenario names its verification kind. `factory check-design` must pass. The pull request changes only the design note. Turn on auto-merge. Stop. The engineer approves, and GitHub merges once all checks pass.
2. **Tests pull request.** The engineer starts a separate session as `acceptance-test-writer`, with its fixed prompt and the skill `writing-acceptance-tests`. It writes one `testscript` file per automated scenario and sees only the design note on `main` and the public interface. `factory check-tests` must pass. The pull request changes only the tests group. It turns on auto-merge. Stop. The engineer approves, and GitHub merges once all checks pass.
3. **Implementation pull request.** The implementing agent writes the code until the tests pass. The pull request changes no path in a group. Stop and ask if a test looks wrong.
4. **Review and demonstration.** The engineer starts a separate session as `increment-reviewer`, with the skill `reviewing-increments`. It reads the design note and the diff against `main`. It writes `review.md` into the implementation pull request: findings on special cases, test detection, work outside the scenarios, and changed docs that contain their own history or elements outside their format, then 2 or 3 comprehension questions, at least one about debugging. The engineer answers there before the reviewer explains. The same session runs the demonstrated scenarios and puts their recordings into `evidence/`.
5. **Manual scenarios.** The engineer runs the manual scenarios.
6. **Retro and summary.** Write `retro.md`: metrics, agent version, corrections, and every point that concerns the factory. From increment 04 on, `factory feedback` prints a prefilled issue link for each factory point, and the engineer submits it. Then a summary of at most 10 lines. Turn on auto-merge. Stop. The engineer approves, and GitHub squash-merges once all checks pass.

Do not write implementation code before the design and tests pull requests are merged.

## Verification kinds

- **automated:** a `testscript` file from step 2
- **demonstrated:** the review session runs it and records it, and the engineer views the recording
- **manual:** the engineer runs it

## The implementing agent never

- changes a path in a group: acceptance tests, the test runner, `go.mod`, `go.sum`, anything under `.github/` or `.factory/`, or a design note on `main`
- writes code that detects a test run, or that hard-codes test cases
- uses the engineer's own credentials, or reads `~/.ssh`, `~/.config/gh` or the SSH agent. Agents act as the machine account.
- approves a pull request. It may merge one, or turn on auto-merge, once GitHub allows the merge: all required checks pass, including the engineer's approval. It does not merge while the engineer has switched off the ruleset.
- counts its own unit tests as evidence. They are scaffolding.
- reads an issue in the factory repo that does not have the label `accepted`

## Reports

Every report starts with one line for each pull request that the agent opened or changed since the last report: its title and its link. Before building a change, the agent weighs its benefit against its cost. It drops a change with little benefit, also one that a review suggested, and names what it dropped in the report. Every report names the kind of each piece of evidence: approved scenario, demonstrated, manual, review, ADR confirmation, agent test (scaffolding), or not verified.

The agent asks the engineer directly, never through a section in a document. A question says what it is about, so that a reader without context understands it, and names the consequences of each option. The agent decides implementation details itself and does not ask about them.

## Pull requests

With squash merges, the title and description of a pull request become the commit message on `main`. The description gives a reader without context an overview of the change.

- **Title:** a complete sentence in the imperative, at most 50 characters (72 at most), capitalized, without a final period. It names what the change does, not which files it touches, and it is clear without the description. It contains no internal numbers such as increment or scenario ids. If the title needs "and" to join two topics, split the pull request.
- **Description:** the sections of `.github/pull_request_template.md`, in this order. Delete an optional section that does not apply.
  - **Summary** (required): two or three sentences that a reader without context understands: what changes, and which problem it solves.
  - **Context** (optional): the problem, and what a reader must know to follow the change.
  - **What changes** (required): the change at a high level. Core changes first, then mechanical ones such as renames or generated files.
  - **Diagram** (optional): a Mermaid sketch of about 7 elements that answers one question (rule 18).
  - **Merge risk** (required): one-way door or two-way door; what the change could break beyond the diff; the one fact that makes it safe, and how far that fact is proven; how to undo the change.
  - **Trade-offs** (optional): costs and shortcomings that are not obvious.
  - **Evidence** (required): each piece with its kind (see Reports).
- The description is self-contained and describes the current state. Before the approval, check that it still matches the diff.
- One problem per pull request. If the description gets long, split the pull request.
- Every plan names the target pull request of each change. First check the topics and path groups of the open pull requests. A change that fits an open pull request goes there. A new pull request names why no open one fits.
- Notes for the engineer's approval (what to read, decide or run) go into the first comment of the pull request, not into the description.
- Commits on a work branch use the same title rules. Their body explains what and why, not how.

## Source of truth and push

- The engineer's clone at `~/dev/software-factory` is the source of truth. Before each change, compare GitHub with the local clone. If the engineer works locally without pushing, ask for a push or a copy and build on that state.
- Push the work branch after every commit, fast-forward only, never with force. Never push to `main`. Changes reach `main` only by squash merge of a pull request, and a branch rule enforces this. After the merge, the branch is deleted. A branch whose content moved to other branches is deleted too.

## Docs

- Docs state general rules. They contain no "you" and no "I". The roles are engineer and agent.
- Decision records use MADR, in `docs/decisions/<category>/`. The categories are `architecture` and `process`. A process decision gets a record only when it had real alternatives; its rule in `docs/principles.md` links to the record. Numbers count per category. Context and decision drivers are value-neutral facts. A record has no status field: on `main` it is in force, because the engineer merged it. A record in force is never rewritten. A new record that replaces it, in the same category, adds `status: superseded by ADR-<nnnn>` to the old one, in the same pull request. A record that stops applying without a replacement gets `status: deprecated` and a one-line reason under More Information.
- Docs describe the current state. How a doc came about stays in git and in the pull request. Research notes keep rejected options with their reason, not the course of the research.
- A doc that follows a named format (MADR, the Agent Skills specification, the design-note format) uses only the elements of that format. An addition names its source in the pull request. `factory check-adr` must pass.
- Process and policy are rules in `docs/principles.md`. Each rule names its source, what it depends on, how it is enforced, and when to re-check it.
- `README.md` stays one page. Update its status in every increment.

## Code

- The CLI is Go. It uses only the standard library at runtime. Test-only dependencies are allowed.
- `core/` and the tool-neutral packages name no specific agent tool. Tool-specific parts live in `adapters/<tool>/`.
- Output is short and grep-able. Each line starts with the command's prefix, for example `verify:`.
- Skill names use the -ing form, for example `designing-increments`.
- Run `go test ./...` before every commit.
