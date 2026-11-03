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
