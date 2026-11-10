# AI software factories: what the field says, and a review of this factory

Research and review, 2026-11-10, against `main` at 2461855. Part 1 collects what current sources say about setting up a software factory around coding agents, about writing skills, and about how others do it. Part 2 holds this factory against it. Part 3 lists recommendations in priority order. Sources are at the end.

This note builds on [docs/research/2026-10-factory-basis.md](2026-10-factory-basis.md) and repeats a source from it only where the review needs it.

**How the review was checked.** The docs, skills, workflows and the Go packages `prcheck`, `paths`, `review`, `guard`, `crossreview` and the Claude Code adapter were read. `factory` was built from 2461855 and `check-tests`, `check-adr` and `scenarios` were run. Claims about the code and the repo are measured. Two things could not be determined from the repo and are marked as open questions: whether the rulesets of `factory setup` are active on GitHub, and where the hooks of the factory repo itself are configured. The first was read through the GitHub API on 2026-11-10: the ruleset `main` is active, but requires neither a code owner review nor a status check.

## Summary

- **The centre is right.** Verification first, deterministic checks before model judgment, tests written apart from the code and before it, a judge that is not the worker, the engineer approving intent on the server side. This is the field's consensus, and the factory records the limits of each practice better than most of its sources do.
- **The gate is weaker than the docs read.** Three gaps: the code of the checks, the hooks and the reviewer's skill can change in a pull request that merges without the engineer (F2); no hard check reads the result of the fresh-context review, so a result with findings does not hold a merge (F3); and the scenarios that prove the hard guards are manual, and none is recorded as run (F1). On `main`, `check-tests --base main` fails with 7 errors, so a required check is red for every code pull request, and the repo does not say how the last two code pull requests merged.
- **The process costs more than it has to.** Seven reviewer rounds and four cross-review rounds for a change to two skills, with no cap and no severity on reviewer findings (F4). Gate logic is written twice, in Go and in workflow shell (F5).
- **The factory has grown ahead of its evidence.** 42 rules in two weeks, 15 of them enforced, one of them sourced from an observed failure, and no project has run an increment on the factory yet (F7). For four of six skills, the evidence that the skill helps is incomplete (F6).
- **First moves:** record and close the manual scenarios of the trust core; make the state derivable with `factory status`; own the trust roots; give reviewer findings a severity that a check reads, and cap the rounds; move the gate predicates into the CLI.

## Part 1: What the field says

### 1.1 Verification is the constraint

- Thoughtworks' retreat report, as Fowler quotes it: "Code generation is no longer the bottleneck — verification is." Addy Osmani: "Verification, not generation, is the real constraint on a factory", and a loop gets "as much autonomy as you can cheaply and reliably verify, and not one inch more." A loop earns "fully automated status" only when its check "is cheap, runs at high frequency, and relies on something that can't be easily faked out."
- Böckeler's harness article gives the vocabulary: guides (feedforward) and sensors (feedback), each computational (deterministic, cheap) or inferential (model judgment, "expensively and probabilistically. Not on every commit"). Feedback alone repeats mistakes; guides alone never learn whether they worked.
- The behaviour harness is the unsolved part: a green AI-written test suite "puts a lot of faith into the AI-generated tests, that's not good enough yet." Fowler: "Conformance tests (sensors) are more valuable than specifications (guides), but it's hard to imagine all the conformance tests that are needed to say what shouldn't happen", and "We can outsource many things, but not the acceptance criteria."
- Naresh Jain measured a behavioural harness on one service, 10 runs per arm: a frontier model passed 10 of 10 with and without it, at 10.5 times the tokens with it; a cheaper model rose from 3 of 10 to 10 of 10. His rule: keep an independent acceptance check in every case, and for a strong model run it "once as a deterministic final verification gate."

### 1.2 Where the human stands

