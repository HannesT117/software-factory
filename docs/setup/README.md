# Setup

The engineer does these steps once per repo, in this order; step 3 needs the engineer's personal account with admin rights on the repo. Options and reasons: [process decision 0002](../decisions/process/0002-approval-through-code-owners-and-an-agent-account.md).

1. **Agent account.** Create a GitHub account for agents, with its own email and two-factor authentication. GitHub allows one free machine account next to a personal account, used only for automated tasks.
2. **Code owners and workflows, before the rulesets.** Copy `.github/CODEOWNERS`, `.github/workflows/` and `.factory/` from the factory repo, put the engineer's personal account in `CODEOWNERS`, and push them to `main`. Once step 3 has imported the ruleset `main`, every change needs a pull request whose required checks pass, and these workflows provide the checks.
3. **Admin settings.** In the factory repo, install `factory` with `go install ./cmd/factory`; it goes to `$(go env GOPATH)/bin`, which must be on `PATH`. In the project's clone, logged in to `gh` with the engineer's personal account, run
   ```bash
   factory setup --agent <agent account> --dry-run
   factory setup --agent <agent account> < file-with-the-opper-api-key
   ```
   It sets the merge settings, invites the agent account with the role Write, creates the environment `cross-review` for `main` with the secret `OPPER_API_KEY`, and imports the rulesets with their required checks. The dry run prints what a run would change. A second run changes only what differs from the target. The values live in `internal/setup/` of the factory repo. Then accept the invitation as the agent account.
4. **Other environments.** Settings → Environments → create each other environment that a workflow names with `environment:`, with "Deployment branches and tags" limited to `main`, and add the secrets that its job uses as environment secrets, never as repo secrets: a repo secret reaches every workflow on every branch, and the agent account may push workflows.
5. **Cloud sessions.** Create a token (classic) of the agent account with the scopes `repo` and `workflow` and an expiry date. Sign in `gh` with it in a separate config directory, start Claude Code there, and run `/web-setup`. It must print `Connected as <agent account>`.
   ```bash
   GH_CONFIG_DIR=~/.config/gh-agent gh auth login --with-token
   GH_CONFIG_DIR=~/.config/gh-agent claude    # then: /web-setup
   ```
   All cloud sessions of the engineer then act as the agent account, also in other projects, so the agent account needs access to each of those repos. The Claude GitHub App does not serve here: it connects only repos that the connected account owns, so Claude Projects threads and auto-fix are not available. Repeat this step before the token expires.

   Two cloud environments: one for design sessions, with open network access and the environment variable `FACTORY_SESSION=design`, and one for all other sessions, whose network access allows only the hosts of `.factory/allowed-hosts`.
6. **Local agent sessions.** Before the first one: merge `adapters/claude-code/user-settings.json` into the Claude Code user settings, and keep its list of hosts equal to `.factory/allowed-hosts`. Give git and `gh` a token (classic) of the agent account with the scopes `repo` and `workflow`; a fine-grained token cannot be used on a repo where its account is only a collaborator. From an agent session, `gh api user` must return the agent account. A local design session starts with `FACTORY_SESSION=design claude`.
7. **Python.** `factory validate-skills` needs Python 3.10 or later with PyYAML as `python3`, on the engineer's computer and in CI.

To check: from a cloud session and from a local agent session, `gh api user` returns the agent account. A pull request that changes a design note waits for the engineer's review as code owner, and a pull request that changes only code merges once its checks pass.

A repo that still uses the environment of [process decision 0001](../decisions/process/0001-approval-through-a-protected-environment.md) switches after step 3: delete `.github/workflows/approval.yml` and the environment `approval`.
