# Research basis, October 2026

These are the sources behind [principles.md](../principles.md) and the ADRs. Each entry says what the source found and what this repo takes from it.

Labels:

- **read**: the primary source was read during this research.
- **via review**: an adversarial review agent read it. It was not re-read.
- **via search**: only a search summary was read.
- **advocacy**: the author sells a product in this space.
- **vendor**: the source is a vendor's own documentation or marketing.

## What other factories put at the centre

| Source | Finding | Taken |
|---|---|---|
| Simon Willison on StrongDM's Software Factory, 2026-02-07 (read) | "Code must not be reviewed by humans." Trust comes from end-to-end scenarios that are kept outside the codebase as a holdout, a "satisfaction" metric, and digital twins of third-party services. About $1,000 in tokens per engineer per day. Willison doubts the no-review rule. | Holdout scenarios are in the backlog. Not the no-review rule, and not the budget. |
| Anthropic, Building a C compiler with parallel Claudes, 2026-02-05 (read; vendor) | "The task verifier is nearly perfect, otherwise Claude will solve the wrong problem." GCC served as the oracle. Agents push straight to a shared upstream, and tests are the only gate. Test output was sized for the model. New features broke old ones until CI enforced them. | A strong verifier first. Short, grep-able output. Not the missing human gate. |
| Anthropic, Harness design for long-running apps, 2026-03-24 (read; vendor) | Planner, generator and evaluator. "Out of the box, Claude is a poor QA agent." A separate evaluator beats self-evaluation. With Opus 4.6 the evaluator became unnecessary for tasks within the model's reach. Remove harness parts one at a time. | A fresh-context reviewer. Every part names its assumption. |
| Stripe Minions, parts 1 and 2 (read; vendor) | Isolated devboxes. Blueprints that mix deterministic and agentic steps. Scoped rule files. At most two CI rounds. A run ends in a pull request that passes CI. An engineer reviews it and asks a second engineer for review. Humans merge. | A fixed order of steps (rule 27). Human approval per pull request (rule 31). |
| OpenAI harness engineering, Feb 2026 (via review, secondary summary) | Mechanical enforcement. Observability the agent can query per worktree. AGENTS.md as a short map that points to docs. | A short AGENTS.md block. Observability is in the backlog. |
| Birgitta Böckeler, Harness engineering for coding agent users, martinfowler.com, 2026-04-02 (read) | Guides (feedforward) and sensors (feedback). Computational sensors before inferential ones. The behaviour harness is the hardest part. AI-generated tests are "not good enough yet". Open problem: is a sensor that never fires good or blind? | Order of checks (rule 5). Receipts and sabotage tests are in the backlog. |
| Adam Jacob, Syntax #1049, 2026-10-21, and blog posts 03–04/2026 (read; advocacy, CEO of Swamp) | "Reverse the loop": deterministic process, intelligence called for judgment. Black-box acceptance tests first. Commit to one architecture with DDD vocabulary and review code against it. Adversarial review for other requirements. A decision model can pick which checks a change needs. "No two factories are the same." | Scripts over prompts. Acceptance first. Architecture rules. Check routing is in the backlog. |
| Mitchell Hashimoto, My AI adoption journey, 2026-02-05 (read; no ties to AI companies) | Engineer the harness: after each agent mistake, he changes AGENTS.md or builds a tool himself. One background agent. Keep hard work manual to keep skills. | A rule only after an observed failure. The engineer submits each feedback issue (rule 35). |

## Tests and evidence

