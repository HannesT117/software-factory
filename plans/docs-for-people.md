# Plan: documentation of the whole system for people

- **State:** research. The result is a set of options for the engineer, not a decision.
- **Target pull requests:**
  1. Research notes in `docs/research/2026-11-docs-for-people.md`, from the branch `claude/research-docs-for-people`.
  2. After the engineer's decision: a design pull request for an increment. Its number follows the order of work, and the planned increments after it shift.
- **Runs as:** a design session (`FACTORY_SESSION=design`, may read the open web).

## Context

- Each increment has a design note, `docs/increments/<nn>-<slug>/design.md`. It is written for the test writer, the implementing agent and `factory check-design`: Interface, format rules, scenarios, and one diagram that answers one question.
- `README.md` is one page with the status. `docs/domain/` holds the glossary and the context map. `docs/principles.md` holds the rules.
- Rule 18: agents draw only small sketches of about 7 elements. Building-block views come from a model that the engineer owns. The backlog has "Architecture model as code, checked against the import graph".
- Rule 36 and `AGENTS.md`: docs describe the current state, and a doc points to a config file instead of restating it.
- Increment 04 moves the comprehension questions to the design review. After the code, the reviewer explains deviations and draws a sketch only where the built flow deviates.

## The engineer's statement, 2026-11-03

One diagram per design note may be too little. Separate the implementation plan for a model from the documentation for people? The documentation for people should:

- keep the whole system in view, instead of thinking in increments
- be written for people, not for agents
- show flows, with more than one small sketch
- be easy to read

## Questions for the research

1. How do others separate specs for agents from documentation for people? Candidates: GitHub Spec Kit (spec, plan, tasks), Kiro (requirements, design, tasks), `AGENTS.md` next to `README.md`, Diátaxis (explanation against reference), arc42 (runtime view), C4 (dynamic diagrams), living documentation and docs as code.
2. How does a system doc stay current while agents change the system? Who updates it, in which pull request, and which check shows that it matches the code: diagrams generated from code, or drawn by hand?
3. Docs describe the current state. Does the system doc change in the design pull request, which describes a future state, or in the code pull request?
4. Which flows a factory and its projects need, how many sketches a reader takes in, the limits of Mermaid, and what rule 18 and its sources say about agents that draw diagrams.
5. What the engineer reads per increment today, and after each option.

## Steps

1. Read `README.md`, `AGENTS.md`, `docs/principles.md` (rules 18, 19 and 36), the design note of 02, `core/skills/designing-increments/` and `docs/domain/`.
2. Research the questions on the open web, with a link for each source. Keep rejected options with their reason, not the course of the research.
3. Write 2 to 4 options. Each says what changes, what the engineer reads, what it costs, and its consequences, with no recommendation.
4. Commit the research notes on the work branch and push. Send the options to the orchestrating session, and wait for the engineer's decision. Then open the research pull request.

## Decisions of the engineer

### 2026-11-03, round 1 (options A to D)

- Option C: hand-written flows, updated in the code pull request and checked by the reviewer, plus generated structure (package graph, CLI commands, scenarios in force) with a check that fails when it differs from the code.
- Not one file `docs/system.md`: a sensible split of the docs for people. The next round proposes 2 to 4 splits, each with the files, what goes where, what the engineer reads and the consequences, from best practice (Diátaxis, arc42, C4, docs as code, how projects split explanation, flows, reference and how-to).
- The paths of these docs are owned by the engineer as code owner, so each code pull request that changes them waits for the engineer.
- A skill for writing these docs may be needed. Each split says whether it needs such a skill (see the backlog row for `reviewing-architecture` and `writing-docs`) and what it would hold.
- Rule 18 needs a decision for option C (building-block views come from a model that the engineer owns; code is not such a model). It goes as a question with consequences into the same round.

### 2026-11-03, round 2 (split and rule 18)

- Split S2, arc42 lean: `docs/architecture/context.md`, `docs/architecture/runtime.md`, links to the ADRs and the glossary instead of copies, and a smaller skill `writing-docs`.
- Rule 18: no building-block view for the factory's own docs ("fairly pointless, because it is a collection of scripts"). Rule 18 stays as it is.
- Open for increment 05 or later, not decided: whether projects that use the factory get a generated building-block view.
- Settled by the agent (rule 38, benefit against cost): the factory's docs get no generated files. The CLI commands are the external interface of arc42 section 3, so `context.md` names the CLI and points to `factory` without arguments, which prints the commands; the scenarios in force are printed by `factory scenarios`. A generated copy would restate what a command prints (rule 36, AGENTS.md: a doc points instead of restating), and as an owned path it would make code pull requests wait for the engineer. So for the factory, the generator and its check of option C are dropped; `building-blocks.md` does not exist. Whether projects get generated parts stays open with the building-block view.

## Checkpoints
- 2026-11-03T00:52Z step 1 done; claude/research-docs-for-people@7759cd2 (no commits yet, from main); next: step 2, research question 1; open: none. Note: `docs/domain/` does not exist in the factory repo; only the skill assets define its format.
- 2026-11-03T00:55Z step 2 done; claude/research-docs-for-people@1401381; next: step 3, write the options; open: none
- 2026-11-03T00:56Z step 3 done; claude/research-docs-for-people@241f36b; next: step 4, wait for the engineer's decision on options A to D (sent to the orchestrating session), then open the research pull request; open: which option
- 2026-11-03T01:05Z round 2 written; claude/research-docs-for-people@cbe5495; next: wait for the engineer's choice of split (S1 to S3) and the rule 18 answer (R1 to R3), sent to the orchestrating session; then record it under Decisions and open the research pull request; open: split, rule 18
