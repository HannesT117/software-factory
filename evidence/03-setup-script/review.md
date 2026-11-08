# Review of increment 03: setup script

Reviewed commit: 5ef86b24672307255ac2ffd51a80094a479456cd
Pull request: https://github.com/HannesT117/software-factory/pull/54
Result: no findings

## Findings

None.

`internal/setup/gh.go:100` (`list`) reads invitations, deployment branch policies and rulesets with `gh api --paginate`, as `internal/setup/setup.go:217`, `:258` and `:342` call it. `gh api --help` of gh 2.89.0 states the format it relies on: "Each page is a separate JSON array or object", and only `--slurp` wraps the pages into one outer array. `list` decodes the output as a stream of JSON documents until `io.EOF`: an array endpoint gives one array per page, and the branch-policies endpoint gives one object per page whose `branch_policies` holds the items. Empty output and `[]` give no items. Pagination follows GitHub's `Link` header, which carries `per_page=100` on to later pages. The note's rule that a second run reaches the target state (design note, Goal) now holds for repos with more than 100 invitations, branch policies or rulesets; no input is treated differently from that rule.

A page without the key counts as an empty page, as before the change, where `lookup` also gave nil. The fake gh in `test/acceptance/gh_test.go:497` ignores `--paginate` and answers one page, so no approved scenario covers more than one page; the note asks for none.

Round 2's note on the ruleset comparison is settled: design pull request #56 is on `main`, and the step `ruleset main` there says "holds every key and value that `internal/setup/rulesets/main.json` sets", which `contains` in `internal/setup/compare.go` does.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- Commit 5ef86b2 changes only `internal/setup/gh.go` and `internal/setup/setup.go`; no test, CI, `.factory/` file or design note. The commits between round 2 and it change only `docs/increments/03-setup-script/retro.md` and this file.
- `list` against a stub gh that prints `[{"id":1},{"id":2}][{"id":3}]`, two `branch_policies` objects on separate lines, and `[]`: 3, 3 and 0 items (agent test, scaffolding, not committed)
- `factory check-pr --base origin/main` (main's binary): `check-pr: ok 24 checked`
- `factory check-tests --base origin/main` (main's binary): `check-tests: ok 30 checked`
- `factory check-design docs/increments/03-setup-script` (main's binary): `check-design: ok 9 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before this result: it asked for a result for 5ef86b24672307255ac2ffd51a80094a479456cd
- The real paging of the GitHub API and the GitHub checks at 5ef86b2 were not run: gh's GitHub access from this session answered HTTP 403.

## Demonstrated scenarios

None; the note has no demonstrated scenario. Scenario 03.9 is manual.
