# Principles

Each rule names its source, what it depends on, how it is enforced, and when to check it again.

"Depends on" has four values:

- **model**: true for the current models, may change with the next one
- **tool**: tied to a tool or platform feature
- **human**: about how people learn, review and keep an overview
- **evidence**: backed by studies or field reports that do not depend on one model

"Enforced by" names the check or hook that makes the rule hold. "planned (increment nn)" means the check does not exist yet. "not enforced" means only the text and the engineer's review hold the rule. Each retro tries to turn "not enforced" rules into checks (rule 4).

Sources are listed in [research/2026-10-factory-basis.md](research/2026-10-factory-basis.md).

| # | Rule | Source | Depends on | Enforced by | Re-check when |
|---|---|---|---|---|---|
| 1 | Verification and a deterministic pipeline are the centre. Skills are a thin layer for judgment. | StrongDM, Anthropic C compiler, Stripe Minions, Böckeler | evidence | not enforced | a factory with a different centre shows better results |
| 2 | Work that repeats becomes a script. Intelligence is called only for judgment, with typed inputs and outputs. | Jacob, Syntax #1049 (advocacy) | evidence | not enforced | a script costs more to keep than the agent run it replaces |
| 3 | A factory fits one product. Add a rule only after an observed failure. | Jacob; Hashimoto | human | not enforced | every retro |
| 4 | A correction made twice becomes a check. | Hashimoto; `health` plan | human | not enforced | every retro |
| 5 | Computational checks run before inferential checks. | Böckeler, Harness engineering | evidence | planned (increment 04): order in `run-checks` | never, unless inferential checks become cheap and deterministic |
| 6 | Acceptance scenarios are black-box, approved by the engineer, and locked. | Böckeler ("approved scenarios"); Jacob; StrongDM | evidence | planned (increment 04): lock and workflow gate | never |
| 7 | Part of the scenarios is hidden from the implementing agent (holdout). | StrongDM; EvilGenie | model | planned (increment 04): holdout fetch in CI | agents stop gaming visible tests in measured runs |
| 8 | The engineer states properties. The agent writes only the generators. | Anthropic property-based testing | model | planned (increment 04): lock covers property statements | agents state correct properties for implicit assumptions |
| 9 | Agent-written unit tests are scaffolding, never evidence. | VibeCheck; arXiv 2602.07900 | model | planned (increment 04): receipts classify evidence | a new study with current models shows otherwise |
| 10 | Agents get no test-first or red-green instruction. | Böckeler, TDD in the agent loop (Sonnet 4.6, Opus 4.8) | model | not enforced | a model change, or a larger study |
| 11 | Tests count when they are independent of the implementation. | arXiv 2607.05139; StrongDM; Jain | evidence | planned (increment 04): workflow gate (rule 27) | never |
| 12 | A reviewer in a fresh context checks the work. Self-review is not evidence. | Anthropic, Harness design | model | not enforced | the model passes independent checks without a reviewer for 5 increments |
| 13 | Inferential reviews run once per increment, as a final gate. | Jain, Are harnesses worth the cost? | model | planned (increment 04): workflow checkpoint "review done" | a weaker or cheaper model is used, or billing moves to the API |
| 14 | Every check writes a receipt with the number of tests it ran. A missing receipt or a drop blocks. | Böckeler (open problem: a sensor that never fires); sabotage tests | evidence | planned (increment 04): receipts | never |
| 15 | Check output is short and grep-able. Details go to log files. | Anthropic C compiler | model | not enforced | context windows grow so that long output no longer hurts |
| 16 | The merge gate runs on the server: a GitHub ruleset with PRs only, required checks and an empty bypass list. Git hooks and agent hooks give early feedback only. | git (`--no-verify`); GitHub docs; adversarial review | tool | planned (increment 03): ruleset, active after the engineer imports it | the hosting platform changes, or the repos become private (rulesets then need GitHub Pro) |
| 17 | Architecture rules are deterministic first. Review covers only what rules cannot express. A finding seen twice becomes a rule. | Böckeler; Jacob | evidence | not enforced | never |
| 18 | Agents draw only small sketches (about 7 elements). Building-block views come from a model the engineer owns. | Simon Brown; arXiv 2608.28403; project experience | model | not enforced | LLMs reliably get relations right in a measured test |
| 19 | Before merge, the agent asks 2–3 questions, at least one about debugging. The engineer answers before the agent explains. | Anthropic RCT (largest gap in debugging) | human | planned (increment 04): workflow checkpoint "comprehension done" | after 5 increments: did the questions find real gaps? |
| 20 | Measure each increment: interventions, defects found after merge, quota used, time to green, files touched. | Böckeler; Anthropic, Demystifying evals | human | planned (increment 06): PR record template | every retro |
| 21 | Record the agent version per increment. Vendor changes to the harness can lower quality. | Anthropic postmortem 2026-04-23 | tool | planned (increment 06): PR record template | never |
| 22 | Every harness part assumes something the model cannot do. On a model change, remove parts one at a time. | Anthropic, Harness design | model | not enforced | every model change |
| 23 | Chaos runs only in disposable copies or isolated environments, with a manifest before acting and an exact log. | Principles of Chaos Engineering | human | planned (increment 05): sabotage tests in temp copies, `run-chaos` manifest | never |
| 24 | No vendor lock-in: standard skill fields, scripts without dependencies, thin adapters. | project requirement; agentskills.io | tool | planned (increment 02): ADR 0001 confirmation check | a cross-tool hook standard appears |
| 25 | The engineer keeps the overview: a one-page README as the entry point, and ADRs only for architecture decisions, in MADR format, short, accepted only by the engineer. Process and policy live here and in executable files, not in ADRs. | project requirement; Fowler, ADR; MADR | human | not enforced | the README no longer fits on one page |
| 26 | Public repos hold no private data. Personal details and holdout scenarios live in private repos. A check fails on listed private paths, and secret scanning with push protection is on. | decision for public repos, 2026-10-26 | human | planned (increment 04): private-paths check; secret scanning after the engineer enables it | the repos become private |
| 27 | Each increment runs in this order: the engineer approves the scenarios, a separate agent writes the acceptance tests, the engineer reviews and locks them, then the implementing agent implements. A report names the kind of each piece of evidence. | incident 2026-10-26; arXiv 2607.05139; Böckeler | human | planned (increment 04): workflow gate and computed evidence; until then the engineer's approval in chat | every retro |
| 28 | Automated acceptance tests are written by a separate agent in a fresh context. It sees only the approved scenarios and the public interface, never the implementation. | engineer decision 2026-10-26; arXiv 2607.05139; StrongDM | evidence | planned (increment 02): guard hook denies the test agent reads of the implementation; planned (increment 04): provenance log from hooks | never |
| 29 | The implementing agent never changes acceptance tests, locks or approved scenarios. Code never detects a test run or hard-codes test cases. A test that looks wrong goes back to the engineer. | engineer decision 2026-10-26; EvilGenie | model | planned (increment 02): test lock with hashes, guard hook, check for test references in code; planned (increment 04): holdout | agents stop gaming tests in measured runs |
| 30 | Each scenario names how it is verified: automated (a test by a separate agent), demonstrated (a separate agent runs and records it, the engineer views the recording), or manual (the engineer). | engineer decision 2026-10-26; Anthropic, Harness design | human | planned (increment 02): design note check | every retro |
