---
name: writing-commit-messages
description: Writes the message of a commit on a work branch. Use before each commit.
---

# Writing commit messages

A commit on a work branch tells the reviewer what changed and why. On `main`, the squash merge replaces these messages with the pull request's title and description, so a commit message serves the review.

## Before the commit

- Run the project's tests, for example `go test ./...`.
- One topic per commit. A commit with tests contains no code, and a commit with code contains no test file.

## Subject line

- A complete sentence in the imperative, at most 50 characters (72 at most), capitalized, without a final period.
- It names what the change does, not which files it touches.
- It contains no internal numbers such as increment or scenario ids.

## Body

- A blank line after the subject.
- Lines of at most 72 characters.
- It explains what changes and why, not how. The diff shows how.
- It names a decision that the design note did not cover, and the reason for it.

## After the commit

Push the work branch, fast-forward only, never with force.

Example:

```
Count review results only for the reviewed commit

The stop hook let a session end when any review result existed. A
result for an older commit no longer counts, so each code commit
gets its own review.
```
