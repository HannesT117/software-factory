# Increment 04: handover and feedback loop

Status: proposed

This is a draft. It is rewritten with the `designing-increments` skill after increment 03. Only then does it go to the engineer.

## Goal

- Every increment in a project ends with a retro. Points that concern the factory reach the factory repo as issues. Nothing gets lost, and the project never patches factory files.
- A factory change reaches a project as a new version, in a pull request with a visible diff.
- After this increment, `health` starts.

## Out of scope

- automatic fixes in the factory. An issue is input for a normal factory increment.
- feedback from other projects than `health`

## Interface

To be written with `designing-increments`. Planned commands:
- `factory check-retro <increment-dir>`
- `factory feedback <increment-dir>`

## Scenarios

- **scenario 04.1** (automated) **A complete retro passes.** Given a `retro.md` with the required sections and an issue link for every point marked `factory`, `factory check-retro` exits with 0.
- **scenario 04.2** (automated) **Missing points fail.** Given a retro that lists fewer guard-hook denials or review findings than the increment's records, or a factory point without an issue link, `factory check-retro` exits with 1 and names what is missing.
- **scenario 04.3** (automated) **Feedback becomes an issue.** Given a retro with a new factory point, `factory feedback` creates one issue in the factory repo with project, increment, observation, evidence and a proposed check, and writes the issue link back into `retro.md`.
- **scenario 04.4** (automated) **A repeated point is flagged.** Given a factory point that matches an existing issue label, `factory feedback` adds a comment to that issue instead of a new one, and labels it `seen-twice`.
- **scenario 04.5** (manual) **An update reaches the project.** A new factory version reaches `health` through `factory install` in a pull request, and the diff shows every changed factory file.
- **scenario 04.6** (manual) **End to end.** In a fixture project, one full increment runs: design note, tests, implementation, review, retro, feedback issue, squash merge.

## Evidence

To be written with `designing-increments`.

## Open questions

To be written with `designing-increments`.