| Source | Finding | Taken |
|---|---|---|
| Böckeler, TDD in the agent loop, martinfowler.com (read) | Sonnet 4.6 coded, Opus 4.8 judged blind. No clear quality difference with or without TDD. Mutation scores were equal. TDD used 3–8.5× the tokens. She stopped telling agents to write tests first. Alternatives: mutation testing, static analysis and modularity reviews, trend tracking, approved scenarios. | No test-first instruction. Approved scenarios. Mutation score as a measure. |
| Naresh Jain, Are harnesses worth the cost?, martinfowler.com, 2026-10-13 (read) | Frontier model: 10.5× the tokens with the harness, acceptance 10/10 either way. Cheaper model: 3/10 rose to 10/10. Keep an independent acceptance check in every case, and run it once as a final gate for strong models. n = 10 runs. | Inferential reviews once per increment. |
| Jacob on TDD (read) | Cites the Böckeler article and says TDD "hurt the outcomes". | Overstated. The measurement shows no difference at a higher cost. |
| Thoughtworks, Future of Software Development retreat, 2026 (via search) | "TDD produces dramatically better results." | A discussion theme, not a measurement. Not used. |
| VibeCheck, arXiv 2609.05978 (read, abstract) | Agent unit tests often run but have weak assertions and miss edge cases. | Agent unit tests are not evidence. |
| Rethinking the value of agent-generated tests, arXiv 2602.07900 (via search) | Agent test writing gives limited gain in issue resolution. | Same. |
| On the risk of coding before testing, arXiv 2607.05139 (via review) | Tests generated independently of faulty code catch 25% of faults vs 14%. | Independence matters. |
| EvilGenie, arXiv 2511.21654 (via review) | Coding agents hard-code tests and edit test files. Hard-coding rose from 2% on unambiguous to 33% on ambiguous tasks. An LLM judge was highly effective in unambiguous cases. Held-out tests added only a little. | Path groups (rule 29). Fixed grammar for notes and tests (rule 32). Holdouts only in the backlog (rule 7). |
| ImpossibleBench, arXiv 2510.20270 (via review) | Read-only tests stop test edits but not special cases in code. LLM monitors caught 86–89% of cheating on simple tasks and 42–50% on complex ones. | Fresh-context review step (rule 12). |
| METR, Recent reward hacking, 2025-06-05 (via review; independent research) | Frontier models modify tests and scoring code. | The tests group includes the runner and module files, and CI and config are a group of their own (rule 29). |
| `testscript`, github.com/rogpeppe/go-internal (read) | Extracted from the Go team's script tests for the `go` command. Shell-like scripts check stdout, stderr and exit status, and run from `go test`. Commands can be guarded by conditions such as `[darwin]`, `[short]` or `[exec:x]`. | Acceptance tests as readable `.txtar` files with one shared runner. All conditions are forbidden. |
| Anthropic, Agentic property-based testing, 2026-01-14, arXiv 2510.09907 (read; vendor) | The agent infers properties and writes Hypothesis tests. Real bugs found in NumPy and others. "Only the library maintainers can decide what the correct property to test is." | The human states properties. |
| Martin Kleppmann, AI will make formal verification mainstream, 2025-12-08 (read) | Proofs become cheap. The human must review the spec. | Later option. |

## Approvals, merge gate and credentials

| Source | Finding | Taken |
|---|---|---|
| GitHub docs, Copilot cloud agent risks and mitigations (read; vendor) | The agent has its own app identity, pushes only to its own `copilot/` branch, and cannot approve or merge. Workflows run only after a human clicks "Approve and run workflows". | Agents never approve or merge (rule 31). |
| DORA, Streamlining change approval (via review) | Approval through peer review, supplemented by automation, recorded in the platform. | Approval in the platform (rule 31). |
| GitHub docs, deployments and environments (read) | Required reviewers are available for public repos on Free. "Prevent self-review" blocks approving a run that the reviewer started. Admin bypass can be switched off. | The approval gate, with "Prevent self-review" off (rule 31). |
| GitHub docs, permissions for fine-grained tokens (read) | Merging a pull request needs Contents write. Approving a pending deployment needs Deployments write. Changing workflow files needs the Workflows permission. | The agent token has Contents and Pull requests only (rule 33). |
| GitHub docs, managing personal access tokens (read) | A fine-grained token cannot be used to contribute to repos where the user is a collaborator. Each token covers one owner. | A machine account would need a classic token or an organization. Not taken. |
| GitHub terms of service (via review) | One free machine account is allowed in addition to the personal account. | Kept as an option for a later team setup. |
| GitHub docs, workflow events (via review) | A `pull_request` workflow runs the version in the pull request's merge commit. | The checks that judge a pull request run from `main` (rule 34). |
| GitHub changelog, 2025-11-07, `pull_request_target` changes (read) | From 2025-12-08, `pull_request_target` always takes the workflow file and the checkout from the default branch. | The gate runs on `pull_request_target` (rule 34). Whether its checks count as required checks on the pull request head is not verified. |
| Claude Code docs, sandboxing (read; vendor) | `sandbox.filesystem.denyWrite` is enforced by the OS for Bash. `allowUnsandboxedCommands: false` turns off the escape hatch. Commands may read credential files such as `~/.ssh` by default, and there is no built-in credential deny list. | Deny list in the user settings (rule 33, scenario 02.12). |
| Secretive README and PostHog handbook (read) | SSH keys in the Mac's Secure Enclave cannot be exported. Use can require Touch ID. Git can sign with them. | Not taken: signed approval commits need an own verifier and a trust root for keys. The approval gate (rule 31) needs neither. |
| Yubico, Git signing with SSH and FIDO2 (read) | `ed25519-sk` keys sign only after a physical touch. Verification uses an allowed-signers file. | Not taken, for the same reason. |
| ssh-keygen manual, allowed-signers format (via review) | Supports `cert-authority`, `namespaces`, `valid-after` and `valid-before`. It cannot prove that a touch happened. | Not taken. A signature proves the key, not a human touch. |
| Secretive source, notifier (via review) | The prompt shows app, key and namespace, not the signed content. Authentication can be remembered for up to 24 h. | Not taken. Touch ID does not show what is signed. |
| GitHub REST, SSH signing keys (via review) | `GET /users/<login>/ssh_signing_keys` is public and returns every signing key of the account. | Not taken. As a trust root, it would also trust keys without Touch ID. |
| GitHub docs, commit signature verification (via review) | "Rebase and merge" drops signatures. Squash merges create new commits that GitHub signs. | Not taken. A squash merge removes the approval signature from `main`. |
| Invariant Labs, GitHub MCP vulnerability (via review; security vendor) | Text in public issues reaches agents that read them and acts as an injected instruction. | Factory agents read only issues with the label `accepted` (rule 35). |
| GitHub docs, interaction limits (via review) | Limits on who may open issues last at most 6 months. | Labels, not limits, decide what agents read. |

