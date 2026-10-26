# Setup

The engineer does these steps once per repo, in this order. Agents cannot do them, because the agent token has no Administration or Workflows rights.

1. **Agent token.** Create a fine-grained token for the factory repo and `health`: Contents and Pull requests write, nothing else.
2. **Environment.** Settings → Environments → New environment `approval`. Required reviewers: the engineer. "Prevent self-review": off, because agents act under the engineer's account. "Allow administrators to bypass configured protection rules": off.
3. **Approval workflow.** Copy `approval.yml` to `.github/workflows/approval.yml`, open a pull request, and merge it. The workflow must be on `main` before the ruleset requires its check.
4. **Ruleset.** Settings → Rules → Rulesets → Import a ruleset → `ruleset.json`. Check that the required check `approval` has GitHub Actions as its source.
5. **Local agent sessions.** Before the first one: in the Claude Code user settings, deny reading `~/.ssh`, `~/.config/gh` and the SSH agent socket, set `allowUnsandboxedCommands: false`, and give git the agent token. Increment 02 adds a template for this.

To check: a pull request shows the job `approval` as waiting, and the merge button stays blocked until the engineer approves.
