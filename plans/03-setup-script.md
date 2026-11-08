# Plan: increment 03, setup script

- **Folder:** `docs/increments/03-setup-script/`
- **Target pull requests:**
  1. Design pull request with `design.md`, from the branch `claude/03-setup-script-design`.
  2. Tests and code on one branch, `claude/03-setup-script`, after the design note is on `main`.
- **Runs as:** a design session (`FACTORY_SESSION=design`, may read the open web) for steps 1 to 4; then a test-writer session and an implementing session, as `AGENTS.md` describes.

## Context

- `docs/setup/README.md` lists the engineer's one-time setup per repo. Most steps need the engineer's account with admin rights, because the agent account has the role Write only:
  - repo settings
  - inviting the agent account
  - the environment `cross-review` with the secret `OPPER_API_KEY`, which only `main` may use
  - importing the rulesets `docs/setup/ruleset.json` and `docs/setup/ruleset-branches.json`
  - the required checks
- Required checks are added today only after they have run once (step 5 of the setup). Whether the API accepts them earlier, by name and the GitHub Actions app id, is not verified.
- The path groups are defined in the design note of 02, section Format (tests; design notes; CI and config). `factory check-pr`, the guard hook (`factory hook pre-tool-use`) and `.github/CODEOWNERS` use them. Scenarios 02.13 and the path-group scenarios of 02 fix their behaviour; a change replaces them through this note, and `factory scenarios` lists the scenarios in force.

## Decisions of the engineer, 2026-11-03

- A script in `docs/setup/` that the engineer runs once per repo with their own `gh`.
- The engineer owns `docs/setup/` as code owner. `docs/setup/` becomes a path group of its own, so that the guard hook refuses agent edits there and `check-pr` keeps it apart from other changes.
- It is an increment of its own, number 03. Install and distribution becomes 05, the handover 06.

## Decisions of the engineer, round 1, 2026-11-03

- Form: a `factory` subcommand, with the consequence that agents change `cmd/` and `internal/` without the engineer's approval and the engineer runs that code with admin rights.
- Steps: all four through the API: merge settings; invite the agent account; environment `cross-review`, limited to `main`, with the secret `OPPER_API_KEY`; import both rulesets.
- Required checks: at once, with the import of `ruleset.json`. The workflows `checks` and `test` run only on `pull_request`, so the note says nothing about checks on the latest commit of `main`.
- Open from the engineer: "Why don't we use `factory install` for this, then?" Asked in round 2 (Q4), together with what still holds of the owned `docs/setup/` (Q5), second run and dry run (Q6), and README overlap (Q7).

## Decisions of the engineer, round 2, 2026-11-03

- Q4: two commands. `factory setup` (admin, the engineer, once per repo) in increment 03; `factory install` (files, agents, every update) in increment 05. Increment 03 stays.
- Q6: each run sets every step to the target state again; rulesets are updated by their name. A dry run prints what it would change without changing anything.
- Q7: `docs/setup/README.md` lists only the steps that stay manual and names the command for the rest; the values live only in the command and the ruleset files.
- Q5 (what holds of the owned `docs/setup/`): open. The engineer asked whether all of `docs/` should be owned; the orchestrating session answers that and sends Q5's answer.

## Decisions of the engineer, Q5, 2026-11-03

- Ownership of `docs/`: as today, plus `docs/setup/`. Not all of `docs/`. The CODEOWNERS line for `docs/setup/` comes with a CI pull request, not with the design pull request.
- Rulesets: embedded in the binary from a code path that the engineer does not own. Protection: approved scenarios with tests, reviewer, cross-review, and the dry run the engineer reads before applying.
- Agent decision: the ruleset files move to `internal/setup/rulesets/` (Go embed reaches only files below the package folder); `docs/setup/ruleset*.json` are removed.
- Open, round 3 (Q8): whether `docs/setup/` stays a path group of its own when it holds only `README.md`.

## Decisions, round 3, 2026-11-03

