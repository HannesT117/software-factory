# Review of increment 06: gate and status

Reviewed commit: e30dda011b555807d757e4a6753d5402b888d486
Pull request: https://github.com/HannesT117/software-factory/pull/69
Result: 2 blocking, 1 notes

## Findings

- **blocking** Tests that miss their scenario, `test/acceptance/06-gate-and-status/scenario-06.11.txtar:29-36`: scenario 06.11 sends the tests "with instructions that ask whether each test checks what its scenario says". The test checks the model, the note of 08, its scenario titles and the diff, but nothing of the instructions. With `internal/crossreview/crossreview.go:183` changed to send the instructions of the code review instead, the test still passes (run in a scratch worktree). Before the code it failed only on its line `cross-review: model code-model`, which scenario 06.12 covers. The code does send the instructions of the tests review (`crossreview.go:51`), but no test holds them: a later change can drop the review of the tests, the point of this scenario, with every test green.
- **blocking** Special cases, `internal/review/review.go:55-70` and `internal/review/check.go:38`: `Find` returns the first `evidence/*/review.md`, in path order, that names the hash, and `check-review` reads only that file. A pull request with two increments, as the Rollout of this note plans for 05 and 06, gets one review result per increment for the same latest code commit, and a blocking finding in the second file does not hold the merge. On this branch, with this result committed, `factory check-review --base simbase` reads only `evidence/05-tests-per-increment/review.md` and passes, while this file counts 2 blocking findings; a scratch repo with two such results gives the same. That is a failure signal that the note names: `check-review` passes a result with a blocking finding. Doubtful only because the Interface speaks of "the file" that names the hash, in the singular.
- **note** Special cases, `.github/CODEOWNERS:12-20`: the new lines own `/cmd/`, `/internal/`, `/adapters/`, `/core/` and `/test/`, but not `/.claude/`, through whose symlinks the sessions load the skills (`.claude/skills/*`) and the agents `increment-reviewer` and `acceptance-test-writer` (`.claude/agents/*`), and where Claude Code reads project settings. A pull request can point `.claude/skills/reviewing-increments` at an unowned copy, or add `.claude/settings.json`, and so change what judges the next pull request without the engineer, which the Goal of the note excludes. Doubtful: the Interface lists exactly the paths that the file now has, so the gap is in the note.

## Checks run

Each `factory` command ran with a binary built from the branch, whose code equals e30dda0, unless the line says otherwise. `simbase` is `origin/main` with the design notes of #67 and #68 merged; it is not an ancestor of the branch, so its merge base with the branch is 5af365e, the head of #67, and the note of 06 counts as changed against it.

- `go vet ./...`: no output
- `go test $(go list ./... | grep -v /test/acceptance)`: no package outside the acceptance tests has test files
- `go test -count=1 -v ./test/acceptance/`: `ok  	github.com/HannesT117/software-factory/test/acceptance`, 47 scripts pass
- `go test -count=1 -run 'TestAcceptance/06-gate-and-status/scenario-06.11$' ./test/acceptance/` with `crossreview.go:183` set to `instructions = codeInstructions`, in a scratch worktree: `ok`
- `factory check-design docs/increments/06-gate-and-status`: `check-design: ok 15 checked`
- `factory check-design --base simbase`: `check-design: ok 15 checked`; `--base origin/main`: `check-design: ok 23 checked`
- `factory check-tests --base simbase`: `check-tests: ok 47 checked`; `--base origin/main`: `check-tests: 16 errors`, because the notes of #67 and #68 are not on `origin/main`
- `factory check-pr --base simbase`: `check-pr: 46 errors`, `mixed-pr` on every changed path, with `.github/CODEOWNERS`, as the engineer allowed, and with the note of 06 and `docs/research/2026-11-factory-review.md` from the merge base above; `--base origin/main`: `check-pr: 47 errors`
- `factory check-review --base simbase` before these results: `check-review: 1 errors`, `missing-review` for e30dda011b555807d757e4a6753d5402b888d486; with them committed: `check-review: reviewed e30dda011b555807d757e4a6753d5402b888d486 in evidence/05-tests-per-increment/review.md, round 1` and `check-review: ok 1 checked`
- `factory status` before these results: `status: ok 5 checked`, with 05 and 06 at stage `review` and `status: review-open e30dda011b555807d757e4a6753d5402b888d486`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before these results: blocks and names e30dda011b555807d757e4a6753d5402b888d486
- with `factory` built from `simbase`, as CI builds it from the base: `check-design --base simbase` ends with a usage error, and `check-review` is no command yet, as the Rollout says
- the workflow patch of the first comment of the pull request applies with `git apply --check` on the branch
- `git merge-tree --write-tree HEAD origin/main` and `HEAD simbase`: no conflict

## Demonstrated scenarios

- **scenario 06.13**: `evidence/06-gate-and-status/acceptance-run.md`. `go test -count=1 -v ./test/acceptance/` at e30dda0 prints `acceptance: coverage 85.3% of statements` and three functions that no script reaches: `internal/report/report.go:44 Warning`, `internal/setup/setup.go:62 Error` and `internal/skills/skills.go:140 lastLine`; all 47 scripts pass.
