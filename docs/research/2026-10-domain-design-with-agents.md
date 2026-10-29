# Domain-driven design with coding agents, October 2026

These are the sources behind the domain model, the Domain section of design notes, the interview in `designing-increments` and the architecture checks for projects. Each entry says what the source found and what this repo takes from it.

Labels:

- **read**: the primary source was read for this note.
- **via research agent**: a research agent read it. It was not re-read.
- **tried**: run on a small test tree, 2026-10-29.
- **measured**: a study with numbers. **practitioner**: experience without measurement. **consensus**: a consultancy's shared view.

## What is measured

| Source | Finding | Taken |
|---|---|---|
| Sharma, ContextCov, arXiv 2603.00822, 2026-02-28 (read; measured) | Rules from agent instruction files, turned into executable checks, reached "88.3% constraint compliance (vs. 67.0% and 50.3%)" for prompt-only and LLM reflection. | Architecture rules as checks, not as text. |
| Dente, Satriani, Papotti, arXiv 2605.06445, 2026-05-07 (via research agent; measured) | Each added structural constraint (architecture pattern, database, ORM) lowers the agents' pass rate, by 27.28 points on average when fully specified. | One rule first: domain isolation. No full layering everywhere. |
| Eisenreich, Jusic, Wagner, arXiv 2603.26244, 2026-03-27 (via research agent; measured, one company, 3 experts) | LLM drafts of a glossary, of EventStorming results and of bounded contexts were usable. Aggregates and the mapping to an architecture were impractical. The LLM split contexts more finely than the experts found practical. | The agent drafts; the engineer decides boundaries. One context is a valid start. |
| Gloaguen et al., arXiv 2602.11988, 2026-02-12 (via research agent; measured) | Context files such as AGENTS.md did not generally raise task success and cost more than 20% extra. They help for practices an agent would not guess. | A pointer to the glossary, not a copy of the model in AGENTS.md. |
| Canedo, arXiv 2608.21747, 2026-08-22 (via research agent; measured, 6 models) | On the strongest models the format of an architecture spec barely matters. On weaker models, code-near formats close most of the gap. | No special spec format. |
| Li et al., arXiv 2508.06414; Le et al., arXiv 2510.03178, 2025 (via research agent; measured) | Removing descriptive names costs up to 30 points in code generation. | Names matter to models. No study measures a glossary itself. |

No study measures a glossary, a term check or context boundaries in agent work. The language side rests on practitioner experience. A glossary may mainly help the humans; Matt Pocock's skill docs name this objection.

## Ubiquitous language

| Source | Finding | Taken |
|---|---|---|
| Michael Plöd, `como-to-ai-demo`, KanDDDinsky 2026, 2026-10-15 (read; practitioner) | One `glossary.yaml` per context with words to avoid, each with a reason and a decision link. `GlossaryLint` splits identifiers and fails on an avoided word; the README warns it is token-based, so false positives occur. `PromptLint` asks back when a request is ambiguous. AGENTS.md and agent permissions are generated from the model. CODEOWNERS guards the model files. No license file. | Avoided words with a reason. Model files that agents cannot change alone. Ideas only, no code. |
| Matt Pocock, skills `domain-modeling` and `grilling`, MIT, 2026-10-27 (via research agent; practitioner) | Glossary in Markdown: `**Term**: definition. _Avoid_: x, y`. The interview numbers open questions in rounds. The agent finds facts, the user decides. Each question comes with a recommended answer. | The Markdown form and the interview in rounds. Not the recommended answer: it invites agreement instead of a decision. |
| Nick Tune, automatic code reviewer, 2025-11-30 (via research agent; practitioner) | "Coding agents will ignore instructions in the system prompt on a regular basis." Lint rules carry conventions; an LLM review checks names and domain logic that leaks out of the core. | Review checks names against the glossary. |
| cspell 10.4.0, `flagWords` (tried) | Reported `Forbidden word (Booking)` inside `BookingRefundService`. Its default dictionaries are too broad to report unknown domain words. | Later: a check for avoided words. A check for unknown terms needs custom code. |
| ddd-crew, `ai-ddd-prompts-and-rules`, CC BY-SA 4.0, 2026-08 (via research agent) | One rule on terminology mismatch, empty prompt sections. | Nothing. |