- Q8, decided by the orchestrating session for the engineer: `docs/setup/` is owned but not a path group. The code pull request of this increment updates `docs/setup/README.md`, removes the two ruleset files, and waits for the engineer's approval as code owner.
- From now on the design session settles details itself and records each in the note or the pull request; it asks the engineer only about direction: scope, who owns or approves what, and cost.

## Open, to ask the engineer in rounds

Each question names the consequences of each option and carries no recommendation (skill `designing-increments`, section Who decides):

- Form of the script: a shell script with `gh`, or a `factory` subcommand. A subcommand lives in `cmd/` and `internal/`, which agents change without the engineer's approval.
- Which steps the script covers, and which stay manual (agent account, cloud environments, tokens).
- Behaviour on a second run (idempotent update), a dry run that prints the changes, and how it reports each step.
- How the script and `docs/setup/README.md` avoid stating the same thing twice (`AGENTS.md`: a doc points to a config file instead of restating it).
- The required checks: in the script, or still by hand after the first run.

## Steps

1. Read `AGENTS.md`, `README.md`, `docs/principles.md`, `docs/setup/` with its files, the design note of 02 (Format, Interface, the scenarios on path groups and the guard hook), and `core/skills/designing-increments/SKILL.md` with its assets.
2. Research existing tools, with links, for the section Existing tools: the GitHub REST API for repo settings, collaborators, environments, deployment branch policies, environment secrets, rulesets and required status checks; `gh` commands; tools that apply repo settings from a file, such as the Probot settings app, safe-settings, the Terraform GitHub provider, or `gh` extensions. Record each tool checked, the decision and the reason.
3. Ask the open decisions in rounds, through the orchestrating session.
4. Write `design.md` with the skill `designing-increments`. Run `factory check-design docs/increments/03-setup-script`. Have a fresh context explain each scenario, and rewrite each scenario that it cannot explain alone. Open the design pull request (skill `writing-pull-requests`), turn on auto-merge, and stop. The engineer approves as code owner.
5. After the merge: the test writer, the code, the review and the retro, as `AGENTS.md` describes. Each step gets its own plan section here before it starts.

## Step 5a: acceptance tests

- Runs as: a test-writer session (agent `acceptance-test-writer`), on the branch `claude/03-setup-script` from `main`.
- One `testscript` file per automated scenario, 03.1 to 03.8, in `test/acceptance/03-setup-script/`. The runner gets a stand-in for `gh` (`test/acceptance/gh_test.go`) and two JSON commands (`test/acceptance/json_test.go`).
- Done when every test fails because `factory setup` does not exist, the tests of 02 pass, and `factory check-tests --base main` passes. Then the pull request "Add the command for the repo admin setup", and a report to the orchestrating session.

## Step 5b: code

