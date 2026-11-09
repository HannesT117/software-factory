---
tags: [designing-increments]
max_turns: 30
timeout_seconds: 900
allowed_tools: [Read, Glob, Grep, Skill, Write, Edit, Bash, Agent]
---

Earlier in this session we wrote the note together. We first called the flag `--loud` and then renamed it, and I said that I may want a greeting in German later.

We are in the design step of greeter's increment 02. Its design pull request for docs/increments/02-shout/design.md (branch design-shout) is open, and you showed me these comprehension questions:

1. What does `greet --shout` print when no name follows the flag?
2. Which output tells that the flag `--shout` is not recognised?
3. What does `greet --shout Ada` print?

My answers:

1. It prints `greet: error args: missing-name` and exits with 1.
2. I don't know.
3. `greet: HELLO, ADA!`

Explain the gaps, as the skill designing-increments says.

Input files: README.md, go.mod, docs/, evidence/ (a git repository: main holds increment 01, the branch design-shout adds the note of increment 02)
