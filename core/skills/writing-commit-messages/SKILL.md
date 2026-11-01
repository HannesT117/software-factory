---
name: writing-commit-messages
description: Writes commit messages on a work branch, with an imperative subject and a body that explains what changes and why, one topic per commit. The agent uses it before every commit, when it splits work into commits, or when it fixes a message before the push. Not for the title and description of a pull request, which writing-pull-requests covers.
---

# Writing commit messages

The squash merge replaces commit messages on `main` with the pull request's, so a commit message on a work branch serves only the review.

- **One topic per commit.** Test files and code never share a commit. The tests come first, so that they can be shown to fail before the code exists.
- **Subject:** a complete sentence in the imperative, at most 50 characters, capitalized, without a final period. It names what the change does, not which files it touches, and contains no increment or scenario ids.
- **Body:** after a blank line, lines of at most 72 characters. It explains what changes and why; the diff shows how. It names each decision that the design note did not cover, with its reason.
- The project's tests pass before the commit. The work branch is pushed after it, fast-forward only.

Example:

```
Count review results only for the reviewed commit

The stop hook let a session end when any review result existed. A
result for an older commit no longer counts, so each code commit
gets its own review.
```
