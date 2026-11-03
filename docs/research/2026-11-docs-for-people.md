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
