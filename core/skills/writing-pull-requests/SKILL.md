---
name: writing-pull-requests
description: Writes the title and description of a pull request, which become the squash-merge commit on main, and the first comment with notes for the person's approval. The agent uses it when it opens a pull request, when the diff or purpose of one changes before the merge, or to decide whether a change belongs in an open pull request. Not for commit messages, which writing-commit-messages covers.
---

# Writing pull requests

With squash merges, the title and description of a pull request become the commit message on `main`.

## Scope

One problem per pull request. A change that fits the topic of an open pull request goes there; a new pull request says in its first comment why no open one fits. When the title needs "and" to join two topics, or the description grows long, the change is two pull requests. The agent drops a change with little benefit for its cost, also when a review suggested it, and names it in its report.

## Title

- A complete sentence in the imperative, at most 50 characters, capitalized, without a final period.
- It names what the change does, not which files it touches, and is clear without the description.
- It contains no internal numbers such as increment or scenario ids.

Example: `Reject pull requests that mix path groups`, not `Update prcheck.go` or `Increment 02: check-pr`.

## Description

The headings of the repo's pull request template, in its order. An optional section that does not apply is left out.

- **Summary** (required): two or three sentences that a reader without context understands: what changes, and which problem it solves.
- **Context** (optional): what a reader must know to follow the change.
- **What changes** (required): the change at a high level. Core changes first, then mechanical ones such as renames or generated files.
- **Diagram** (optional): a Mermaid sketch of about 7 elements that answers one question.
- **Merge risk** (required): one-way or two-way door; what the change could break beyond the diff; the one fact that makes it safe, and how far that fact is proven; how to undo the change.
- **Trade-offs** (optional): costs and shortcomings that are not obvious.
- **Evidence** (required): each piece with its kind: approved scenario, demonstrated, manual, review, ADR confirmation, agent test (scaffolding), or not verified.

The description is self-contained and describes the current state, not how the change came about. It matches the diff at the time of the approval.

## First comment

What the person should read, decide or run before the approval goes into the first comment, not into the description: the description outlives the approval, the notes do not.
