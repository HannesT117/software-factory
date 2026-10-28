# Setup

These steps need a person with admin rights on the repo, because the agent token has no Administration or Workflows rights. Do them once per repo, in this order.

1. **Repo settings, before the first merge.** Settings → General → Pull Requests:
   - "Allow squash merging": on, with the default commit message "Pull request title and description".
   - "Allow merge commits" and "Allow rebase merging": off.
   - "Automatically delete head branches": on.
2. **Agent token.** Create a fine-grained token for the repos that agents work in: Contents and Pull requests write, nothing else.
3. **Environment.** Settings → Environments → New environment `approval`. Required reviewers: the people who approve pull requests. "Prevent self-review": off, because the agent token belongs to the approver's account. "Allow administrators to bypass configured protection rules": off.
4. **Approval workflow.** Copy `.github/workflows/approval.yml` from the factory repo to the same path in the repo, open a pull request, and merge it. In the factory repo itself, the file is already there. The workflow must be on `main` before the ruleset requires its check.
5. **Ruleset.** Settings → Rules → Rulesets → Import a ruleset → `ruleset.json`. Check that the required check `approval` has GitHub Actions as its source. GitHub offers that source only after the job has run once: if it is missing, approve the job on one pull request first. Add the required check only after the workflow on `main` runs without errors, or no pull request can pass it, also not the one that fixes the workflow.
6. **Local agent sessions.** Before the first one: in the Claude Code user settings, deny reading `~/.ssh`, `~/.config/gh` and the SSH agent socket, both in the sandbox and as permission rules for the Read tool, because the sandbox covers only Bash. Set `allowUnsandboxedCommands: false`, and give git the agent token. `gh` and git can keep tokens in the system keychain: from an agent session, `gh auth token` and `git credential fill` must both fail. Increment 02 adds a template for this.
7. **Cloud sessions.** Connect GitHub through the Claude GitHub App. Do not use `/web-setup`: it sends the local `gh` token to cloud sessions, and a token with the `repo` scope can approve the job `approval`.

To check: a pull request shows the job `approval` as waiting, and the merge button stays blocked until the engineer approves.