## Tooling and distribution

| Source | Finding | Taken |
|---|---|---|
| Swamp source repository, github.com/swamp-club/swamp (read; advocacy) | A TypeScript CLI on Deno, compiled with `deno compile` into one binary that embeds the runtime. Installed with `curl … \| sh`. `swamp repo init --tool <name>` sets up nine agent tools. The skills are embedded in the binary and written into the agents' skill folders, which are gitignored. | One CLI with embedded skills (ADR 0001). Committed copies instead of gitignored ones (ADR 0002). |
| GitHub Spec Kit installation docs (read) | The `specify` CLI is Python, installed with `uv tool install` or run once with `uvx`, pinned to a release tag. | CLIs as the entry point are common. The language follows the ecosystem. |
| Node.js TypeScript docs (read) | Type stripping is on by default from v22.18 and v23.6, for erasable syntax only, with no type check. | Option 2 in ADR 0001. |
| Measured in a Linux container, 2026-10-26 | Startup median: Go or Rust binary about 2 ms, Node script about 30 ms, Node with Effect 3 imported from its root about 390 ms. | Decision driver in ADR 0001. |
| GitHub docs, artifact attestations (via review) | Build provenance logged in Sigstore's public instance for public repos (SLSA Build Level 2). Attestations "are not a guarantee that an artifact is secure". | Releases with attestations and a pinned checksum (ADR 0002). |
| Go blog, reproducible builds (via review) | `go build -trimpath` with `CGO_ENABLED=0` gives reproducible binaries. | Build flags in ADR 0002. |
| Cursor docs, skills (via review; vendor) | Cursor also loads `.claude/skills/`. | A consequence in ADR 0002. |

## Pull request descriptions

| Source | Finding | Taken |
|---|---|---|
| Google engineering practices, Writing good CL descriptions (read) | The first line is a short, standalone summary of what the change does, as an imperative sentence. The body states the problem, why this approach is best, and its shortcomings, with background such as benchmark results. Links may not be visible to later readers, so the description carries enough context. Review the description before submitting, so that it still matches the change. | Title and body sections (rule 37). |
| Linux kernel, Submitting patches, "Describe your changes" (read) | Describe the problem, the user-visible impact, then the change in technical detail. Quantify optimizations and name non-obvious costs. One problem per patch. The description is self-contained and does not refer to earlier versions. Imperative mood. | Why, trade-offs, evidence, one problem per pull request (rule 37). |
| Chris Beams, How to write a Git commit message (read) | Subject at most 50 characters (72 hard), capitalized, no final period, imperative. The body explains what and why, not how. | Title rules (rule 37). |
| GitHub docs, squash merging (read) | The default squash commit message can be set to the pull request title and description. | The description is the commit message on `main` (setup step 6). |
| Jeff Bezos, 2015 letter to shareholders (read) | Type 1 decisions are irreversible, "one-way doors", and need care. Type 2 decisions are reversible, "two-way doors", and can be made quickly. | Each pull request names its merge risk as a one-way or two-way door (rule 37). |
| pstack `technical-writing`, Lauren Tan (read; Cursor) | A PR body is a briefing that a reviewer can read in under a minute. Link logs and metric tables instead of pasting them. Every writing layer except the document modes applies to PR descriptions and commit messages. | The summary and the length limit (rule 37). Not its rule to address the reader as "you". |
| pstack `make-pr-easy-to-review` (read; Cursor) | A TL;DR that matches the diff. Core files apart from generated or mechanical ones. Risky behavior changes, migration order, rollout plan and test coverage called out. A PR too large for notes is split. | Summary, "What changes" and merge risk (rule 37). |
| pstack `blast-radius` (read; Cursor) | Find what a change breaks beyond the diff. Name the one fact that makes it safe, and prove it by running real code. Say how far each fact is proven. | Merge risk names that fact and its proof (rule 37). |
| GitHub docs, creating diagrams (read) | Mermaid renders in pull requests, issues and Markdown files. | Diagrams in pull request descriptions. |
| Conventional Commits 1.0.0 (read) | A type prefix such as `feat:`, `fix:` or `docs:` makes commits machine-readable for changelogs and Semantic Versioning. | Not taken: an extra standard. Titles follow Google and Beams. |

