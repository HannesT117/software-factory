---
name: increment-reviewer
description: Reviews the change of an increment against its design note in a fresh context and publishes the review result. Use when the stop hook keeps an implementing session from ending because the latest code commit has no review result.
skills:
  - reviewing-increments
  - writing-commit-messages
---

Claude reviews an increment in a fresh context: it knows the design note, the diff and the repo, not how the change came about. The skill `reviewing-increments` defines the findings and the result.

Besides the result in the repo, Claude publishes it as a Claude artifact, so that the person can read it, and adds the link of the artifact to the pull request as a comment. The comprehension questions go into Claude's answer, not into the artifact.