## Domain core and dependencies

| Source | Finding | Taken |
|---|---|---|
| Nick Tune, Enforced application architecture for agents and humans, 2026-08-13 (read; practitioner) | "If you're only going to add one convention in your codebase, enforcing domain isolation is probably the best one." "AI does dumb stuff like that all the time even with a million lines of markdown screaming at it not to do that." | The one rule: the domain imports nothing outside the domain. |
| Vaughn Vernon, This is your agent on DDD, WSO2CON, 2026-05-21 (via research agent; practitioner, talk abstract) | Even with significant guidance, the agent falls back to CRUD when it models a domain. | No tactical DDD scaffolding. |
| Thoughtworks Technology Radar Vol. 34, 2026-04 (via research agent; consensus) | Teams write hexagonal architecture into spec constitutions. Deterministic tools plus LLM review fight drift (Assess). No entry on DDD. | Deterministic rules first (rule 17). |
| dependency-cruiser 18.7.0, MIT (tried) | Two rules, `domain-isolated` and `context-via-api`, printed `error domain-isolated: src/sales/domain/ticket.ts → src/sales/adapters/db.ts` and exited non-zero. The same rule blocks npm packages and Node built-ins in the domain unless allowlisted. | Candidate for TypeScript projects, beside an oxlint rule. Decided in increment 03. |

Ports and adapters follow from domain isolation when an adapter is needed. They are not a rule of their own.

## Bounded contexts and context maps

| Source | Finding | Taken |
|---|---|---|
| Plöd (read) | The context map is a closed world: contexts without a relationship on the map may not depend on each other. A side file holds "not responsible for", because CML cannot express it. | "Not responsible for" in the context map. The closed world later. |
| Context Mapper (CML), Apache-2.0, v6.12.0, 2024-08 (via research agent) | The only DSL for context maps. Validates and generates diagrams. Needs JDK 17. Checks no TypeScript or Go code. | Later, from three contexts or a team. |
| ddd-crew Bounded Context Canvas v5, CC BY-SA 4.0, YAML with JSON Schema since 2026-10-15 (via research agent) | The only ddd-crew canvas in a checkable file format: purpose, classification, messages, language, decisions, open questions. | Later, one per context from the second context. |
| ddd-crew Context Mapping, Domain Message Flow, Aggregate Design Canvas, Core Domain Charts (via research agent) | Workshop formats: images, Miro, draw.io. | Not taken: no checkable file. |

## What the engineer owns

| The engineer decides | The agent does |
|---|---|
| each term and its meaning, and which words to avoid | finds terms in code, docs and requests; flags conflicts; drafts entries |
| context boundaries: split, merge, "not responsible for" | proposes splits with the coupling it finds in the code |
| each dependency between contexts and its pattern | drafts the map line and the dependency rule |
| what belongs in the domain core, and which libraries it may use | keeps the code inside the rules; reports a failing check |

Sources: Eisenreich et al. (LLM as a sparring partner, not for full automation); Plöd (the domain experts of a context decide); Miłosz Smółka, Three Dots Labs, 2026-08-20 ("you need someone accountable in the end"); Annegret Junker, codecentric, 2026-03-04 (humans ran the discovery). All via research agent.

## Enforcing the discussion

From strongest to weakest:

1. Model files and rule files are protected; a change needs the engineer's approval (Plöd: CODEOWNERS). In this repo: path groups and the approval gate. They run on the server and work with any agent.
2. A failing check that the agent cannot fix in the same pull request, because the rule file is in a protected group. The agent must stop and ask.
3. A closed world: a new context or dependency needs a model change first.
4. Error messages that name the reason and the decision (Plöd).
5. Agent tool permissions and hooks: early feedback, one tool only.
6. Inferential gates, such as clarification markers in Spec Kit: the model judges itself.

## Options not taken

- Tactical DDD scaffolding and aggregate canvases as rules: LLMs are weakest at aggregates, and agents fall back to CRUD.
- EventStorming simulated by several persona agents (community plugins): no evaluation, and it replaces the human discussion that the method exists for.
- Full hexagonal layering for every module: each structural constraint lowers the agents' pass rate.
- The whole model copied into AGENTS.md: cost without gain.