- Kief Morris: humans "in the loop" inspect every line and become the bottleneck; humans "on the loop" build and maintain the harness. The in-the-loop reflex "is to fix the artefact", the on-the-loop one "is to change the harness that produced the artefact."
- Osmani's lit factory "moves the point of human judgment upstream, to the product, the design, and the architecture before an agent starts a loop." The lights stay on for "large blast radii" and decisions that shape a year of work. The review gate is the one box that does not scale. His failure modes: comprehension debt ("three to six months into a project, you're already drowning in unread code"), and a factory run dark for about four months that had to be torn down.
- Technology Radar volume 35 (November 2026): "risk-tiered human oversight" at Trial ("The riskiest or hardest-to-reverse changes are reserved for humans to apply"); "human-in-the-loop theater" at Caution ("an empty approval step rather than a real check"). Fowler reports Zalando: an LLM rates the rollout risk of each pull request, low risk is auto-approved, lead time falls by 20 to 40 percent, and configuration changes count as high risk.
- Fowler, Agentic Engineering: "Humans are still responsible for what the software does and how it works." His Interrogatory LLM pattern has a model interview the human instead of handing over a document to review, because "people often find reviewing hard."

### 1.3 Instructions and context

- Böckeler sorts context features by who decides to load them: the model (skills, subagents), the human (commands), or the software (hooks, path-scoped rules). Only the third is deterministic. Her advice: "build context like rules files up gradually, and not pump too much stuff in there right from the start", and expect that what was needed half a year ago is no longer needed. She names the failure "illusion of control."
- Claude Code best practices: "If you keep adding emphasis to a CLAUDE.md line and Claude still misses it, that's a sign the line belongs in a hook. The reverse holds too: if you have hooks for things that are only preferences, delete them and let Claude decide." "Keep the enforced set small." An instruction file targets under 200 lines. "A reviewer that's asked to find problems will always find something."
- Radar volume 35 puts "agent instruction bloat" and "AI-driven documentation sprawl" at Caution. Böckeler on spec tools: "I'd rather review code than all these markdown files"; spec-first is "definitely valuable", spec-as-source risks "inflexibility and non-determinism".

### 1.4 Skills: the authoring rules, and this repo against them

From Anthropic's skill authoring best practices. The right column is the state of the six skills in `core/skills/`.

| Rule | This repo |
|---|---|
| `name`: at most 64 characters, lowercase letters, digits, hyphens, no `anthropic` or `claude`; gerund form suggested ("Consider using gerund form") | Met. Errors and the warning `name-not-gerund` in the validator; all six are gerunds. |
| `description`: at most 1024 characters, third person, what the skill does and when to use it, with the terms a person would type | Met, and each also says what it is not for. |
| "Claude is already very smart": only what the model cannot know | Met. Bodies are 23 to 47 lines. |
| Body under 500 lines; references one level deep; a table of contents in a reference over 100 lines | Met, with validator warnings for each. |
| Degrees of freedom: heuristics for judgment, an exact script for fragile or consistency-critical steps | Mostly met. The skills name `factory check-*` as the authority on formats, and the stop hook prints the commit a review must name. |
| Feedback loop: run validator, fix, repeat | Met. |
| No time-sensitive text, one term per concept, forward slashes | Met, with validator warnings. |
| Evaluations first: at least three, with a baseline without the skill | Met in form: 3 to 5 cases per skill, 10 trigger queries each. Incomplete in results, see F6. |
| "Test your Skill with all the models you plan to use it with" | Not met. Evals ran on a mid-size model with a small judge; the sessions run on a larger model. |
| Iterate by watching a fresh instance use the skill on real tasks | Partly. The retros record corrections; scenario 02.16 is open. |
| "You don't need special system prompts or a 'writing skills' skill" | Deliberate deviation: `authoring-skills`, and a hook that requires loading it, see F6. |

### 1.5 Swamp Club

Three layers, which the sources keep apart.

