# Increment 04: handover and feedback loop

This is a draft. It is rewritten with the `designing-increments` skill after increment 03. Only then does it go to the engineer.

## Goal

- Every increment in a project ends with a retro. Points that concern the factory reach the factory repo as issues that the engineer checks and submits. Nothing gets lost, and the project never patches factory files.
- A factory change reaches a project as a new version, in a pull request with a visible diff.

## Out of scope

- automatic fixes in the factory. An issue is input for a normal factory increment.
- feedback from more than one project

## Interface

To be written with `designing-increments`. Planned commands:
- `factory check-retro <increment-dir>`
- `factory feedback <increment-dir>`

## Scenarios

- **scenario 04.1** (automated) **A complete retro passes.** Given a `retro.md` with the required sections (metrics, agent version, corrections, factory points) and an issue link for every factory point, `factory check-retro` exits with 0.
- **scenario 04.2** (automated) **A faulty retro fails.** Given a retro with a missing section, or a factory point without an issue link, `factory check-retro` exits with 1 and names what is missing. It checks the structure only, not the counts.
- **scenario 04.3** (automated) **Feedback becomes a prefilled link.** Given a retro with a factory point, `factory feedback` prints one link to a new issue in the factory repo, with project, increment, observation, evidence and a proposed check filled in. It sends nothing.
- **scenario 04.4** (manual) **The engineer submits.** The engineer opens the link, removes private data, checks for an earlier issue on the same point, and submits. A point seen twice gets the label `seen-twice`. The engineer adds the label `accepted` to issues that factory agents may read.
- **scenario 04.5** (manual) **An update reaches the project.** A new factory version reaches a project through `factory install` in a pull request, and the diff shows every changed factory file.
- **scenario 04.6** (manual) **End to end.** In a fixture project, one full increment runs: the three pull requests, review, retro, a submitted feedback issue, squash merges.

## Evidence

To be written with `designing-increments`.

## Open questions

To be written with `designing-increments`.
