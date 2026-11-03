# Increment 05: docs for people

## Goal

- The engineer sees the whole factory as it is on `main`, in two docs for people next to the design notes for agents: `docs/architecture/context.md` (what the factory is for, its limits, who and what it works with, and its strategy) and `docs/architecture/runtime.md` (its flows, one small sketch each).
- The docs follow arc42 lean: only sections 1 to 4 and section 6. Decisions and rules are linked, not copied. There is no building-block view: the factory is a collection of scripts, and rule 18 stays as it is.
- The implementing agent changes the two docs in the code pull request when the code changes a flow or an external interface, with a new skill `writing-docs`. The reviewer reports a doc that no longer matches the diff. The engineer owns both paths, so a code pull request that changes them waits for the engineer.
- `AGENTS.md`, `docs/principles.md` and a process decision record name the new docs, where they change and why; `README.md` points to `context.md` as the entry to the whole system.

## Out of scope

- A building-block view, and docs generated from the code, for the factory: the CLI prints its commands (`factory` without arguments) and the scenarios in force (`factory scenarios`), so a generated copy would only restate them. For projects that use the factory, this stays open for install and distribution (increment 06) or later.
- Flows that are not built yet: install and update (increment 06), the feedback loop (increment 07). The increment that builds a flow adds it to `runtime.md`.
- The entry for `/docs/architecture/` in `.github/CODEOWNERS`: `.github/` is its own path group, so a separate pull request adds it, after the CI pull request has created the file.
- A check that counts the nodes of a sketch: a Mermaid parser costs more than the review that already reads each sketch. It goes to the backlog if the reviewer misses an oversized sketch twice (rule 4).
- A glossary for the factory: `docs/domain/` does not exist in this repo, and this increment adds no terms to it.

## Domain

- Terms: none.
- Contexts: none.
- Dependencies: none.

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Structure of the docs | arc42 sections and its lean and essential tips; Diátaxis; C4; matklad, `ARCHITECTURE.md`; the research notes `docs/research/2026-11-docs-for-people.md` | arc42 lean: sections 1 to 4 in `context.md`, section 6 in `runtime.md` | The engineer chose it from three splits. arc42 names the sections, so the docs use no other elements (rule 36). |
| Sketches | Mermaid on GitHub; C4 dynamic diagrams; Structurizr and LikeC4 | Mermaid, at most 7 nodes, one question each | Mermaid renders on GitHub without a build step. A model tool needs Java or Node and draws no flows across people and sessions. |
| Generated parts | `go list`, go-arch-lint, the CLI's own output | not used | The factory's commands and scenarios in force are printed by the CLI already; the engineer decided that the factory needs no building-block view. |
| Keeping the docs current | OpenSpec (current specs change after the code), Write the Docs on docs as code, an LLM drift sweep (precision 0.22 on maintained files) | the code pull request changes the docs, and the reviewer checks them against the diff | The doc on `main` then matches the code on `main` (rule 36). A sweep outside the pull request raises false alarms on maintained docs. |
| Skill for the docs | community skills for Diátaxis; the backlog rows `writing-docs` and `reviewing-architecture` | a new skill `writing-docs`; no `reviewing-architecture` | No community skill was checked against this repo's rules, and skills change only with `authoring-skills`. Without a building-block view, there is no architecture model to review. |
| Record of the decision | MADR; rule 25 | process decision record 0003 | The decision had real alternatives: four placements and three splits. |

## Interface

**`docs/architecture/context.md`.** Line 1 is `# Context`. It has the arc42 sections `## Introduction and Goals`, `## Constraints`, `## Context and Scope` and `## Solution Strategy`, in this order, and no other headings. Each section is a short list. Context and Scope names every external interface of the factory: the engineer, the agent sessions, GitHub (pull requests, rulesets, code owners, Actions), the model of another vendor, and the CLI. For the CLI, it points to `factory` without arguments for the commands and to `factory scenarios` for the scenarios in force, and lists neither. Decisions and rules are links to `docs/decisions/` and `docs/principles.md`, not copies.

**`docs/architecture/runtime.md`.** Line 1 is `# Runtime View`. Each flow is a `##` heading with the name of the flow. Under it come: one line, right after the heading, with the question that the sketch answers; one Mermaid block with at most 7 nodes, whose names are the roles, sessions, commands and GitHub features as `AGENTS.md` and `docs/principles.md` name them; and at most 5 sentences or list items. The flows on `main` after this increment are: the increment workflow, the merge gate, write permissions per session, and credentials and network. The text describes the state on `main`; it does not say what is new or what changed.