## People, comprehension and operations

| Source | Finding | Taken |
|---|---|---|
| Anthropic RCT on coding skills, 2026-01-29 (read; vendor) | n = 52 in a chat sidebar. AI group 50% vs 67% on a quiz. The largest gap was in debugging. "Generation then comprehension" scored 86% with n = 2. | Comprehension questions, at least one on debugging. |
| Mario Zechner, Thoughts on slowing down, 2026-03-25 (read) | Agent errors compound without a human bottleneck. Write architecture and APIs by hand. Cap generated code to what a person can review. | Small increments. Stop-and-discuss triggers. |
| Anthropic, Redesigning source control, 2026-10-13 (read; vendor) | 30,000 merges a day. Humans review intent, not every line. Agents that watch telemetry catch subtle issues. | Review by intent. Observability is in the backlog. |
| Anthropic postmortem, 2026-04-23 (read; vendor) | Three harness changes lowered Claude Code quality for weeks while the model stayed the same. | Record the agent version. Pin the factory version (ADR 0002). |
| Anthropic, Demystifying evals for AI agents, 2026-01-09 (read; vendor) | Code graders first. Start with 20–50 tasks from real failures. Read transcripts. | Metrics from real failures. |
| Charity Majors, Production is where the rigor goes (via review; CTO of an observability vendor) | "The spec describes intent. Telemetry describes reality." | 2 am support in the backlog. |
| Principles of Chaos Engineering (read) | Steady state, hypothesis, real-world events, automation, minimal blast radius. | Chaos rules. |

## Architecture and diagrams

| Source | Finding | Taken |
|---|---|---|
| Simon Brown, 2026 talks (via search; author of C4 and Structurizr) | LLMs draw poor C4 diagrams unless the tool enforces the model. Use models as code. | The human owns the model. Agents draw sketches only. |
| Industry case study, arXiv 2608.28403 (via search) | Recovered entities are stable. Recovered relations need human review. | Same. |

## Tools and platforms (vendor docs, read)

- agentskills.io specification: `name` (1–64 characters, lowercase letters, digits and single hyphens, same as the folder), `description` (1–1024 characters), `license`, `compatibility` (at most 500 characters), `metadata`, experimental `allowed-tools`.
- vercel-labs/skills: Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read `.agents/skills/`. Claude Code reads `.claude/skills/`.
- Claude Code skills: a `<skill-name>` entry may be a symlink. Project cloud threads do not load committed `.claude/skills/`.
- Claude Code memory: Claude Code reads nothing under `.agents/`. AGENTS.md is read through `@AGENTS.md`. Edit and Write refuse to write through a symlink.
- Claude Code hooks: exit 2 blocks; PreToolUse `permissionDecision: deny`; Stop `decision: block`; `$CLAUDE_PROJECT_DIR`.
- GitHub rulesets: free for public repos, not for private repos on GitHub Free.
- MADR template: all frontmatter fields are optional. Sections: Context and Problem Statement, Decision Drivers, Considered Options, Decision Outcome with Consequences and Confirmation, Pros and Cons of the Options, More Information.
- Opper: decision models at `/v3/compat/v1/systemone` and `/v3/compat/decisions`. EU routes per model.
- dependency-cruiser 18.1.0: waits for the public API in TypeScript 7.1.