- **The product stance.** "Behavior-centric design separates understanding from execution": typed, versioned behaviors compose into deterministic workflows that "call an LLM only when a specific step requires reasoning." The published numbers are vendor figures. In their code-review case study, a skill with a coordinator and 23 subagents (4.6 million tokens, 12.8 minutes) became a workflow with 3 agents (0.5 million tokens, 6.6 minutes); "Most of the savings came from removing the coordinator agent." The second case on the same page fails in both versions.
- **The `@swamp/software-factory` extension.** A generic state machine: "stages, transitions, gates, required artifacts, review skills, agent commands, system prompts" are data, and the engine "ships no lifecycle concepts." Gates are checked in code on every advance: `artifact-exists`, `artifact-fresh`, `findings-clear` with blocking severities, `human-approval` with the approver's identity, `evidence-recorded`, `cooldown`, `workflow-succeeded` ("verified against swamp's own workflow run records, not driver attestation"). Every stage has `maxCycles`, default 5; past it, the run parks for a human override. A journal is append-only, and a summary is rendered without a model. In both example lifecycles the human approves the plan and the result, never the tests, and a review stage clears only when no `critical` or `high` finding is open.
- **How Swamp builds Swamp.** Issue-driven; outside pull requests are closed, because agents "produce working, reliable source code at a rate that is difficult for a human to understand or review." Lint, tests, build and agent reviews run locally in fresh worktrees, and a generated attestation tells CI that they ran. CI itself runs an adversarial review for core changes, a security review for workflow changes and a "review integrity" check for trust-root changes. Pull requests auto-merge unless they carry the label `hold`. Swamp's docs state the weak point: "nothing yet stops a document being posted that it did not write." Skills are scored in CI: a review score of at least 90 percent, and trigger evals with should-trigger and should-not-trigger queries at a 90 percent pass rate, on the default model per pull request and on five models weekly. Acceptance tests live in a separate repo and run against the compiled binary.

### 1.6 A harness goes stale, so it is measured

- Radar volume 35: "harness quality evaluation" at Assess (guardrails help "but only if they actually work"), "feedback flywheel" at Trial, "coding throughput as a measure of productivity" at Caution. Cole Murray's measures: "cycle time, accepted-output yield, escaped defects, review load, and cost per accepted change."
- Böckeler's open questions stand: how a harness stays coherent as it grows, and "If sensors never fire, is that a sign of high quality or inadequate detection mechanisms?" One open factory answers the second mechanically: every rule is "written twice: once as prose, once as an enforced check with a mutation test proving it fires."
- The counter-position, which Fowler quotes from Robert Martin: with stronger models "the need for any but the most liberal of harnesses may be obviated." Jain's frontier-model result points the same way.

### 1.7 Sensors for maintainability

Böckeler, May 2026, on a TypeScript codebase, without measurements:

- Computational sensors "impressed me most at the file and function level": type checks, lint thresholds for length and complexity, dependency rules. Lint messages are written as guidance for the agent, with the reason and the allowed way out, and dependency errors recap the layering.
- Mutation testing found 13 surviving mutants in a file with 100 percent statement coverage: "coverage tells us that a line was executed, but not that its impact was verified." And: "acceptance tests increase coverage, but are often not very assertion-heavy."
- Cross-file problems (duplication, inconsistent calls) were found by a model review, run at a slow cadence; "computational sensors alone cannot help us much" there. Fowler's retreat notes agree: "we often find duplication in LLM generated code, together with mixing of concerns", also with a good harness.
- Risks: feedback overload, agents that raise thresholds, and "a false sense of security and an illusion of quality."

## Part 2: Review of this factory

### 2.1 What matches or exceeds the field

- **Centre and order.** Rules 1, 2 and 5 are the consensus of 1.1. The cross-vendor review as a single final gate is Jain's recommendation, and rule 13 records the study in which such a reviewer made code worse.
- **Tests before code, by a check.** `check-pr` and `tests-start-red` prove on the server that the tests came first and failed. No reviewed source enforces this order; most ask for it in a prompt.
- **Checks built from `main`.** `checks` and `cross-review` build `factory` from the base branch, the cross-review key sits in an environment that only `main` uses, and config is read from the base. Swamp's attestation model trusts the agent's machine and says so; this model does not.
- **Approved scenarios in the engineer's words,** with a fresh context explaining each one before the pull request. This is the spec form Böckeler recommends, and the cautions of 1.3 about generated specs do not apply to it.
- **Tests as data.** `testscript` files with a fixed first line and no `skip`, `stop` or conditions; the Radar lists the technique at Assess.
- **The assumption columns.** "Depends on" and "Re-check when" in `docs/principles.md` put 1.6 into practice. No reviewed source records which rules are tied to the current model.
- **One source per rule in the skills.** `authoring-skills` keeps a trace table and ships `check_sources.py`; the other skills name a `check-*` command as the authority instead of restating formats.
- **Retros with numbers.** Review rounds per source, interventions, eval runs and their cost, and the agent version are recorded today, ahead of rule 20's plan.
- **Evidence labels.** The research note says how each source was read, and reports name the kind of each piece of evidence.

