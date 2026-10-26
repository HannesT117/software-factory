# Research basis, October 2026

These are the sources behind [principles.md](../principles.md) and the ADRs. Each entry says what the source found and what this repo takes from it.

Labels:

- **read**: the primary source was read during this research.
- **via review**: an adversarial review agent read it. It was not re-read.
- **advocacy**: the author sells a product in this space.
- **vendor**: the source is a vendor's own documentation or marketing.

## What other factories put at the centre

| Source | Finding | Taken |
|---|---|---|
| Simon Willison on StrongDM's Software Factory, 2026-02-07 (read) | Code is neither written nor reviewed by humans. Trust comes from end-to-end scenarios that are kept outside the codebase as a holdout, a "satisfaction" metric, and digital twins of third-party services. About $1,000 in tokens per engineer per day. Willison doubts the no-review rule. | Holdout scenarios. Not the no-review rule, and not the budget. |
| Anthropic, Building a C compiler with parallel Claudes, 2026-02-05 (read) | "The task verifier is nearly perfect, otherwise Claude will solve the wrong problem." GCC served as the oracle. Test output was sized for the model. New features broke old ones until CI enforced them. | A strong verifier first. Short, grep-able output. |
| Anthropic, Harness design for long-running apps, 2026-03-24 (read) | Planner, generator and evaluator. "Out of the box, Claude is a poor QA agent." A separate evaluator beats self-evaluation. With Opus 4.6 the evaluator became unnecessary for tasks within the model's reach. Remove harness parts one at a time. | A fresh-context reviewer. Every part names its assumption. |
| Stripe Minions part 2 (via review) | Isolated devboxes, blueprints that mix deterministic and agentic steps, scoped rule files, at most two CI rounds, then human review. | Workflow as data. |
| OpenAI harness engineering, Feb 2026 (via review, secondary summary) | Mechanical enforcement. Observability the agent can query per worktree. AGENTS.md as a short map that points to docs. | A short AGENTS.md block. Observability is in the backlog. |
| Birgitta Böckeler, Harness engineering for coding agent users, martinfowler.com, 2026-04-02 (read) | Guides (feedforward) and sensors (feedback). Computational sensors before inferential ones. The behaviour harness is the hardest part. AI-generated tests are "not good enough yet". Open problem: is a sensor that never fires good or blind? | Order of checks. Receipts and sabotage tests. |
| Adam Jacob, Syntax #1049, 2026-10-21, and blog posts 03–04/2026 (read; advocacy, CEO of Swamp) | "Reverse the loop": deterministic process, intelligence called for judgment. Black-box acceptance tests first. Commit to one architecture with DDD vocabulary and review code against it. Adversarial review for other requirements. A decision model can pick which checks a change needs. "No two factories are the same." | Scripts over prompts. Acceptance first. Architecture rules. Check routing is in the backlog. |
| Mitchell Hashimoto, My AI adoption journey, 2026-02-05 (read; no ties to AI companies) | Engineer the harness: every agent mistake becomes an AGENTS.md line or a tool. One background agent. Keep hard work manual to keep skills. | A rule only after an observed failure. |

## Tests and evidence

| Source | Finding | Taken |
|---|---|---|
| Böckeler, TDD in the agent loop, martinfowler.com (read) | Sonnet 4.6 coded, Opus 4.8 judged blind. No clear quality difference with or without TDD. Mutation scores were equal. TDD used 3–8.5× the tokens. She stopped telling agents to write tests first. Alternatives: mutation testing, static analysis and modularity reviews, trend tracking, approved scenarios. | No test-first instruction. Approved scenarios. Mutation score as a measure. |
| Naresh Jain, Are harnesses worth the cost?, martinfowler.com, 2026-10-13 (read) | Frontier model: 10.5× the tokens with the harness, acceptance 10/10 either way. Cheaper model: 3/10 rose to 10/10. Keep an independent acceptance check in every case, and run it once as a final gate for strong models. n = 10 runs. | Inferential reviews once per slice. |
| Jacob on TDD (read) | Cites the Böckeler article and says TDD "hurt the outcomes". | Overstated. The measurement shows no difference at a higher cost. |
| Thoughtworks, Future of Software Development retreat, 2026 (via search summary) | "TDD produces dramatically better results." | A discussion theme, not a measurement. Not used. |
| VibeCheck, arXiv 2609.05978 (read, abstract) | Agent unit tests often run but have weak assertions and miss edge cases. | Agent unit tests are not evidence. |
| Rethinking the value of agent-generated tests, arXiv 2602.07900 (via search) | Agent test writing gives limited gain in issue resolution. | Same. |
| On the risk of coding before testing, arXiv 2607.05139 (via review) | Tests generated independently of faulty code catch 25% of faults vs 14%. | Independence matters. |
| EvilGenie, arXiv 2511.21654 (via search) | Coding agents hard-code tests and edit test files, much more often on ambiguous tasks. | Lock and holdout. |
| Anthropic, Agentic property-based testing, 2026-01-14, arXiv 2510.09907 (read) | The agent infers properties and writes Hypothesis tests. Real bugs found in NumPy and others. "Only the library maintainers can decide what the correct property to test is." | The human states properties. |
| Martin Kleppmann, AI will make formal verification mainstream, 2025-12-08 (read) | Proofs become cheap. The human must review the spec. | Later option. |

