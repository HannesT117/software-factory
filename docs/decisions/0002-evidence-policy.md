# 0002: What counts as evidence that a change works

- Status: proposed
- Date: 2026-10-26

## Context

You want to trust the acceptance tests without reading every line.

- Agent-written unit tests often run but assert little (VibeCheck), and they add little to issue resolution (arXiv 2602.07900).
- Agents edit or game tests, more often on ambiguous tasks (EvilGenie).
- An agent-run TDD loop gave no quality gain at 3–8.5× the tokens (Böckeler, Sonnet 4.6 / Opus 4.8).
- What helps is independence: tests written separately from faulty code catch more faults (arXiv 2607.05139). StrongDM keeps scenarios hidden from its agents.

## Decision

Evidence is:

1. Black-box acceptance scenarios that you approve. They are locked after approval.
2. A holdout part of the scenarios, kept in a private repo that agent sessions do not get. CI reports only scenario ids.
3. Properties that you state. The agent writes only the generators. The property statements are locked.
4. Types.
5. The mutation score of these suites, as a measure of their strength, not as a gate.

Not evidence:

- Agent-written unit tests. Agents may write them as scaffolding. No check requires them.
- An agent's review of its own work.

Process rules:

- Agents get no test-first or red-green instruction.
- An inferential review (LLM review in a fresh context) runs once per slice, as a final gate.

## Consequences

Good:

- Trust rests on artefacts that you approved and the agent cannot change.
- Fewer tokens than per-step reviews or TDD loops.

Bad:

- You must write or approve scenarios and properties. That is the main human cost.
- The holdout repo needs its own access setup.

## Re-check when

A study with current models shows that agent-written tests or agent-run TDD add value, or the model changes.

## Sources

Böckeler, TDD in the agent loop; Jain, Are harnesses worth the cost?; VibeCheck; arXiv 2602.07900, 2607.05139, 2511.21654; StrongDM via Willison. Details in [research](../research/2026-10-factory-basis.md).