**Skill `writing-docs`** in `core/skills/writing-docs/`. The agent uses it when a change adds, removes or reorders a step of a flow in `runtime.md`, or adds, removes or changes an external interface in `context.md`, and when it writes either doc. It holds: which arc42 section a text belongs to; the format of a flow; that the docs describe the state after the change; that a change which touches no flow and no external interface leaves both docs alone. `factory validate-skills --strict` passes for it, and its evals run with and without the skill.

**Skill `reviewing-increments`.** A new kind of finding, `Docs for people out of date`: the diff adds, removes or reorders a step of a flow in `docs/architecture/runtime.md`, or changes an external interface in `docs/architecture/context.md`, and the doc does not match. The finding names the file and the line of the code, and the flow or interface. The result form `evidence/<nn>-<slug>/review.md` stays the same.

**Rules and records.** `AGENTS.md`, workflow step 3: when the code changes a flow or an external interface, the implementing agent changes `docs/architecture/` in the same pull request, with the skill `writing-docs`. `docs/principles.md`: a new rule on the docs for people, with its source, enforcement (the skills `writing-docs` and `reviewing-increments`, and code owner review once `CODEOWNERS` names the path) and re-check; it links to `docs/decisions/process/0003-docs-for-people-as-arc42-lean.md`, which records the options and the decision. `README.md`: a line under "Where things are" for `docs/architecture/context.md`, and the status.

## Scenarios

The demonstrated scenarios 05.1 to 05.3 run in the same prepared repo: a small Go module with a command `push` that sends a file to a server, `docs/architecture/context.md` that lists the server as an external interface, and `docs/architecture/runtime.md` with one flow, "Push", whose sketch has 4 nodes: read file, connect, send, report. Each runs as an eval of the skill, with and without it or with the old and the new version; the run without the skill or with the old version is the baseline, and every run with the new skill must pass. The task prompt of 05.1 and 05.2 names the code change only and does not mention docs.

- **scenario 05.1** (demonstrated) **The agent changes the flow that its code changes.** With the skill `writing-docs`, the agent is asked to make `push` try a failed send up to 3 times. Its commits also change the "Push" flow in `runtime.md`: the line with the question still comes right after the heading, the Mermaid sketch has the retry as a node of its own and at most 7 nodes, and the text describes the push as it runs after the change and uses none of the words "new", "now", "changed" and "no longer". `context.md` and every other part of `runtime.md` stay byte for byte as before.
- **scenario 05.2** (demonstrated) **The agent leaves the docs alone when no flow changes.** With the skill `writing-docs`, the agent is asked to rename an internal helper function whose name appears nowhere in `context.md` or `runtime.md`. After its commits, both files are byte for byte as before.
- **scenario 05.3** (demonstrated) **The reviewer finds a flow that no longer matches.** With the skill `reviewing-increments`, the reviewer reviews a commit that adds the retry of 05.1 to the code but leaves `runtime.md` as it was; the design note in the prepared repo asks for the retry. The review result has a finding that the docs for people are out of date, with the file and a line inside the retry code, and with the name of the flow, "Push". The result line of the review result reports findings, so the pull request would wait for the engineer.
- **scenario 05.4** (manual) **The engineer understands the factory from the two docs.** On the branch of the code pull request of this increment, the engineer reads only `docs/architecture/context.md` and `docs/architecture/runtime.md`. The engineer then judges whether, for each of the four flows (the increment workflow, the merge gate, write permissions per session, credentials and network), the engineer can tell who acts at each step and where the work waits for the engineer; whether each sketch answers the question above it and has at most 7 nodes; and whether `context.md` lists each external interface that the factory has, and no other. The scenario passes when the engineer finds no gap. The engineer names each gap in a comment on the pull request, and the agent closes it before the merge.
- **scenario 05.5** (manual) **A change to the docs waits for the engineer.** When `.github/CODEOWNERS` on `main` names `/docs/architecture/` and the ruleset is active, the agent opens a pull request that changes only `docs/architecture/runtime.md`. On GitHub, the engineer sees that the pull request cannot merge until the engineer approves it as code owner, also with auto-merge on. Until both conditions hold, the scenario is not run, and the retro says so.

## Evidence

- **Approved scenarios:** none; the increment changes skills and docs, not the CLI.
- **Demonstrated:** 05.1, 05.2 and 05.3.
- **Manual:** 05.4 and 05.5.

## Rollout

- **Off switch:** remove workflow step 3's line on `docs/architecture/` from `AGENTS.md` and the finding kind from `reviewing-increments`; the two docs then stay as they are and no longer block a merge.
- **Failure signal:** a code pull request changes a flow or an external interface and merges without a change to `docs/architecture/`; the engineer finds a step in `runtime.md` that the code does not run; or the code pull requests of the next 3 increments wait for the engineer only because of `docs/architecture/`.
