# 0003: Public repos, a GitHub ruleset is the merge gate

- Status: proposed
- Date: 2026-10-26

## Context

Something must stop a merge when checks fail.

- Git hooks run on the client. `--no-verify` skips them, and a fresh clone has none until `core.hooksPath` is set.
- Agent hooks run only in one tool.
- GitHub rulesets with required checks are free for public repos. Private repos need GitHub Pro or higher.
- You chose public repos without health data over paying for Pro.

## Decision

1. The factory repo and `health` are public.
2. Each repo gets the ruleset from `core/ci/github/ruleset.json` on `main`:
   - changes only through pull requests
   - the factory checks are required and must be green
   - the bypass list is empty, so the rule applies to admins too
3. Secret scanning with push protection is on.
4. `.factory/config.json` lists private paths. A check fails if any of them exists in the repo.
5. Personal health details move from `health` docs into a private notes repo. Holdout scenarios live in a private repo too.
6. Git hooks and agent hooks stay, as early feedback only.

## Consequences

Good:

- The gate runs on the server. No local trick bypasses it.
- No extra cost.

Bad:

- All code, plans and research in these repos are visible to everyone.
- The private-paths check finds only listed paths. It does not find health data written into other files. You stay responsible for that.

## Re-check when

You want private repos again (then GitHub Pro, about 4 USD a month), or GitHub changes ruleset availability.

## Verification

Once in phase 2: a pull request with a red check cannot be merged, also not by you.

## Sources

GitHub docs on ruleset availability; git documentation on `--no-verify`; adversarial review of the plan. Details in [research](../research/2026-10-factory-basis.md).
