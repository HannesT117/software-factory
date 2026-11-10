# Review of increment 06: gate and status

Reviewed commit: 16e4808119f0ab47fd9adabffe583e1f021d587b
Pull request: https://github.com/HannesT117/software-factory/pull/69
Result: 1 blocking, 1 notes

## Findings

- **blocking** Protected paths, `test/acceptance/06-gate-and-status/scenario-06.11.txtar`: commit 16e4808 changes this test file after the first code commit c1f93d8. The new test does check its scenario: with `internal/crossreview/crossreview.go:183` set to send the instructions of the code review, it fails at line 42, and with the words "including each case where something must not happen" removed from the instructions of the tests review, it fails at line 43. No check reports the change on this pull request, because `check-pr` stops at the mixed path groups before it looks at the order of tests and code. Doubtful: the coordinator reports that the engineer allowed this test change for this pull request, but nothing in the repo or in the pull request records it; after two rounds the engineer decides.
- **note** Special cases, `.github/workflows/cross-review.yml:52-54` after the patch in the first comment of the pull request: the workflow posts a comment only when the output has a line `cross-review: model`. The Interface says it posts its output "unless the review does not apply"; so a pull request whose code waits for a review result, or a run without a key or configuration, gets no comment, only a failed job with the reason in its log. The engineer then sees on the pull request no reason why `cross-review` is red. Doubtful: the change answers a deviation of the previous round and may be what the engineer wants; then the Interface of the note should say so.

## Checks run

The measure is the note of 06 at the head of #68 (bcc08bf), which changes the sketch and the owned paths against the note on this branch. Each `factory` command ran with a binary built from 16e4808, unless the line says otherwise. `simbase` is `origin/main` with the design notes of #67 (5af365e) and #68 (12c6316) merged; it is not an ancestor of the branch, so its merge base with the branch is 5af365e, and the note of 06 counts as changed against it.

- `go vet ./...`: no output
- `go test -count=1 -v ./test/acceptance/`: `ok  	github.com/HannesT117/software-factory/test/acceptance`, 47 scripts pass
- scenario 06.11 alone, in a scratch worktree of 16e4808: with `crossreview.go:183` set to `instructions = codeInstructions`, `FAIL: 06-gate-and-status/scenario-06.11.txtar:42: no match for (?i)checks what its scenario says`; with "must not happen" removed from the instructions of the tests review, `FAIL: ...:43: no match for (?i)must not happen`
- `factory check-design docs/increments/06-gate-and-status` in a worktree of bcc08bf: `check-design: ok 15 checked`
- `factory check-design --base simbase`: `check-design: ok 15 checked`; `--base origin/main`: `check-design: ok 23 checked`
- `factory check-tests --base simbase`: `check-tests: ok 47 checked`; `--base origin/main`: `check-tests: 16 errors`, because the notes of #67 and #68 are not on `origin/main`
- `factory check-pr --base simbase`: `check-pr: 46 errors`, `mixed-pr` on every changed path, with `.github/CODEOWNERS`, as the engineer allowed, and with the note of 06 and `docs/research/2026-11-factory-review.md` from the merge base above; `--base origin/main`: `check-pr: 47 errors`
- `factory check-review --base simbase` before these results: `check-review: 1 errors`, `missing-review` for 16e4808119f0ab47fd9adabffe583e1f021d587b; with them committed: `check-review: reviewed 16e4808119f0ab47fd9adabffe583e1f021d587b in evidence/05-tests-per-increment/review.md, round 2`, the same for this file, `check-review: error evidence/06-gate-and-status/review.md: blocking 1 in round 2; after 2 rounds the engineer decides` and `check-review: 1 errors`
- `factory status` before these results: `status: review-open 16e4808119f0ab47fd9adabffe583e1f021d587b`, `status: ok 5 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before these results: blocks and names 16e4808119f0ab47fd9adabffe583e1f021d587b
- the workflow patch of the first comment of the pull request (SHA-256 `a68f361b68ad6da9a3ad68a154278b713b4eb1df2048ba6d3a7abb9332cc17b3`) applies on 16e4808
- `git merge-tree --write-tree HEAD origin/main` and `HEAD bcc08bf`: no conflict

## Demonstrated scenarios

- **scenario 06.13**: `evidence/06-gate-and-status/acceptance-run.md`. `go test -count=1 -v ./test/acceptance/` at 16e4808 prints `acceptance: coverage 85.4% of statements` and three functions that no script reaches: `internal/report/report.go:44 Warning`, `internal/setup/setup.go:62 Error` and `internal/skills/skills.go:140 lastLine`; all 47 scripts pass.
