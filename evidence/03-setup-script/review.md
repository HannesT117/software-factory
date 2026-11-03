# Review of increment 03: setup script

Reviewed commit: 21f81eb27790eaaace2b6937b7832b25a2f8a6e6
Pull request: https://github.com/HannesT117/software-factory/pull/54
Result: 1 finding

## Findings

- **Failing checks**, `.github/workflows/checks.yml:14` (with `internal/prcheck/prcheck.go:144`): the required check `checks` fails on every commit of the pull request since the first code commit, at 21f81eb too, with `check-pr: 11 errors`: each test file is "changed in commit 414296aff62f after the first commit outside the tests group". The workflow checks out GitHub's merge ref, so `--first-parent` from the merge base sees one merge commit that changes tests and code together, and the pull request reads as mixed. On the branch head, `check-pr` passes. Rebuilt locally: a `--no-ff` merge of 21f81eb into origin/main gives the same 11 errors with main's `factory`. The pull request cannot merge while the check fails. Doubtful because the cause is on `main`, not in this diff: the workflow (for example `ref: ${{ github.event.pull_request.head.sha }}` at the checkout) or `check-pr`'s reading of a merge ref needs a fix in a pull request of its own, and every pull request with tests and code hits it.

Round 1's finding (rulesets compared by containment while the note said "equals") is no longer a finding: design pull request #56 changes the step `ruleset main` to "holds every key and value that the file sets", and says that keys the file does not set are not compared. `contains` in `internal/setup/compare.go:34` does that. This holds once #56 is on `main`; until then the note on `main` still says "equals".

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- Commits after round 1 (5a0ca33, 7602992, 21f81eb) change only `docs/setup/README.md`, `docs/increments/03-setup-script/retro.md` and `README.md`; no code, test, CI, `.factory/` or design note. The tests still end at 70ff719, where round 1 saw all 8 tests of 03.1 to 03.8 fail without the code.
- `factory check-pr --base origin/main` (main's binary) on the branch head: `check-pr: ok 24 checked`; on a merge of the branch into origin/main, as in CI: `check-pr: 11 errors`
- `factory check-tests --base origin/main`: `check-tests: ok 30 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/03-setup-script` on #56's branch: `check-design: ok 9 checked`; this pull request changes no design note
- `factory hook stop` before this result: it asked for a result for eda89847bf8c7edcefb9d4736b1a5030fe8753f7, the latest commit with code; 21f81eb contains it and changes only docs after it, so this result covers it. Round 1 named 62b966c, which changes no code, so the hook and the CI job `cross-review` did not accept it.
- GitHub checks at 21f81eb: `test` and `tests-start-red` pass; `checks` fails as above; `cross-review` waited for a review result for eda8984

## Demonstrated scenarios

None; the note has no demonstrated scenario. Scenario 03.9 is manual.
