# Agent identity on GitHub

Question: under which GitHub identity do agents push and open pull requests, when the machine account of [process decision 0002](../decisions/process/0002-approval-through-code-owners-and-an-agent-account.md) can be suspended. Read on 2026-11-08.

## Facts

- **Terms of Service, B.3:** "We do permit machine accounts." A machine account is set up by a human who "is responsible for its actions" and is "used exclusively for performing automated tasks". "You may maintain no more than one free machine account in addition to your free Personal Account." ([GitHub Terms of Service](https://docs.github.com/en/site-policy/github-terms/github-terms-of-service))
- **GitHub Apps:** "GitHub Apps are not tied to a user account and do not consume a seat." They act "independently of a user", with fine-grained permissions, per-installation repository selection and short-lived tokens. ([Deciding when to build a GitHub App](https://docs.github.com/en/apps/creating-github-apps/about-creating-github-apps/deciding-when-to-build-a-github-app))
- **Installation token:** a JWT signed with the app's private key, then `POST /app/installations/<id>/access_tokens`. The token can be limited to repositories and permissions, and expires after 1 hour. ([Generating an installation access token](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/generating-an-installation-access-token-for-a-github-app)) For a terminal, the `gh` extension [gh-token](https://github.com/Link-/gh-token) makes such a token from the key file; `gh` reads it from `GH_TOKEN`, and `gh auth setup-git` lets git use it.
- **Claude Code on the web:** GitHub is connected through the Claude GitHub App or `/web-setup`, which sends the token that `gh auth token` prints. A proxy adds the connected credentials to each GitHub request; network secrets do not apply to GitHub. Auto-fix replies "are posted using your GitHub account". ([Claude Code on the web](https://code.claude.com/docs/en/claude-code-on-the-web), [Cloud environments](https://code.claude.com/docs/en/cloud-environments))
- **claude-code-action:** runs Claude Code in GitHub Actions. With a custom GitHub App, `actions/create-github-app-token` makes the token from the secrets `APP_ID` and `APP_PRIVATE_KEY`, and the action gets it as `github_token`. The app needs Contents, Issues and Pull requests read and write. With a `prompt` input, it runs on any event without a mention. It authenticates to Claude with an API key or with `CLAUDE_CODE_OAUTH_TOKEN` from `claude setup-token` (Pro, Max, Team, Enterprise). ([GitHub Actions](https://code.claude.com/docs/en/github-actions), [setup.md](https://github.com/anthropics/claude-code-action/blob/main/docs/setup.md))

## Options

| Option | Where agents run | Identity | Consequences |
|---|---|---|---|
| A. Machine account, as today | cloud sessions | `HannesT117-ai` | Allowed by the Terms, but one suspension stops all agents; the reason is not public. |
| B. GitHub App, agents in GitHub Actions | `claude-code-action`, started by a workflow | `<app>[bot]` | No account to suspend. Questions to the engineer go through pull request and issue comments; no chat, no messages between sessions. A job ends after at most 6 hours on hosted runners. |
| C. GitHub App, agents on the engineer's computer | local Claude Code, steered from the phone through Remote Control | `<app>[bot]` | A token helper renews the 1-hour token. The computer must be on. The sandbox and the user settings of the Claude Code adapter keep the session away from the engineer's own credentials. |
| D. GitHub App in cloud sessions | cloud sessions | not possible | The proxy uses the connected user account for GitHub; an app token in `GH_TOKEN` would reach only `gh`, expire after an hour, and sit in plain text in the environment. |
| E. The engineer's own account | any | the engineer | Rejected: agents could not be told apart from the engineer, and nobody could approve the engineer's own pull requests as code owner. |

## Not verified

- Whether commits through `claude-code-action` with a custom app show the app's `[bot]` user or need `bot_id` and `bot_name`.
- Whether an approval by a GitHub App counts toward required approvals of a ruleset. An app with Pull requests write can submit reviews; it cannot be a code owner, so it cannot satisfy code owner review.
- Which permission an app needs to push changes to `.github/workflows/`: by GitHub's permission model, Workflows write.
