# Principles

Each rule names its source, what it depends on, and when to check it again. "Depends on" has four values:

- **model**: true for the current models, may change with the next one
- **tool**: tied to a tool or platform feature
- **human**: about how people learn, review and keep an overview
- **evidence**: backed by studies or field reports that do not depend on one model

Sources are listed in [research/2026-10-factory-basis.md](research/2026-10-factory-basis.md).

| # | Rule | Source | Depends on | Re-check when |
|---|---|---|---|---|
| 1 | Verification and a deterministic pipeline are the centre. Skills are a thin layer for judgment. | StrongDM, Anthropic C compiler, Stripe Minions, Böckeler | evidence | a factory with a different centre shows better results |
| 2 | Work that repeats becomes a script. Intelligence is called only for judgment, with typed inputs and outputs. | Jacob, Syntax #1049 (advocacy) | evidence | a script costs more to keep than the agent run it replaces |
| 3 | A factory fits one product. Add a rule only after an observed failure. | Jacob; Hashimoto | human | every retro |
| 4 | A correction made twice becomes a check. | Hashimoto; `health` plan | human | every retro |
| 5 | Computational checks run before inferential checks. | Böckeler, Harness engineering | evidence | never, unless inferential checks become cheap and deterministic |
| 6 | Acceptance scenarios are black-box, approved by the human, and locked. | Böckeler ("approved scenarios"); Jacob; StrongDM | evidence | never |
| 7 | Part of the scenarios is hidden from the implementing agent (holdout). | StrongDM; EvilGenie | model | agents stop gaming visible tests in measured runs |
| 8 | The human states properties. The agent writes only the generators. | Anthropic property-based testing | model | agents state correct properties for implicit assumptions |
| 9 | Agent-written unit tests are scaffolding, never evidence. | VibeCheck; arXiv 2602.07900 | model | a new study with current models shows otherwise |
| 10 | Agents get no test-first or red-green instruction. | Böckeler, TDD in the agent loop (Sonnet 4.6, Opus 4.8) | model | a model change, or a larger study |
| 11 | Tests count when they are independent of the implementation. | arXiv 2607.05139; StrongDM; Jain | evidence | never |
| 12 | A reviewer in a fresh context checks the work. Self-review is not evidence. | Anthropic, Harness design | model | the model passes independent checks without a reviewer for 5 slices |
| 13 | Inferential reviews run once per slice, as a final gate. | Jain, Are harnesses worth the cost? | model | a weaker or cheaper model is used, or billing moves to the API |
| 14 | Every check writes a receipt with the number of tests it ran. A missing receipt or a drop blocks. | Böckeler (open problem: a sensor that never fires); sabotage tests | evidence | never |
| 15 | Check output is short and grep-able. Details go to log files. | Anthropic C compiler | model | context windows grow so that long output no longer hurts |
| 16 | The merge gate runs on the server: a GitHub ruleset with PRs only, required checks and an empty bypass list. Git hooks and agent hooks give early feedback only. | git (`--no-verify`); GitHub docs; adversarial review | tool | the hosting platform changes, or the repos become private (rulesets then need GitHub Pro) |
| 17 | Architecture rules are deterministic first. Review covers only what rules cannot express. A finding seen twice becomes a rule. | Böckeler; Jacob | evidence | never |
| 18 | Agents draw only small sketches (about 7 elements). Building-block views come from a model the human owns. | Simon Brown; arXiv 2608.28403; the engineer's experience | model | LLMs reliably get relations right in a measured test |
| 19 | Before merge, the agent asks 2–3 questions, at least one about debugging. The engineer answers before the agent explains. | Anthropic RCT (largest gap in debugging) | human | after 5 slices: did the questions find real gaps? |
| 20 | Measure each slice: interventions, defects found after merge, quota used, time to green, files touched. | Böckeler; Anthropic, Demystifying evals | human | every retro |
| 21 | Record the agent version per slice. Vendor changes to the harness can lower quality. | Anthropic postmortem 2026-04-23 | tool | never |
| 22 | Every harness part assumes something the model cannot do. On a model change, remove parts one at a time. | Anthropic, Harness design | model | every model change |
| 23 | Chaos runs only in disposable copies or isolated environments, with a manifest before acting and an exact log. | Principles of Chaos Engineering | human | never |
| 24 | No vendor lock-in: standard skill fields, scripts without dependencies, thin adapters. | engineer's requirement; agentskills.io | tool | a cross-tool hook standard appears |
| 25 | The engineer keeps the overview: a one-page README as the entry point, and ADRs only for architecture decisions, short, accepted only by the engineer. Process and policy live here and in executable files, not in ADRs. | engineer's requirement; Fowler, ADR | human | the README no longer fits on one page |
| 26 | Public repos hold no private data. Personal details and holdout scenarios live in private repos. A check fails on listed private paths, and secret scanning with push protection is on. | engineer's decision for public repos | human | the repos become private |
