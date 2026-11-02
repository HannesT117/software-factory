---
tags: [designing-increments]
max_turns: 30
timeout_seconds: 900
allowed_tools: [Read, Glob, Grep, Skill, Write, Edit, Bash]
---

Write the design note for the next increment of greeter. I have answered every question:
- A flag `--lang` picks the language. Allowed values: `en` and `de`; without the flag, `en`.
- German greets with `Hallo, <name>!`.
- An unknown value prints `greet: error --lang: unknown-language <value>` and exits with 2.
- The new term is Language: the language of a greeting. It joins the glossary.

Input files: README.md, docs/domain/glossary.md, docs/increments/01-greeting/design.md, evidence/01-greeting/review.md, go.mod
