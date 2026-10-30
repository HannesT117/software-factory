---
name: writing-pull-requests
description: Writes the title and description of a pull request so that they work as the squash-merge commit message on main. Use when opening a pull request or when its diff changes.
---

# Writing pull requests

With squash merges, the title and description of a pull request become the commit message on `main`. The description gives a reader without context an overview of the change.

## Before opening one

- Look at the open pull requests. A change that fits the topic of an open pull request goes there.
- One problem per pull request. If the title needs "and" to join two topics, or the description gets long, split the pull request.
- Weigh the benefit of each change against its cost. Drop a change with little benefit, also one that a review suggested, and name what was dropped in the report.

## Title

- A complete sentence in the imperative, at most 50 characters (72 at most), capitalized, without a final period.
- It names what the change does, not which files it touches, and it is clear without the description.
- It contains no internal numbers such as increment or scenario ids.

Example: `Reject pull requests that mix path groups`, not `Update prcheck.go` or `Increment 02: check-pr`.

## Description

Use the sections of the repo's pull request template, in this order. Delete an optional section that does not apply.

- **Summary** (required): two or three sentences that a reader without context understands: what changes, and which problem it solves.
- **Context** (optional): the problem, and what a reader must know to follow the change.
- **What changes** (required): the change at a high level. Core changes first, then mechanical ones such as renames or generated files.
- **Diagram** (optional): a Mermaid sketch of about 7 elements that answers one question.
- **Merge risk** (required): one-way door or two-way door; what the change could break beyond the diff; the one fact that makes it safe, and how far that fact is proven; how to undo the change.
- **Trade-offs** (optional): costs and shortcomings that are not obvious.
- **Evidence** (required): each piece with its kind: approved scenario, demonstrated, manual, review, ADR confirmation, agent test (scaffolding), or not verified.

The description is self-contained and describes the current state, not how the change came about. Before the approval, check that it still matches the diff.

## Notes for the approval

What the engineer should read, decide or run goes into the first comment of the pull request, not into the description.
