# Setup

These steps need the engineer's personal account with admin rights on the repo. Do them once per repo, in this order. Options and reasons: [process decision 0002](../decisions/process/0002-approval-through-code-owners-and-an-agent-account.md).

1. **Repo settings, before the first merge.** Settings → General → Pull Requests:
   - "Allow squash merging": on, with the default commit message "Pull request title and description".
   - "Allow merge commits" and "Allow rebase merging": off.
   - "Allow auto-merge": on.
   - "Automatically delete head branches": on.
2. **Agent account.** Create a GitHub account for agents, with its own email and two-factor authentication. GitHub allows one free machine account next to a personal account, used only for automated tasks. Settings → Collaborators → invite it with the role Write, never Admin.
3. **Code owners.** Copy `.github/CODEOWNERS` from the factory repo, put the engineer's personal account in it, and merge it.
4. **Workflows.** Copy `.github/workflows/` and `.factory/` from the factory repo. Settings → Environments → create each environment that a workflow names with `environment:`, with "Deployment branches and tags" limited to `main`, and add the secrets that its job uses as environment secrets, never as repo secrets: a repo secret reaches every workflow on every branch, and the agent account may push workflows.
5. **Rulesets.** Settings → Rules → Rulesets → New ruleset → Import a ruleset, once for `ruleset.json` and once for `ruleset-branches.json` in this folder. Add the required checks only after their workflows ran on `main` without errors, or no pull request can pass them, also not the one that fixes a workflow. Check that each required check has GitHub Actions as its source; GitHub offers that source only after the check has run once.
6. **Cloud sessions.** Create a token (classic) of the agent account with the scopes `repo` and `workflow` and an expiry date. Sign in `gh` with it in a separate config directory, start Claude Code there, and run `/web-setup`. It must print `Connected as <agent account>`.
   ```bash
   GH_CONFIG_DIR=~/.config/gh-agent gh auth login --with-token
   GH_CONFIG_DIR=~/.config/gh-agent claude    # then: /web-setup
   ```
   All cloud sessions of the engineer then act as the agent account, also in other projects, so the agent account needs access to each of those repos. The Claude GitHub App does not serve here: it connects only repos that the connected account owns, so Claude Projects threads and auto-fix are not available. Repeat this step before the token expires.

   Two cloud environments: one for design sessions, with open network access and the environment variable `FACTORY_SESSION=design`, and one for all other sessions, whose network access allows only the hosts of `.factory/allowed-hosts`.
7. **Local agent sessions.** Before the first one: merge `adapters/claude-code/user-settings.json` into the Claude Code user settings, and keep its list of hosts equal to `.factory/allowed-hosts`. Give git and `gh` a token (classic) of the agent account with the scopes `repo` and `workflow`; a fine-grained token cannot be used on a repo where its account is only a collaborator. From an agent session, `gh api user` must return the agent account. A local design session starts with `FACTORY_SESSION=design claude`.
8. **Python.** `factory validate-skills` needs Python 3.10 or later with PyYAML as `python3`, on the engineer's computer and in CI.

To check: from a cloud session and from a local agent session, `gh api user` returns the agent account. A pull request that changes a design note waits for the engineer's review as code owner, and a pull request that changes only code merges once its checks pass.

A repo that still uses the environment of [process decision 0001](../decisions/process/0001-approval-through-a-protected-environment.md) switches after step 5: delete `.github/workflows/approval.yml` and the environment `approval`.
