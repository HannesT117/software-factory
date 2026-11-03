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

## Checkpoints
- 2026-11-03T00:52Z step 1 done; claude/research-docs-for-people@7759cd2 (no commits yet, from main); next: step 2, research question 1; open: none. Note: `docs/domain/` does not exist in the factory repo; only the skill assets define its format.
- 2026-11-03T00:55Z step 2 done; claude/research-docs-for-people@1401381; next: step 3, write the options; open: none
