---
name: designing-increments
description: Writes the design note of an increment, with goal, existing tools, interface, scenarios for the person and rollout, and updates the glossary and context map when terms or bounded contexts change. The agent uses it when the person asks to start, plan, scope or specify the next piece of work, or to rework a design note after comments, before any test or code exists. Not for the acceptance tests, which writing-acceptance-tests covers.
---

# Designing increments

A design note is the spec of one increment. The person approves it, a separate session writes the acceptance tests from it alone, and the implementing agent builds against those tests. The note is done when:

- the person can approve it without the conversation that led to it,
- the test writer finds in it every case and every output that a test must check,
- and `factory check-design <increment-dir>` passes.

## Who decides

The person owns the terms, the boundaries of bounded contexts, the dependencies between contexts, and each choice whose consequences reach beyond the code. The agent finds facts in the repo and its docs itself, and decides implementation details itself.

The agent asks the person directly in the conversation, never through the note. Each question says what it is about, so that it reads without context, and names the consequences of each option. A question carries no recommended answer, so that the person decides instead of agreeing. Questions come in rounds; each round builds on the answers to the last one.

## What the note holds

`assets/design-note.md` shows the structure. `factory check-design` is the authority on the format and names each mistake. Beyond the format:

- **Existing tools:** for each part that the increment builds, the tools, platform features and authoritative guidance checked, the decision and the reason. Building a part needs a reason. Only a design session reads the open web; other sessions reach only the allowed hosts.
- **Interface:** everything that a test must match: commands, inputs, output lines, error codes and exit codes.
- **Scenarios:** behaviour in plain sentences for the person, with a concrete example where it helps. They name no exit or error codes and point to no other section; the test writer takes those from the Interface. Each names how it is verified: automated (a test file from the test writer), demonstrated (a separate session runs and records it) or manual (the person runs it).
- **Rollout:** how the change is switched off once in use, and the signal that shows that it fails.
- **Replaced scenarios:** when the increment changes behaviour that a scenario of an earlier increment fixes, the new note replaces that scenario and gives the reason. The earlier note stays as it is. `factory scenarios` lists the scenarios in force.

Before the pull request, a fresh context that reads only the scenarios explains each one in its own words. The agent rewrites each scenario that it cannot explain alone.

## Domain model

When the increment adds or changes terms, bounded contexts or dependencies between contexts, the glossary and the context map change with it; `assets/glossary.md` and `assets/context-map.md` show their format. A new context or a new dependency between contexts also gets an architecture decision record. A term with two meanings, a new context or a new dependency is a question for the person.

## Boundaries

The design pull request changes only the design note and the domain model. After opening it, the agent stops: the person approves it as code owner.
