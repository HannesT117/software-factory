# Retro of increment 03: setup script

## Metrics

| Metric | Value |
|---|---|
| Scenarios | 8 automated, 1 manual (03.9, open) |
| Test writer's questions | 0 |
| Tests changed after review | 0 |
| Code commits | 3; 663 lines in `cmd/` and `internal/`, with the two moved ruleset files |
| `go test ./...` runs until green | 1 |
| Review rounds | 2; round 1: 1 finding (ruleset comparison); round 2: 1 finding outside the diff (the required check `checks` fails on the merge ref) |
| Engineer's interventions during the code | 0; the orchestrating session settled the one question (how a ruleset is compared) |
| Agent | Claude Code 2.1.317, cloud session; a test-writer session, an implementing session, and the subagent `increment-reviewer` |

## Corrections

- The design note says that a ruleset on GitHub "equals" the embedded file. The code compares by containment, because GitHub returns default keys that the file lacks; the reviewer reported this as its finding. The orchestrating session kept containment, and design pull request #56 changes the note to "holds every key and value that the file sets".
- `docs/setup/README.md` told the engineer to build `factory` and then run it in another folder, where it is not on `PATH`. Fixed with `go install`.
- The implementing session asked the reviewer to name the branch head as the reviewed commit. Round 1 named a docs commit, so neither the stop hook nor the CI job `cross-review` accepted it; round 2 names the latest commit with code as well.
- The implementing session deleted the reviewer's uncommitted `review.md` in the shared worktree, and restored it from the reviewer's commit.

## Points for the factory

1. **A reviewer subagent shares the implementing session's worktree.** While it wrote `evidence/`, the stop hook asked the implementing session to commit the untracked file, which would have committed an unfinished review. Fix: the reviewer runs in a worktree of its own, or the hook ignores `evidence/` while a review runs.
2. **A design note can fix the shape of an external answer that nobody has seen.** "Equals" assumed that GitHub returns a ruleset as it was sent. Fix: the skill `designing-increments` asks the note to mark each assumption about an external API as unverified and to name the manual scenario that settles it.
3. **The retro asks for the agent version, and the cloud harness forbids model names in pushed files.** This retro names the tool version only. Fix: `AGENTS.md` says which version counts.
4. **The acceptance stand-in for `gh` answers like the documented API, not like the real one.** It returns a ruleset exactly as it was stored, so no test can show the default keys of point 2. Fix: record one real answer per endpoint in scenario 03.9 and let the stand-in add the same defaults.

5. **`checks` judges the merge ref, not the branch.** The workflow checks out GitHub's merge commit of the pull request into `main`; `factory check-pr` then sees one commit that changes tests and code together and reports each test file. Every pull request with tests and then code fails this required check. Fix, in its own pull request: check out the head of the pull request in `.github/workflows/checks.yml`, or let `check-pr` read a merge ref.
6. **"The latest commit with code" is easy to get wrong.** The reviewer's prompt and the skill `reviewing-increments` differ on which commit a review names. Fix: `factory hook stop` prints the hash that the review must name.

## Summary

`factory setup` applies the admin settings of a repo in six steps, with a dry run and with repair on a second run. The tests of 03.1 to 03.8 passed on the first full run. Review 1 found one deviation from the note: rulesets are compared by containment, not equality; #56 changes the note to match. The required check `checks` fails on every pull request with tests and then code; that fix needs its own pull request. The manual scenario 03.9 on a real repo is open; it also shows which default keys GitHub adds to a ruleset.
