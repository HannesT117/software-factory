# Setup

These steps need a person with admin rights on the repo, because the agent token has no Administration or Workflows rights. Do them once per repo, in this order.

1. **Repo settings, before the first merge.** Settings → General → Pull Requests:
   - "Allow squash merging": on, with the default commit message "Pull request title and description".
   - "Allow merge commits" and "Allow rebase merging": off.
   - "Automatically delete head branches": on.
2. **Agent token.** Create a fine-grained token for the repos that agents work in: Contents and Pull requests write, nothing else.
3. **Environment.** Settings → Environments → New environment `approval`. Required reviewers: the people who approve pull requests. "Prevent self-review": off, because the agent token belongs to the approver's account. "Allow administrators to bypass configured protection rules": off.
4. **Approval workflow.** Copy `approval.yml` to `.github/workflows/approval.yml`, open a pull request, and merge it. The workflow must be on `main` before the ruleset requires its check.
5. **Ruleset.** Settings → Rules → Rulesets → Import a ruleset → `ruleset.json`. Check that the required check `approval` has GitHub Actions as its source.
6. **Local agent sessions.** Before the first one: in the Claude Code user settings, deny reading `~/.ssh`, `~/.config/gh` and the SSH agent socket, set `allowUnsandboxedCommands: false`, and give git the agent token. Increment 02 adds a template for this.

To check: a pull request shows the job `approval` as waiting, and the merge button stays blocked until the engineer approves.