### 2.2 Findings

**F1. The evidence for the hard guards is open, and the state of the gate cannot be read from the repo (high).**

- 48 scenarios are in force: 35 automated, 4 demonstrated, 9 manual. No file records the result of a manual scenario. Open are 02.10 (only the engineer approves), 02.11 (a pull request cannot fake a check), 02.12 (agent identity and credentials), 02.20 (new tests start red), 02.27, 02.16, 03.9, 04.4 and 04.5. The first four are the claims that rules 16, 31, 33 and 34 rest on. Demonstrated scenario 02.15 is recorded as "not complete". Increment 02 has no retro.
- `factory check-tests --base origin/main` on `main` reports 7 errors: the six missing tests of increment 05 and the replaced test of 02.7. The required check `checks` therefore fails for every pull request that changes code. Two such pull requests merged after that state began (#65, #59). The design of increment 02 says the next retro records each use of break glass; the retros of 03 and 04 record none. Open question: was the ruleset `main` active for those merges?
- Eight entries in `docs/principles.md` still read "planned (CI pull request of increment 02)", though the workflows, `CODEOWNERS` and the rulesets are on `main`. The README rows of increments 02, 04 and 05 describe the state before their pull requests merged. A status written inside a pull request is one step behind once it merges.
- Field: the Radar's "human-in-the-loop theater" starts where an approval or a manual step leaves no record. Swamp's `evidence-recorded` and `human-approval` gates exist for this.
- Do: record each manual and demonstrated result in a file under `evidence/<nn>-<slug>/` that the engineer owns, one line per scenario with date and result, and each break-glass use beside it. Close 02.10, 02.11, 02.12 and 02.20 before `factory install` hands these guards to a project.

**F2. The factory repo's trust roots are not owned (high).**

- `CODEOWNERS` covers `.github/`, `.factory/`, `go.mod`, the notes and the principles. It does not cover `internal/` (the checks themselves, the path groups in `internal/paths`, the review instructions in `internal/crossreview`, the rulesets in `internal/setup/rulesets/`), `adapters/claude-code/` (hooks, agent prompts) or `core/skills/` (the reviewer's skill).
- CI builds `factory` from `main`, so a pull request cannot change the checks that judge it. It can change the checks that judge the next one, and merge without the engineer when its own checks pass. Rule 34 holds for one pull request, not for two in a row.
- In a project this gap is closed by design: the factory version and checksum are pinned in an owned file (ADR 0002). It is open only in the factory repo, which is the one place where code is the gate.
- Field: Swamp adds a separate review for "trust-root changes"; the Radar reserves the riskiest changes for humans; Osmani keeps the lights on where the blast radius is large.
- Do: add the gate packages, the adapter and the skills of the reviewer and the test writer to `CODEOWNERS` in the factory repo. The cost is one approval per factory code pull request, which is the right tier for this repo.

**F3. The fresh-context review does not hold a merge (high).**

- `review.Pending` and the job `cross-review` ask only whether some `evidence/*/review.md` contains the hash of the latest code commit. Nothing reads its `Result:` line. The result of increment 04 on `main` reads "Result: 1 finding", and its pull request merged.
- The skill says "any finding makes it wait", but what makes it wait is the reviewer not turning on auto-merge. The implementing session may merge "once GitHub allows".
- The review result lies under `evidence/`, which goes with any path group and is not owned, and the guard hook matches the Write and Edit tools, not the shell. An implementing session can write the file.
- Test authorship has the same shape: the order of tests and code is hard, who wrote the tests is not, because CI sees one agent account. The design note of increment 02 lists both as soft parts; rule 12 and rule 28 read as if they were enforced.
- So of the two inferential reviews, only the other vendor's is a gate. That may be acceptable (rule 13 has the reasons), but then the docs should say it.
- Do: either state in rules 12 and 28 that these are early feedback, or make it hard with a small check that reads the result of the latest code commit and fails on a blocking finding (with F4).

**F4. Review rounds have no cap, and reviewer findings have no severity (high).**

- Increment 04, two skills and their evals: seven reviewer rounds, four cross-review rounds, about forty eval runs. From the third round on, the retro calls the findings "doubtful"; several concern the wording of the retro and the README. Increment 03: four rounds, and a new branch with the code cherry-picked, because a test could not change after the first code commit.
- The skill asks for this: "a missed finding costs more than a doubtful one." The cross-review has `blocking` and `note`; the reviewer has only findings. Every pull request with code pays the same price, whatever it could break.
- Field: "A reviewer that's asked to find problems will always find something", so the prompt says what to skip. Swamp clears a review stage on blocking severities only and parks a run after `maxCycles`. The Radar tiers oversight by risk. Rule 13 already says once per pull request.
- Do: give the reviewer the two severities of the cross-review; only `blocking` holds the merge, notes go into the report. After two rounds on one pull request, the third is a decision for the engineer. Name what the reviewer skips: wording in retros and in the README. Run no cross-review on a pull request that changes no code path. Have the cross-review read the tests against the note at the commit where the tests end, so that a missing "does not" case surfaces before the first code commit.

**F5. Gate logic is written twice, in Go and in workflow shell (medium).**

- `checks.yml` turns `.factory/code-paths` into a regular expression with `sed` to decide whether `check-tests` runs. `paths.Code.Is` decides the same in Go, for the stop hook. `cross-review.yml` decides with `grep` whether a pull request is a design pull request; `paths.GroupOf` decides it in Go. `prcheck.readCommits` and `review.Pending` each walk the first-parent commits and filter merge commits by the branch's paths.
- `adapters/claude-code/user-settings.json` and `.factory/allowed-hosts` hold the same host list, kept equal by hand (setup step 6).
- The shell copies run in the most trusted place, have no acceptance test, and the engineer has to review a `sed` program by eye. When two places must agree on a predicate, the next change moves one of them.
- Do: let the commands decide for themselves. `factory check-tests --base` prints that it does not apply when the pull request changes no code; `factory cross-review` picks its own mode and waits for the review result itself. The workflows shrink to calls, and the predicates get scenarios. Generate the sandbox host list from `.factory/allowed-hosts` at install.

**F6. For four of six skills, the evidence is incomplete (medium).**

- Scenario 02.15 stopped at a usage limit: two skills ran, four did not, and no trigger query ran.
- The retro of increment 04 found that the eval judge sees only the first and last 12 lines of a trace. 75 of 92 grader files are still judged by a model on the trace; 17 match tool calls. The fix reached five cases of two skills. `authoring-skills` still tells the agent to leave tool-call graders alone, which the retro marks as wrong.
- Runs per case are 1 or 2, on one model that is not the model of the sessions. Evals run by hand, not in CI and not on a schedule.
- The guard hook requires that a session loads `authoring-skills` before it edits a skill, by parsing the transcript of one agent tool. Anthropic's page says no such skill is needed, and the Claude Code guidance says to delete hooks for preferences. The deterministic value is in the two scripts; `quick_validate` runs in CI, `check_sources` does not.
- Do: move the remaining trace graders to tool-call graders and correct `authoring-skills`; finish 02.15 on the model the sessions use, with three runs per case; rerun when the agent version changes (rules 21 and 41 name that trigger). Run `check_sources` in `checks` for changed skills, then drop the loading requirement from the hook.

**F7. Rules have grown faster than observed failures, and no project has used the factory (medium).**

- 42 rules after 45 commits in 16 days: 9 enforced, 6 enforced with a part still planned, 12 planned, 15 not enforced. Rule 3 says a rule is added only after an observed failure; one rule (27) names an observed failure as its source. The backlog holds 22 more parts.
- Every increment so far is the factory building itself. `factory install` is increment 06 of 08. Rule 3 also says a factory fits one product; so far the product is the factory.
- Field: Böckeler builds context up gradually; the Radar warns of instruction bloat; Jain measured a tenfold token cost with no gain for a frontier model; Robert Martin expects most harness parts to become unnecessary.
- Do: run one real project increment on the factory before adding process, with a hand copy of the skills and workflows if install is not ready. Until then, add no rule without an observed failure. Give "not enforced" a clock: after two retros such a rule names its check and increment, moves to the backlog with a trigger, or shrinks to a one-line principle without a re-check.

**F8. The acceptance tests are the only sensor on the code, and their strength is not measured (medium).**

- About 3,050 lines of Go in `cmd/`, `internal/` and `adapters/` have no unit tests, by design, and `go vet` as the only static check. The 35 acceptance tests run a separately built binary.
- The backlog parks the mutation score because coverage-based tools see nothing of that binary. Go measures coverage of a built binary since 1.20 (`go build -cover`, `GOCOVERDIR`), so the coverage half needs only the standard toolchain. The reviewer of increment 03 did the mutation half by hand, with three reverts.
- Field: 1.7. The Radar has mutation testing at Adopt.
- Do: build the binary under test with `-cover` and report which packages and functions no acceptance test reaches. Keep the mutation runner in the backlog with its trigger. When increment 06 builds lint rules with fix-naming messages for projects, apply them to the factory's own Go as well.

**F9. The current behaviour is spread over notes that amend each other (low now, grows).**

- The Interface of increment 05 reads "as in increment 02, with this change." `factory scenarios` lists the titles in force, not the text, and nothing assembles the interface in force. A test writer for a later increment reads every earlier note that touches the same command.
- Field: the Radar's documentation sprawl; Böckeler's spec-anchored level, where the kept spec has to stay one thing.
- Do: decide this in increment 07. A generated view from the notes fits rule 36 better than a second document kept by hand.

### 2.3 What the field does not have, and the factory should keep

- **Comprehension questions before the design approval,** from a session that knows only the note, with the engineer answering first, and a deviation list after the code. It is the one mechanism in the reviewed material that keeps the human able to debug the system. Fowler's Interrogatory LLM is the nearest pattern, and it supports the direction.
- **The "Existing tools" section** in every design note, with rejected options and reasons.
- **Reports that start with the pull requests** and name the kind of each piece of evidence.
- **A known residual, kept as is:** one machine account with a classic token serves all of the engineer's sessions, across projects. Rule 33 names the re-check; a second project that holds anything private is that moment.

## Part 3: Recommendations in order

| # | Recommendation | Finding | Changes |
|---|---|---|---|
| 1 | Record manual and demonstrated results, and each break-glass use, in owned files under `evidence/`; close 02.10, 02.11, 02.12 and 02.20 before increment 06 | F1 | `CODEOWNERS`, `AGENTS.md` step 5, retro of 02 |
| 2 | Pull `factory status` forward: per increment the stage, open manual scenarios, increments waiting for tests, the commit a review must name, what waits for the engineer. The README status points to it | F1 | `cmd/factory`, backlog, README |
| 3 | Own the trust roots of the factory repo | F2 | `.github/CODEOWNERS`, rule 34 |
| 4 | Reviewer findings get `blocking` or `note`; a check reads the result of the latest code commit; two rounds, then the engineer decides; a skip list for the reviewer; no cross-review without a code path; cross-review of the tests before the code | F3, F4 | `reviewing-increments`, `checks`, rules 12 and 13 |
| 5 | Move the gate predicates out of the workflows into the commands, and add scenarios for them | F5 | `checks.yml`, `cross-review.yml`, `internal/` |
| 6 | Finish the skill evidence: tool-call graders, 02.15 on the session model with three runs, rerun on an agent version change; `check_sources` in CI, then drop the loading requirement from the hook | F6 | skills' `evals/`, `authoring-skills`, `internal/guard`, scenario 02.33 |
| 7 | One table per retro with a row per review source: findings, findings that changed code or a note, cost. After three increments, leave the source with the lowest yield out once (rule 22) | F4, F7 | `retro.md` |
| 8 | A project increment before more process; no rule without an observed failure until then; a clock for "not enforced" rules | F7 | order of increments, rule 4 |
| 9 | Coverage of the acceptance run with `go build -cover`; lint for the factory's own Go with increment 06 | F8 | `test/acceptance/acceptance_test.go` (tests group), `test.yml` |
| 10 | A generated view of the interface and scenarios in force | F9 | scope of increment 07 |

Considered and dropped:

- **Generator commands beside the checks** (`new-increment`, `new-test`): the retros show no failure on formats; the test writer asked no question in increment 03, and the checks name each mistake. Three commands for no observed problem.
- **Holdout scenarios:** EvilGenie found little gain; stays in the backlog with its trigger.
- **No human review, as at StrongDM, or auto-merge by default with a `hold` label, as at Swamp:** one engineer has to keep the overview, and the design approval is the place for it.
- **A dependency on Swamp, or on Spec Kit or OpenSpec:** the Existing tools table of increment 02 gives the reasons, and the ideas of 1.5 transfer without the tool.
- **A lifecycle state machine as data:** the state of an increment can be derived from git and the evidence files; a status command reads it without a second store.
- **A shorter `AGENTS.md`:** 63 lines under rule 41 is already the short map the field asks for.

## Sources

**read**: fetched and read for this review through a fetch agent. On 2026-11-10 each quote and claim about the Anthropic, Claude Code, martinfowler.com, Thoughtworks, Swamp, Osmani, Go and GitHub pages was checked again against the page; the wording above follows the pages where it differed. **vendor** and **advocacy** as in the research note. Not retrieved: OpenAI, Harness engineering (HTTP 403); the two Thoughtworks retreat reports (PDF), taken through Fowler's fragments.

Skills and instructions

- Anthropic, Skill authoring best practices, platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices (read; vendor)
- Anthropic, Best practices for Claude Code, and the memory guide, code.claude.com/docs/en/best-practices (read; vendor)

Swamp Club (all read; advocacy)

- `@swamp/software-factory` 2026.10.23.1, swamp-club.com/extensions/@swamp/software-factory; its README, `SKILL.md`, and the examples `feature-factory.yaml` and `sdlc-classic.yaml`, github.com/swamp-club/swamp-extensions
- Press release, 2026-10-14, and Behavior-centric AI design, swamp-club.com/behavior
- Code review token reduction case study, swamp-club.com/case-studies/code-review-token-reduction
- github.com/swamp-club/swamp: README, `CONTRIBUTING.md`, `AGENTS.md`, `CLAUDE.md`, `contributing/skills-pipeline.md`, `agent-constraints/verification-conventions.md`

martinfowler.com (all read)

- Böckeler: Harness engineering for coding agent users, 2026-04-02; Maintainability sensors for coding agents, 2026-05-27; Context engineering for coding agents, 2026-02-05; Understanding spec-driven development, 2025-10-15
- Kief Morris, Humans and agents in software engineering loops, 2026-03-04
- Naresh Jain, Are harnesses worth the cost?, 2026-10-13
- Fowler: Agentic Engineering, 2026-05-21; Interrogatory LLM, 2026-05-14; Future of Software Development, 2026-07-02; fragments of 2026-07-06, 07-13, 07-21, 08-24 and 09-16

Factories and measures

- Thoughtworks Technology Radar, volume 35, November 2026, thoughtworks.com/radar: the blips named in the text (read; the lower part of the techniques list through a summary)
- Addy Osmani, Software factories, light and dark, 2026-07-20, addyosmani.com/blog/software-factories (read; the author works on Claude Code)
- Simon Willison on StrongDM's software factory, 2026-02-07 (read)
- Cole Murray, What is a software factory?, 2026-07-08, murraycole.com/posts/software-factory (read)
- varun1505, awesome-software-factories, the entry for nicolasmelo1/software-factory (read; the project itself was not)
- LeadDev, The engineer's guide to building a software factory, 2026-09-10, and Warp, Software factory vs. agentic loop, 2026-08-13 (read; vendor content, not used for a finding)
