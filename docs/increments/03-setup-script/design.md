# Increment 03: setup script

## Goal

- The engineer applies the admin part of the setup of a repo with one command, `factory setup`, with their own `gh`: merge settings, the invitation of the agent account, the environment `cross-review` with its secret, and both rulesets with their required checks.
- A second run sets every step to the target state again and repairs settings that changed by hand. A dry run prints what a run would change and changes nothing.
- `docs/setup/README.md` lists only the steps that stay manual and names the command for the rest.

## Out of scope

- Copying `.github/CODEOWNERS`, the workflows and `.factory/` into a project: `factory install`, increment 05.
- Steps without an API or outside the repo: creating the agent account and its token, `/web-setup`, the two cloud environments, the local Claude Code settings, Python. They stay manual steps in `docs/setup/README.md`.
- Settings of an organization, such as teams or organization rulesets.
- Reading back the value of a secret. GitHub does not return it; the command only knows whether the secret exists.

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Applying repo settings from a file | [Probot settings app](https://github.com/repository-settings/app) | not used | Needs a GitHub App with admin rights installed on the repo; it has no section for rulesets. |
| Applying repo settings from a file | [safe-settings](https://github.com/github/safe-settings) | not used | A self-hosted GitHub App with an admin repo, made for organizations; the engineer's repos are personal. |
| Applying repo settings from a file | [Terraform GitHub provider](https://registry.terraform.io/providers/integrations/github/latest/docs) (`github_repository_ruleset`, `github_repository_environment`, `github_repository_environment_deployment_policy`) | not used | Another runtime and a state file per repo for a setup that the engineer runs rarely. |
| Calling GitHub | `gh api`; a Go HTTP client with a token | `gh api` | Uses the engineer's existing `gh` login, so the command handles no token; the runtime stays the standard library. |
| Importing a ruleset | `gh ruleset` (only `list`, `view`, `check`); `gh` extensions; [rulesets endpoint](https://docs.github.com/en/rest/repos/rules) | rulesets endpoint, through `gh api` | No subcommand or extension creates a ruleset. |
| Merge settings | [Update a repository](https://docs.github.com/en/rest/repos/repos#update-a-repository); [`gh repo edit`](https://cli.github.com/manual/gh_repo_edit) | Update a repository, through `gh api` | Both can set every field; one form of call for all steps keeps the stand-in `gh` of the tests small. |
| Agent account | [Add a repository collaborator](https://docs.github.com/en/rest/collaborators/collaborators#add-a-repository-collaborator); [Get repository permissions for a user](https://docs.github.com/en/rest/collaborators/collaborators#get-repository-permissions-for-a-user) | used, through `gh api` | Its `permission` field applies only to repos of an organization; on a personal repo a collaborator has write access, never admin. |
| Environment `cross-review` | [Create or update an environment](https://docs.github.com/en/rest/deployments/environments); [deployment branch policies](https://docs.github.com/en/rest/deployments/branch-policies) | used, through `gh api` | One `PUT` creates or updates the environment; `custom_branch_policies` with one policy for `main` limits it to `main`. |
| Environment secret | [Create or update an environment secret](https://docs.github.com/en/rest/actions/secrets#create-or-update-an-environment-secret); [`gh secret set --env`](https://cli.github.com/manual/gh_secret_set) | `gh secret set` | The endpoint needs the value encrypted with a libsodium sealed box, which the standard library lacks; `gh secret set` encrypts it and reads the value from standard input. |
| Required checks | [Available rules for rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets); rulesets endpoint | in the ruleset `main`, imported at once | The engineer's decision. The docs say that the app of a check "must have recently submitted a check run"; whether the endpoint accepts `integration_id` before that is verified in the manual scenario. |

## Interface

`factory setup --agent <login> [--repo <owner>/<name>] [--dry-run]`

Output follows the rules of the Interface of increment 02: every line starts with `setup:`, error lines are `setup: error <where>: <code> <detail>`, all output goes to stdout. Exit codes: 0 ok, 1 errors, 2 usage or environment.

**Calls.** The command calls GitHub only through `gh`, which must be on `PATH` and logged in: `gh api` for every step except the secret, and `gh secret set OPPER_API_KEY --env cross-review --repo <owner>/<name>` with the value on standard input. Without `--repo`, the repo is the one of the current directory, from `gh repo view --json nameWithOwner --jq .nameWithOwner`.

**Rulesets.** The binary carries `internal/setup/rulesets/main.json` and `internal/setup/rulesets/branches.json`, embedded at build time; the command reads no ruleset file at run time. Their content, in the format of the [rulesets endpoint](https://docs.github.com/en/rest/repos/rules), is part of this interface. They replace `docs/setup/ruleset.json` and `docs/setup/ruleset-branches.json`, which are removed.

**Steps**, in this order. For each, the command reads the current state, compares it with the target, and writes only when they differ.

| Step | Target |
|---|---|
| `merge-settings` | `allow_squash_merge` true, `allow_merge_commit` false, `allow_rebase_merge` false, `allow_auto_merge` true, `delete_branch_on_merge` true, `squash_merge_commit_title` `PR_TITLE`, `squash_merge_commit_message` `PR_BODY`. |
| `agent-account` | `<login>` is a collaborator or has a pending invitation. If missing, `PUT repos/<repo>/collaborators/<login>` with `permission` `push`. |
| `environment` | The environment `cross-review` exists with `deployment_branch_policy` `{"protected_branches": false, "custom_branch_policies": true}` and exactly one branch policy, `main`. Other branch policies are deleted. |
| `secret` | The environment `cross-review` has the secret `OPPER_API_KEY`. If missing, the command reads its value as one line from standard input. If present, the command does not read standard input and leaves it as it is. |
| `ruleset main` | A ruleset named `main` in the repo equals `internal/setup/rulesets/main.json` in `name`, `target`, `enforcement`, `bypass_actors`, `conditions` and `rules`. If missing, `POST repos/<repo>/rulesets`; if different, `PUT repos/<repo>/rulesets/<id>`. Rulesets with other names stay as they are. |
| `ruleset branches` | The same for `internal/setup/rulesets/branches.json` and the ruleset named `branches`. |

**Output lines**, one per step: `setup: unchanged <step>`, `setup: changed <step>`, or with `--dry-run` `setup: would change <step>`. A changed or would-change line of `merge-settings` and of a ruleset names the fields that differ, as `<field> <current> -> <target>` separated by `, `. The last line is `setup: ok 6 checked` or `setup: <n> errors`.

A writing call is a `gh api` call with a method other than `GET`, or `gh secret set`.

**Secret value.** Before it reads, the command prints `setup: reading OPPER_API_KEY from standard input`. It reads up to the first newline or the end of input and strips surrounding white space. An empty result is `no-secret`.

**Dry run.** With `--dry-run`, the command makes no writing call and reads no standard input. A missing secret is `setup: would change secret`.

**Errors.** Before the first step, the command checks `gh`, the repo, the admin right, and the rights of `<login>` (`permission` of `GET repos/<repo>/collaborators/<login>/permission`; a 404 means no rights). An error there changes nothing. An error in a step stops the run; the steps after it do not run.

| Code | Exit | When |
|---|---|---|
| `usage` | 2 | `--agent` is missing, or an unknown flag. |
| `no-gh` | 2 | `gh` is not on `PATH`, or is not logged in. |
| `no-repo` | 2 | No `--repo` and the current directory has no GitHub repo, or the repo does not exist. |
| `not-admin` | 2 | The logged-in account has no admin right on the repo (`permissions.admin` of `GET repos/<repo>` is false). |
| `agent-is-admin` | 1 | `<login>` has the permission `admin` on the repo. The command does not change it. |
| `no-secret` | 1 | In the step `secret`: the secret is missing and the value read is empty. |
| `api-failed` | 1 | A call through `gh` exits with a status other than 0. `<where>` is the step; the detail is the method and path of the call, or `gh secret set`, and the first line of `gh`'s error output. |

## Scenarios

- **scenario 03.1** (automated) **The first run sets up a new repo.** On a repo with none of the settings, with a stand-in for `gh` that records each call, `factory setup` with a secret value on standard input reports each of its six steps as changed. Afterwards the stand-in holds: squash merging with the title and description of the pull request as message, no merge commits and no rebase merging, auto-merge, branch deletion, an invitation for the agent account, the environment `cross-review` that only `main` may use, its secret `OPPER_API_KEY`, and the rulesets `main` and `branches` with the content that the binary carries.
- **scenario 03.2** (automated) **A second run changes nothing that is in place.** On a repo where every step already has its target state, `factory setup` reports each step as unchanged and writes nothing. It does not read standard input.
- **scenario 03.3** (automated) **A second run repairs what changed by hand.** On a repo in its target state, where someone then turned merge commits back on and set the ruleset `main` to disabled, `factory setup` reports these two steps as changed, names each field with its old and new value, and sets them back. The other steps stay unchanged. A third ruleset with a name of its own, which the setup does not know, stays as it is.
- **scenario 03.4** (automated) **A dry run shows the changes and makes none.** On the repo of scenario 03.1, `factory setup --dry-run` reports each step that a run would change, with the differing fields, and makes no writing call. It does not read standard input.
- **scenario 03.5** (automated) **Without admin rights, nothing changes.** When the account logged in to `gh` is not an admin of the repo, `factory setup` refuses before the first step, and reports this as a problem of the environment. The same holds when `gh` is missing.
- **scenario 03.6** (automated) **An agent account with admin rights is reported.** When the agent account already has admin rights on the repo, `factory setup` reports it as an error before the first step, changes nothing, and leaves the account's rights as they are.
- **scenario 03.7** (automated) **The secret is asked for only when it is missing.** When the secret is missing, `factory setup` says that it reads the value, reads it from standard input, and passes it to `gh` only on standard input, never as an argument. When standard input is empty or holds only white space, it reports the missing secret as an error and stops; the steps before it stay done, and the rulesets are not imported.
- **scenario 03.8** (automated) **A failed call stops the run, and the next run continues.** On a new repo, GitHub refuses to create the ruleset `main`. `factory setup` reports an error that names the step and the refused call, and does not create the ruleset `branches`. A second run, after GitHub accepts the call, reports the merge settings, the agent account, the environment and the secret as unchanged, does not read standard input, and reports both rulesets as changed.
- **scenario 03.9** (manual) **The setup works on a real repo.** The engineer creates an empty public personal repo without workflows, and runs `factory setup --dry-run` and then `factory setup`, logged in to `gh` as themself and naming the agent account as the one to invite. The settings pages of the repo then show the merge settings, the invitation, the environment with its secret, and both rulesets; the ruleset `main` lists its four required checks with GitHub Actions as their source, although none of them has run. A second run reports every step as unchanged.

## Evidence

- **Approved scenarios:** 03.1 to 03.8.
- **Manual:** 03.9.
- **Not evidence:** unit tests from the implementing agent.

## Rollout

- **Off switch:** none needed; the engineer runs the command by hand. A setting it made is changed back in the GitHub settings, and a ruleset is switched off there.
- **Failure signal:** after a run, a dry run reports a change, or GitHub refuses a call that the docs allow.
