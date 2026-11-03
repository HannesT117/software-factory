# Documentation for people, apart from the specs for agents, November 2026

These are the sources behind the choice of how the factory documents the whole system for the engineer. Each entry says what the source found and what this repo can take from it.

Labels:

- **read**: the primary source was read for this note.
- **via search**: only a search summary or a secondary source was read.
- **measured**: a study with numbers. **practitioner**: experience without measurement. **vendor**: the docs of a tool vendor.

## Today in this repo

- The design note of an increment is written for the test writer, the implementing agent and `factory check-design`. It describes a future state, and one diagram answers one question.
- `README.md` is one page with the status. `docs/principles.md` holds the rules. Projects get `docs/domain/` with a glossary and a context map through the skill `designing-increments`; the factory repo itself has no `docs/domain/`.
- No document describes the whole system as it is now, with its flows, for a person.

## Specs for agents and docs for people

| Source | Finding | Can take |
|---|---|---|
| [GitHub Spec Kit, spec-driven.md](https://github.com/github/spec-kit/blob/main/spec-driven.md) (read; vendor) | Per feature, in `specs/<branch>/`: `spec.md` (what and why), `plan.md`, `tasks.md` for a task agent, and a project-wide `constitution.md`. No system-wide doc for people, no diagrams. | Spec per change. Nothing for a system view. |
| [Kiro, specs](https://kiro.dev/docs/specs/) and [steering](https://kiro.dev/docs/steering/) (read; vendor) | Per feature: `requirements.md`, `design.md` with "sequence diagrams and data flow", `tasks.md`. Whole-project context is in steering files (`product.md`, `tech.md`, `structure.md`) that the agent loads in every interaction. No statement that specs stay current after the feature. | Diagrams in the design of a change. Steering files are context for the agent, not a doc for people. |
| [Böckeler, Understanding spec-driven development: Kiro, spec-kit and Tessl](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html) (read; practitioner) | Three levels: spec-first, spec-anchored ("the spec is kept even after the task is complete") and spec-as-source. Spec Kit and Kiro are spec-first. The markdown was "very verbose and tedious to review"; "I'd rather review code than all these markdown files." | A per-change spec is not a current-state doc. More markdown is more review load for the engineer. |
| [OpenSpec, concepts](https://raw.githubusercontent.com/Fission-AI/OpenSpec/main/docs/concepts.md) (read; vendor, open source) | `openspec/specs/` is the "source of truth" for how the system behaves now. Each change has a folder with proposal, delta specs (ADDED, MODIFIED, REMOVED), design and tasks. Deltas merge into the current specs only when the change is archived, after the implementation. | The split of a current-state doc from a change doc, and the merge after the code. Its current state is requirements text, not flows. |
| [AGENTS.md](https://agents.md/) (read) | "README.md files are for humans"; AGENTS.md is "a README for agents" with build steps, tests and conventions. | Already in use. Neither file is a system description with flows. |
| [Diátaxis, explanation](https://diataxis.fr/explanation/) (read) | Explanation is understanding-oriented, takes a "higher and wider" view, and covers the bigger picture, choices and reasons. Reference is "a close-up view of the machinery". Explanation tends to "absorb other things". | A doc for people is explanation. The design note and the Interface are reference for agents and tests. |
| [arc42, runtime view](https://docs.arc42.org/section-6/) (read) | Scenarios for important use cases, external interfaces, operation and errors. Only "a representative selection", chosen by architectural relevance. Sequence or activity diagrams, or numbered steps in text. | A small set of flows, picked by relevance, as diagram or numbered steps. |
| [C4 model, dynamic diagram](https://c4model.com/diagrams/dynamic) and [diagrams](https://c4model.com/diagrams) (read) | A dynamic diagram shows how elements collaborate for one feature; "should be used sparingly". "The system context and container diagrams are sufficient for most software development teams." | One context and one container view, plus a few dynamic views. |
| [Hilton on Martraire's principles](https://hilton.org.uk/blog/martraire-documentation-principles) (read; practitioner) | Documentation is knowledge transfer. Document knowledge with lasting or broad value. "The Default Is Don't." | Only the stable knowledge goes into a hand-written doc. |
| Martraire, Living Documentation, Addison-Wesley 2019, via [workshop description](https://sessionize.com/s/cyrille-martraire/living-documentation/31423) (via search; practitioner) | Separate stable from volatile knowledge. Evergreen documents hold only stable knowledge and need no upkeep. Living diagrams are generated from the code. | Hand-written text for what changes rarely; generated diagrams for what changes often. |

### Rejected

- **Tessl, spec-as-source:** code is generated from the spec and never edited. Per a [2026 review](https://codemyspec.com/blog/tessl-review) (via search) it was not generally available in mid-2026, and Tessl turned to a registry for agent skills in January 2026. The factory keeps code as the source.
- **Steering files as the doc for people:** they are loaded into every agent session and cost context there; the factory already keeps `AGENTS.md` short (rule 41).

## How a system doc stays current

| Source | Finding | Can take |
|---|---|---|
| [Structurizr, Why "as code"](https://docs.structurizr.com/as-code) (read; vendor) | A DSL holds the model; views render it. "AI agents have limited abilities to interact with UI-driven tools, but they can generate text." Drift detection is a suggested agent workflow, not a feature. Exports Mermaid, PlantUML and a static site. Needs Java or a container. | A model that the engineer owns, as in rule 18. A drift check is extra work in any tool. |
| [go-arch-lint](https://github.com/vsfedorenko/go-arch-lint/releases) (via search) | Checks Go imports against declared layers. A fork exports the dependency graph as Mermaid. | The import graph of a Go module can be drawn from code. `go list -deps -json` gives the same data with the standard library only. |
| [LikeC4, dynamic views](https://likec4.dev/dsl/views/dynamic/) (read; open source) | Dynamic views describe one use case with steps, as a diagram or a sequence. They live only in the view, not in the model. | Flows are drawn by hand even in a model tool. |
| [Write the Docs, docs as code](https://www.writethedocs.org/guide/docs-as-code/) (read) | Docs follow the same review and merge workflow as code; teams can "block merging of new features if they don't include documentation". | The system doc changes in the pull request that changes the behaviour. |
| [Marius-Juston/agent-harness, issue 299](https://github.com/Marius-Juston/agent-harness/issues/299), 2026-10-19 (read; practitioner, one repo) | An LLM sweep for stale claims in specs and design docs: precision 0.87 (34/39) on docs that had drifted, 0.22 (2/9) on maintained files, recall 7 of 9. 31 of 34 earlier drift fixes deleted the doc. The sweep stays advisory. | An LLM drift check gives false alarms on maintained docs. Less hand-written text means less drift. |
| [asimov-agent/llgenie, pull request 646](https://github.com/asimov-agent/llgenie/pull/646) (via search; practitioner) | A bot regenerates `docs/ARCHITECTURE.md` from the git diff since the last refresh, in its own pull request. | A refresh as a separate pull request is a known pattern. |

What can be generated from code and what cannot:

- **Generated:** the packages and their imports, the commands of the CLI, the list of scenarios in force (`factory scenarios`), the rules table. A check can compare them with the code on every pull request.
- **By hand:** flows across sessions, people and GitHub, such as the increment workflow. No code holds them. Only a review checks them.

## In which pull request the system doc changes

| Source | Finding | Can take |
|---|---|---|
| [OpenSpec, concepts](https://raw.githubusercontent.com/Fission-AI/OpenSpec/main/docs/concepts.md) (read) | The change holds a delta (ADDED, MODIFIED, REMOVED). The current specs change only at archive, after the code. | The design note holds the future state; the system doc changes with the code. |
| [Ubl, Design docs at Google](https://www.industrialempathy.com/posts/design-docs-at-google/) (read; practitioner) | "If the designed system hasn't shipped yet, then definitely update the doc." After that, design docs drift and become a set of amendments. A design doc has a system-context diagram. A small change gets a "1-3 page 'mini design doc'". | A design note is a record of a change, not a current-state doc. |
| Rule 36 and `AGENTS.md` (this repo) | Docs describe the current state. ADRs are in force on `main` once merged. | A system doc changed in the design pull request would describe a state that does not exist until the code merges. |

Three placements follow:

1. **Code pull request.** The implementing agent updates the system doc with the code; the reviewer checks it against the diff (rule 12). The doc on `main` always matches the code on `main`. The engineer sees the change only if the system doc is an owned path, which then makes the code pull request wait for the engineer (rule 27).
2. **Design pull request, as a delta.** The design note shows the changed flow as a diagram of the future state. The code pull request moves it into the system doc. The engineer reads the flow before approval, and the doc on `main` stays current.
3. **Separate pull request after the code.** A session refreshes the system doc from `main`. The doc lags behind the code between the merges.

## Flows, diagram size and agents that draw

### Flows that the factory and its projects need

The factory has no runtime of its own beyond the CLI; its flows run across people, sessions and GitHub. From the docs on `main`:

| Flow | Question it answers | Where it lives today |
|---|---|---|
| Increment workflow | Which step comes next, and who acts? | `AGENTS.md` (text), sketch in the design note of 02 |
| Merge gate | Where does a pull request wait for the engineer? | sketch in the design note of 02; rules 16, 27, 31 |
| Write permissions | Which session may change which path? | rules 28, 29; path groups in the design note of 02 |
| Credentials and network | Which identity and which hosts does each session use? | rule 33 |
| Feedback loop | How does a correction in a project reach the factory? | rule 35 (increment 06) |
| Install and update | How does a project get a new factory version? | ADR 0002 (increment 05) |

A project adds the runtime flows of arc42: main use cases, external interfaces, start and stop, errors.

### How much a reader takes in

| Source | Finding | Can take |
|---|---|---|
| [Cowan, The magical number 4 in short-term memory](https://memory.psych.missouri.edu/assets/doc/articles/2001/cowan-bbs-2001.pdf), Behavioral and Brain Sciences 24(1), 2001 (via search; measured) | A central limit of about 4 chunks, 3 to 5 in adults; Miller's 7 was "a rough estimate and a rhetorical device". | About 7 elements fits when a reader can group them; one diagram answers one question. |
| [Störrle, Good, bad, ugly: layout matters](https://www.pst.ifi.lmu.de/~stoerrle/talks/Good_Bad_Ugly-LayoutMatters.pdf) (via search; measured, small samples) | Size measured as number of elements; about 50 is given as a maximum safe size. Layout flaws such as crossings add as much load as elements. | An upper bound, not a target. Layout counts as much as size. |
| [Sözen, Say, Kılıç, working memory and complex diagrams](https://ada.atilim.edu.tr/entities/publication/e4644cf7-f5ad-47c1-a23a-967a8cb0b5cc/full) (via search; measured, students) | Higher working memory capacity helps with complex activity diagrams. | Smaller diagrams help readers with less capacity. |

No study found counts how many diagrams a reader takes in per document. arc42 and C4 advise few views, chosen by relevance.

### Mermaid

| Source | Finding | Can take |
|---|---|---|
| [GitHub docs, creating diagrams](https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/creating-diagrams) (read) | Mermaid renders in Markdown files, issues and pull requests. The supported syntax depends on the Mermaid version that GitHub uses. | Mermaid in the system doc renders without a build step. |
| [Mermaid config schema](https://mermaid.js.org/config/schema-docs/config.html) (via search) | Defaults: `maxTextSize` 50000, `maxEdges` 500. | Size limits are far above what a reader takes in; they do not constrain. |
| [MermaidSeqBench, arXiv 2511.14967](https://arxiv.org/abs/2511.14967), v2 2026-04-25 (read; measured) | 132 samples; LLM judges found "significant capability gaps across models" in sequence diagrams from text. | A sequence diagram drawn by an agent needs a check or a review. |

Layout cannot be pinned in GitHub's Mermaid; there is no model behind the diagrams, so the same element can have a different name in two diagrams.

### Agents that draw diagrams (rule 18)

| Source | Finding | Can take |
|---|---|---|
| [arXiv 2608.28403](https://arxiv.org/abs/2608.28403), Recovering software architecture intent from historical work items (read; measured, one company) | An LLM pipeline recovers C4 diagrams from work items. "High stability for architectural entities but lower stability for their relationships", "with relative variance compounding across generation steps". | Rule 18 holds. The finding is lower stability of relations; the need for human review is this repo's conclusion, not the paper's statement. |
| [Structurizr, DSL generation with AI](https://docs.structurizr.com/ai/dsl-generation) and the [C4 author's post](https://www.patreon.com/posts/152027514) (via search; vendor) | For diagrams from AI, a "models as code" tool gives better results, because the model enforces consistency. | Building-block views from a model that the engineer owns, as in rule 18 and the backlog. |
| Simon Brown, The C4 model, Structurizr vNext, and AI, [BCS Edinburgh 2026-05-20](https://edinburgh.bcs.org/events/2026/260520.htm) (via search; no slides found) | Talk abstract: hand-drawn diagrams turned into DSL, many diagrams from one model. | Same. |

Rule 18 covers building-block views. Flows across people and sessions have no model to generate them from, so they are drawn by hand and checked by review.

### Rejected

- **A full arc42 document:** twelve sections, most of which the factory does not need, such as deployment and quality scenarios; C4 and arc42 both advise only views that add value.
- **Context Mapper or Structurizr now:** the factory has one Go module and one bounded context; the backlog builds a model when a project has two packages or three contexts.

## What the engineer reads per increment today

| Item | Size for increment 02 | Written for |
|---|---|---|
| Design note | about 4,900 words, 1 diagram | agents, tests, the engineer's approval |
| Pull request descriptions | 1 to 2 per increment | the engineer and the history |
| Review result | `evidence/<nn>-<slug>/review.md` | the engineer and the merge gate |
| Recordings of demonstrated scenarios | per scenario | the engineer |
| Retro | `retro.md` | the engineer and the factory |
| `README.md` | about 730 words, the status | the engineer |

No item shows the whole system as it is after the increment. The engineer rebuilds that view from the design notes of all increments.