- Runs as: the implementing session, on the branch `claude/03-setup-script` (pull request #54).
- `factory setup` in `cmd/factory` and `internal/setup/`, with the rulesets in `internal/setup/rulesets/` (moved from `docs/setup/`); `docs/setup/README.md` with only the manual steps; rules 16 and 42 point to the new place.
- Done when `go test ./...` passes; then the subagent `increment-reviewer` until it finds nothing, the retro, the pull request description and the status row of 03 in `README.md`, and a report to the orchestrating session.

## Checkpoints
- 2026-11-03T00:52Z step 1 done; no work branch yet (main@7759cd2, CI PR #41 open); next: step 2; open: none
- 2026-11-03T00:54Z step 2 done; claude/03-setup-script-design@70aa5b5; next: step 3, round 1 of questions; open: required checks before the first run not verified
- 2026-11-03T00:55Z step 3 round 1 sent (form, steps covered, required checks); claude/03-setup-script-design@70aa5b5; next: wait for answers, then round 2 (second run, dry run, report, docs overlap); open: round 1
- 2026-11-03T00:59Z round 1 answered: Q1 factory subcommand; Q2 all four steps; Q3 required checks at once with ruleset.json; claude/03-setup-script-design@70aa5b5; next: round 2 (relation to factory install, what holds of owned docs/setup/, second run and dry run, README overlap); open: round 2
- 2026-11-03T00:59Z step 3 round 2 sent (Q4 relation to install, Q5 owned docs/setup, Q6 second run and dry run, Q7 README overlap); claude/03-setup-script-design@70aa5b5; next: wait for answers, then write design.md; open: round 2
- 2026-11-03T01:06Z round 2 answered except Q5; claude/03-setup-script-design@70aa5b5; next: write design.md parts that do not depend on Q5; open: Q5
- 2026-11-03T01:07Z draft of design.md without Q5 (check-design ok); claude/03-setup-script-design@974f13a; next: add the Q5 answer (ruleset source, docs/setup path group), fresh-context scenario check, PR; open: Q5
- 2026-11-03T01:09Z fresh-context check of scenarios 03.1-03.9 done, scenarios rewritten; claude/03-setup-script-design@5ad835e; next: add the Q5 answer, check new scenarios with a fresh context, open the PR; open: Q5
- 2026-11-03T01:36Z Q5 answered; claude/03-setup-script-design@5ad835e; next: round 3 (Q8 path group of docs/setup), meanwhile write the Q5 parts into design.md; open: Q8
- 2026-11-03T01:37Z round 3 sent (Q8 path group of docs/setup); Q5 written into design.md; claude/03-setup-script-design@6dab254; next: write the Q8 answer, fresh-context check of changed scenarios, open the PR; open: Q8
- 2026-11-03T01:44Z step 4 done: design PR #45 open (https://github.com/HannesT117/software-factory/pull/45); claude/03-setup-script-design@f32707c; auto-merge refused because the PR is already clean, so it is not merged by the agent and waits for the engineer; next: after the merge, step 5 (test writer); open: none
- 2026-11-03T02:06Z step 5a: gh stand-in and tests 03.1 to 03.3 done; claude/03-setup-script@df3d2c6; next: tests 03.4 to 03.8; open: none
- 2026-11-03T02:07Z step 5a: test 03.4 done; claude/03-setup-script@9a1936d; next: test 03.5; open: none
- 2026-11-03T02:07Z step 5a: test 03.5 done; claude/03-setup-script@0ad7884; next: test 03.6; open: none
- 2026-11-03T02:08Z step 5a: test 03.6 done; claude/03-setup-script@531ae03; next: test 03.7; open: none
- 2026-11-03T02:08Z step 5a: test 03.7 done; claude/03-setup-script@0c10a4c; next: test 03.8; open: none
- 2026-11-03T02:08Z step 5a: test 03.8 done; claude/03-setup-script@70ff719; next: pull request; open: none
- 2026-11-03T02:09Z step 5a done: tests 03.1 to 03.8 pushed, each fails because factory setup is missing; pull request #54 (https://github.com/HannesT117/software-factory/pull/54); claude/03-setup-script@70ff719; next: step 5b, the implementing session; open: none
- 2026-11-03T04:01Z step 5b: code and docs done, go test ./... passes; claude/03-setup-script@62b966c; next: increment-reviewer; open: none
- 2026-11-03T04:04Z step 5b: review 1 done (block, 1 finding: ruleset compared by containment, not equality); README fix pushed; claude/03-setup-script@5a0ca33; next: retro draft, then fix and re-review after the engineer's answer; open: ruleset comparison A/B/C sent to the orchestrating session
- 2026-11-03T04:09Z step 5b done: code, review (2 rounds), retro, PR description and README status on #54 (https://github.com/HannesT117/software-factory/pull/54); claude/03-setup-script@e2243cf; reported to the orchestrating session; next: fix of the required check checks on main (own PR, merge ref), merge of #56, manual scenario 03.9 by the engineer; open: checks fails on #54 until main is fixed
- 2026-11-08T23:46Z step 5b: paging (cross-review finding) done; review 3: no findings; claude/03-setup-script@5ef86b2 on GitHub, review d11d701 local in session_01PcDxYGeTSySmjYdvgHL244; next: push d11d701, update the PR description, merge main into the branch (orchestrator); open: GitHub access 403, the account was suspended
