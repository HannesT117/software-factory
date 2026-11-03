# Increment 03: setup script

## Goal

- (draft)

## Out of scope

- (draft)

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Applying repo settings from a file | [Probot settings app](https://github.com/repository-settings/app) | not used | Needs a GitHub App with admin rights installed on the repo; it has no section for rulesets. |
| Applying repo settings from a file | [safe-settings](https://github.com/github/safe-settings) | not used | A self-hosted GitHub App with an admin repo, made for organizations; the engineer's repos are personal. |
| Applying repo settings from a file | [Terraform GitHub provider](https://registry.terraform.io/providers/integrations/github/latest/docs) (`github_repository_ruleset`, `github_repository_environment`, `github_repository_environment_deployment_policy`) | not used | Another runtime and a state file per repo for a setup that runs once. |
| Importing a ruleset | `gh ruleset` (only `list`, `view`, `check`); `gh` extensions | not used | No subcommand or extension creates a ruleset. `gh api --input <file>` against the [rulesets endpoint](https://docs.github.com/en/rest/repos/rules) does. |
| Merge settings | [Update a repository](https://docs.github.com/en/rest/repos/repos#update-a-repository); [`gh repo edit`](https://cli.github.com/manual/gh_repo_edit) | `gh repo edit` | It sets squash merging with `--squash-merge-commit-message pr-title-description`, auto-merge and branch deletion; turning off merge commits and rebase merging needs `--enable-merge-commit=false` and `--enable-rebase-merge=false`. |
| Agent account | [Add a repository collaborator](https://docs.github.com/en/rest/collaborators/collaborators#add-a-repository-collaborator) | used, through `gh api` | Answers 201 for a new invitation and 204 for an existing collaborator. Its `permission` field applies only to repos of an organization; on a personal repo a collaborator always has write access, never admin. |
| Environment `cross-review` | [Create or update an environment](https://docs.github.com/en/rest/deployments/environments); deployment branch policies | used, through `gh api` | One `PUT` creates or updates the environment; `custom_branch_policies` with a policy for `main` limits it to `main`. |
| Environment secret | [`gh secret set --env`](https://cli.github.com/manual/gh_secret_set) | used | Reads the value from standard input, so the value appears in no argument and no file. |
| Required checks | [Available rules for rulesets](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/available-rules-for-rulesets); rulesets endpoint | open | The docs say that the app of a check "must have recently submitted a check run"; the endpoint takes `integration_id` without saying whether it checks that. Not verified. |
