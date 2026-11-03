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

## Checkpoints
- 2026-11-03T00:52Z step 1 done; no work branch yet (main@7759cd2, CI PR #41 open); next: step 2; open: none
- 2026-11-03T00:54Z step 2 done; claude/03-setup-script-design@70aa5b5; next: step 3, round 1 of questions; open: required checks before the first run not verified
- 2026-11-03T00:55Z step 3 round 1 sent (form, steps covered, required checks); claude/03-setup-script-design@70aa5b5; next: wait for answers, then round 2 (second run, dry run, report, docs overlap); open: round 1
- 2026-11-03T00:59Z round 1 answered: Q1 factory subcommand; Q2 all four steps; Q3 required checks at once with ruleset.json; claude/03-setup-script-design@70aa5b5; next: round 2 (relation to factory install, what holds of owned docs/setup/, second run and dry run, README overlap); open: round 2
- 2026-11-03T00:59Z step 3 round 2 sent (Q4 relation to install, Q5 owned docs/setup, Q6 second run and dry run, Q7 README overlap); claude/03-setup-script-design@70aa5b5; next: wait for answers, then write design.md; open: round 2
