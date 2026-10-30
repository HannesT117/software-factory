---
name: designing-increments
description: Writes the design note of an increment, and the domain model when terms or bounded contexts change. Use when an increment starts, before any test or code exists.
---

# Designing increments

The design note is the spec of one increment. The engineer approves it with its scenarios. A separate session writes the acceptance tests from it, and the implementing agent builds against those tests. So each sentence must hold without the rest of the conversation.

## Steps

1. **Find facts.** Read `README.md`, `docs/principles.md`, the decision records, the domain model in `docs/domain/` if it exists, and the notes of earlier increments. Find facts in the code and the docs. Do not ask the engineer for facts that the repo holds.
2. **Search for existing tools.** For each part that the increment builds, search for tools, libraries and platform features that already do it. Record each in `## Existing tools`: the tools checked, the decision (used, not used, later) and the reason. A part with no tool found names what was searched. Building a part needs a reason.
3. **Ask the engineer.** The engineer owns the terms, the context boundaries and the dependencies between contexts. Ask about them in rounds of numbered questions. Each question says what it is about and the consequences of each option. Give no recommended answer: the engineer decides instead of agreeing. Ask the next round only after the answers to the current one.
4. **Write the note** from `assets/design-note.md`. Keep the headings of the template.
5. **Write the scenarios for the engineer.** Each scenario states behaviour in plain sentences, with a concrete example where it helps. It does not point to other sections of the note and names no exit codes or error codes; the test writer takes those from `## Interface`. Each scenario names its verification kind:
   - **automated:** a test file from the test writer
   - **demonstrated:** a separate session runs it and records it, and the engineer views the recording
   - **manual:** the engineer runs it
6. **Have the scenarios explained.** Start a fresh context that reads only the `## Scenarios` section and says in its own words what each scenario means. Rewrite each scenario that it cannot explain without the rest of the note.
7. **Update the domain model** when the increment adds or changes terms, bounded contexts or dependencies between contexts. Use `assets/glossary.md` and `assets/context-map.md`. A new context or a new dependency between contexts also gets an architecture decision record.
8. **Check.** `factory check-design docs/increments/<nn>-<slug>` must pass.
9. **Open the design pull request.** It changes only the design note and the domain model. Stop. The engineer approves and merges it.

## Format of the design note

- Path: `docs/increments/<nn>-<slug>/design.md`. Line 1 is `# Increment <nn>: <title>`, where `<nn>` equals the folder's number.
- Required headings, exact text, any order: `## Goal`, `## Out of scope`, `## Existing tools`, `## Interface`, `## Scenarios`, `## Evidence`. Other headings are allowed. The optional `## Domain` lists new or changed terms, the bounded contexts that the increment touches, and new dependencies between contexts.
- A scenario line matches `^- \*\*scenario (\d{2})\.([1-9]\d*)\*\* \((automated|demonstrated|manual)\) \*\*(.+?)[.?]\*\*`. The bold title ends with a period or a question mark. Any other line that starts with `- **scenario` is an error. Lines in fenced code blocks are ignored.

## Format of the domain model

A project without bounded contexts yet has no domain model.

- `docs/domain/glossary.md`: line 1 is `# Glossary`. One line per term: `- **<Term>**: <definition>.`, optionally followed by `_Avoid_: <word> (<reason>), <word> (<reason>).`
- `docs/domain/context-map.md`: line 1 is `# Context map`. One `##` heading per bounded context, with the lines `- Folder: <path>`, `- Responsible for: <text>`, `- Not responsible for: <text>`, and one line `- Depends on: <context> (<pattern>)` per dependency. The patterns are those of DDD context mapping, for example open host service, anticorruption layer or conformist.

## Stop and ask

Stop and ask the engineer when an architecture check fails, when a term has two meanings, or when the increment needs a new context or a new dependency between contexts.