## People, comprehension and operations

| Source | Finding | Taken |
|---|---|---|
| Anthropic RCT on coding skills, 2026-01-29 (read) | n = 52 in a chat sidebar. AI group 50% vs 67% on a quiz. The largest gap was in debugging. "Generation then comprehension" scored 86% with n = 2. | Comprehension questions, at least one on debugging. |
| Mario Zechner, Thoughts on slowing down, 2026-03-25 (read) | Agent errors compound without a human bottleneck. Write architecture and APIs by hand. Cap generated code to what a person can review. | Small slices. Stop-and-discuss triggers. |
| Anthropic, Redesigning source control, 2026-10-13 (read) | 30,000 merges a day. Humans review intent, not every line. Agents that watch telemetry catch subtle issues. | Review by intent. Observability is in the backlog. |
| Anthropic postmortem, 2026-04-23 (read) | Three harness changes lowered Claude Code quality for weeks while the model stayed the same. | Record the agent version. |
| Anthropic, Demystifying evals for AI agents, 2026-01-09 (read) | Code graders first. Start with 20–50 tasks from real failures. Read transcripts. | Metrics from real failures. |
| Charity Majors, Production is where the rigor goes (via review; CTO of an observability vendor) | "The spec describes intent. Telemetry describes reality." | 2 am support in the backlog. |
| Principles of Chaos Engineering (read) | Steady state, hypothesis, real-world events, automation, minimal blast radius. | Chaos rules. |

## Architecture and diagrams

| Source | Finding | Taken |
|---|---|---|
| Simon Brown, 2026 talks (via search; author of C4 and Structurizr) | LLMs draw poor C4 diagrams unless the tool enforces the model. Use models as code. | The human owns the model. Agents draw sketches only. |
| Industry case study, arXiv 2608.28403 (via search) | Recovered entities are stable. Recovered relations need human review. | Same. |

## Tools and platforms (vendor docs, read)

- agentskills.io specification: `name` (lowercase, hyphens, same as the folder), `description` (at most 1024 characters), `license`, `compatibility`, `metadata`, experimental `allowed-tools`.
- vercel-labs/skills: Codex, Cursor, Copilot, Gemini CLI, OpenCode and Pi read `.agents/skills/`. Claude Code reads `.claude/skills/`.
- Claude Code skills: a `<skill-name>` entry may be a symlink. Project cloud threads do not load committed `.claude/skills/`.
- Claude Code memory: Claude Code reads nothing under `.agents/`. AGENTS.md is read through `@AGENTS.md`. Edit and Write refuse to write through a symlink.
- Claude Code hooks: exit 2 blocks; PreToolUse `permissionDecision: deny`; Stop `decision: block`; `$CLAUDE_PROJECT_DIR`.
- Claude Code sandboxing: `sandbox.filesystem.denyWrite` is enforced by the OS for Bash. `allowUnsandboxedCommands: false` turns off the escape hatch. Edit and Write follow permission rules instead.
- GitHub rulesets: free for public repos, not for private repos on GitHub Free.
- Opper: decision models at `/v3/compat/v1/systemone` and `/v3/compat/decisions`. EU routes per model.
- dependency-cruiser 18.1.0: waits for the public API in TypeScript 7.1.

## Corrections made during the research

- An early claim that Thoughtworks found "the opposite of Jacob" on TDD was wrong. See the TDD rows above.
- "Test-first ordering raises fault detection by 8–18%" was a misreading. The paper's finding is about independence (25% vs 14%).
