---
date: 2026-10-27
status: superseded by ADR-0002
---

# Approvals through a protected GitHub environment

## Context and Problem Statement

Agents open pull requests, and a person approves each one before it reaches `main`. Agents act under the GitHub account of the person who approves. The approval must be a step that only a person can take, also when an agent works with that person's GitHub credentials. How does a pull request get an approval that agents cannot give?

## Decision Drivers

* Agents open pull requests under the approver's account. GitHub does not let the author of a pull request approve it.
* A fine-grained token cannot be used to contribute to a repo where its account is only a collaborator. Each token covers one owner.
* Cloud sessions act with the approver's own GitHub credentials through the Claude Code GitHub proxy, and they can push workflow files.
* A required reviewer can approve a waiting job through the API. A personal access token (classic) or an OAuth token needs the `repo` scope for that.
* A required check matches a name and an app. Every workflow reports as the app GitHub Actions, and GitHub counts a skipped job as a passed required check.
* The repos are public, on GitHub Free.
* Verification code that the factory writes itself must be read line by line before it can be trusted.

## Considered Options

1. GitHub pull request reviews with a required approval
2. A separate GitHub account for agents, with required reviews
3. Approval commits signed with a hardware-backed SSH key, checked by `factory`
4. Approval files with a detached SSH signature, checked by `factory`
5. A protected GitHub environment with required reviewers, which a required check waits for

## Decision Outcome

Chosen option 5, because it needs no verification code of its own, works on GitHub Free for public repos, and agents can approve neither with the agent token nor through the cloud proxy.

* **Workflow:** `.github/workflows/approval.yml` runs on `pull_request_target`, so GitHub takes it from `main`. Its job `approval` names the environment `approval` and waits.
* **Environment:** `approval` has required reviewers. "Prevent self-review" is off, because the agent token belongs to the approver's account. Admin bypass is off.
* **Ruleset:** the check `approval` from GitHub Actions is required for `main`. The bypass list is empty.
* **Each new push** to a pull request starts a new run, which needs a new approval.

### Consequences

* Good, because the factory has no code for signatures, keys or trust roots.
* Good, because the approval is one click in GitHub, also in the mobile app.
* Good, because a second person can join as a required reviewer without a change to the factory.
* Bad, because an agent that gets the approver's own GitHub token, for example a `gh` token sent to cloud sessions with `/web-setup`, can approve, and GitHub cannot tell it from the approver.
* Bad, because the gate holds only while the real job runs on every event that brings a pull request to `main`. A fake job with the same name passes otherwise. The workflow therefore also reacts to `edited`, and its job has no `if:`.
* Bad, because the approval is not recorded in git. It stays in GitHub's deployment history.
* Bad, because environment reviewers on private repos need a paid plan.

### Confirmation

* Scenario 02.10: a pull request cannot be merged without the approval, also not with a fake check of the same name.
* Scenario 02.11: the gate runs from `main`.
* Scenario 02.12: neither a local agent session nor a cloud session can approve the job or change the ruleset.

## Pros and Cons of the Options

### GitHub pull request reviews

* Good, because it is the usual review flow of teams on GitHub.
* Bad, because agents open pull requests under the approver's account, and the author of a pull request cannot approve it.

### Separate GitHub account for agents

* Good, because GitHub can tell agents from the approver, as with the Copilot cloud agent.
* Bad, because the account would need a classic token or an organization, since a fine-grained token cannot be used on a repo where its account is a collaborator.
* Bad, because cloud sessions still act with the approver's credentials.

### Signed approval commits

* Good, because an approval needs a touch on a hardware key, also when an agent has the approver's GitHub token.
* Bad, because `factory` must verify signatures, keys and a trust root, and that code must be read line by line.
* Bad, because a squash merge removes the signature from `main`.

### Signed approval files

* Good, because the signature survives the squash merge.
* Bad, because `factory` still needs verification code and a trust root for keys.

### Protected GitHub environment

* Good, because GitHub enforces it, and the factory only configures it.
* Bad, because it depends on GitHub's environment feature and on the ruleset.

## More Information

* **Revisit when** a second person approves, agents get their own GitHub identity, or the repos become private.
* **Sources:** GitHub docs on pull request reviews, environments, fine-grained token permissions and personal access tokens; GitHub Copilot cloud agent, risks and mitigations; Claude Code docs on the GitHub proxy. Details are in the [research notes](../../research/2026-10-factory-basis.md).
