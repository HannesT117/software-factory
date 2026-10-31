---
date: 2026-10-31
---

# Approvals through code owners and an agent account

## Context and Problem Statement

Agents open pull requests, and the engineer approves only some of them: design notes, the domain model, decision records, the rules for agents, and CI and factory config. All other pull requests merge when their checks pass. The approval must be a step that only the engineer can take. How does a pull request get such an approval, with as little verification code of the factory's own as possible?

## Decision Drivers

* GitHub allows one free machine account next to a personal account, used only for automated tasks.
* The GitHub account that is connected to Claude decides as whom cloud sessions act.
* Code owners are available for public repos on GitHub Free. A ruleset can require a code owner's approval for the files that the code owner owns, and only for those.
* The author of a pull request cannot approve it, and an approval by someone who is not a code owner does not count as code owner review.
* A workflow in a pull request reports as the app GitHub Actions, like the real checks.
* A fine-grained token cannot be used on a repo where its account is only a collaborator. A token (classic) can.
* The engineer approves the scenarios of an increment, not its test code.
* Verification code that the factory writes itself must be read line by line before it can be trusted.

## Considered Options

1. A protected GitHub environment with required reviewers, which a required check waits for
2. Direct pushes to `main`, with the approval asked by an agent hook
3. A protected branch for design notes only, and direct pushes to `main`
4. Approval commits signed with a hardware-bound key, checked outside the repo
5. A machine account for agents, with code owner review

## Decision Outcome

Chosen option 5, because GitHub enforces the approval with the mechanism that teams use, the engineer approves only owned paths, and no verification code of the factory's own guards the approval.

* **Accounts:** agents act as a machine account with write access and no admin rights. Cloud sessions act as the account connected to Claude; local sessions use a token (classic) of the machine account. The engineer approves with a personal account.
* **Code owners:** `.github/CODEOWNERS` makes the engineer the owner of design notes, the domain model, decision records, `AGENTS.md`, `docs/principles.md`, `.github/` and `.factory/`.
* **Ruleset on `main`:** pull requests only, squash merges only, a code owner's approval when an owned path changes, approvals dismissed on a new push, the required checks from GitHub Actions, no force push or deletion, and an empty bypass list.
* **Merge:** a pull request that changes no owned path merges as soon as its checks pass. Agents may merge or turn on auto-merge; GitHub refuses the merge until the rules hold.

### Consequences

* Good, because GitHub tells agents and the engineer apart, also in the history.
* Good, because the approval is one click in GitHub, also in the mobile app, and only for owned paths.
* Good, because the environment, the workflow on `pull_request_target` and a list of approval-free paths are no longer needed.
* Good, because a second person can join as a code owner without a change to the factory.
* Bad, because all agents share one account, so GitHub cannot tell the test writer from the implementing agent.
* Bad, because all of the engineer's Claude sessions act as the machine account, also in other projects, and that account needs access to each of their repos.
* Bad, because local sessions need a token (classic) with the `repo` scope for the machine account.
* Bad, because paths that nobody owns, test files among them, merge without a person; there the guards are checks.
* Bad, because a second account, with its own email and two-factor authentication, needs care.

### Confirmation

* Scenario 02.10: a pull request that changes a design note cannot be merged without the engineer's approval; a pull request that changes only code merges without it.
* Scenario 02.11: a pull request cannot fake a check, because `.github/` is owned.
* Scenario 02.12: agent sessions act as the machine account.

## Pros and Cons of the Options

### Protected environment

* Good, because GitHub enforces it, and it works with one account.
* Bad, because it cannot tell paths apart, so the factory needs a list of approval-free paths and its own logic to decide when the engineer is needed.
* Bad, because it stays out of the agents' reach only through a workflow on `pull_request_target`, and only while that job runs on every event that brings a pull request to `main`.

### Direct pushes with a hook

* Good, because the engineer approves inside the session, and there are no pull requests.
* Bad, because an agent can get around a hook through the shell, and nothing on GitHub checks the approval.
* Bad, because checks run only after the push, so `main` can hold failing code.

### Protected branch for design notes

* Good, because GitHub enforces the approval of design notes, and `main` takes direct pushes.
* Bad, because the design notes live apart from the code, and the tools read them from a second branch.
* Bad, because the gate still needs the environment or one account per role.

### Signed approval commits

* Good, because an approval needs a touch on a hardware key, and there are no pull requests.
* Bad, because it only detects an unsigned change after the push, and the checker must run outside the repo, with a trust root out of the agents' reach.
* Bad, because the engineer can approve only at the computer that holds the key.

### Machine account with code owners

* Good, because it uses GitHub's standard review for teams.
* Bad, because it needs a second account and a token (classic) for local sessions.

## More Information

* **Replaces** [process decision 0001](0001-approval-through-a-protected-environment.md).
* **Revisit when** a second person approves, GitHub offers an identity per agent session, or the repos become private.
* **Sources:** GitHub Terms of Service on machine accounts; GitHub docs on code owners, rulesets and the limits of fine-grained tokens; Claude Code docs on the GitHub connection of cloud sessions.
