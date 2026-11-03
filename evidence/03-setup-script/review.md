# Review of increment 03: setup script

Reviewed commit: 62b966ca6dc2ca762b2c358a35ab684aa24bbcc6
Pull request: https://github.com/HannesT117/software-factory/pull/54
Result: 1 finding

## Findings

- **Special cases**, `internal/setup/setup.go:368` (with `internal/setup/compare.go:34`): the note's step `ruleset main` (and `ruleset branches`) requires that the ruleset *equals* the embedded file in `name`, `target`, `enforcement`, `bypass_actors`, `conditions` and `rules`. The code compares by containment: an object in GitHub's ruleset may hold keys that the file lacks. A key added by hand is therefore neither reported nor repaired, against the goal "a second run repairs settings that changed by hand" and the failure signal "after a run, a dry run reports a change". Example, checked with a probe against `contains`: a `required_reviewers` entry, or `automatic_copilot_code_review_enabled: true`, added to the parameters of the rule `pull_request` in `main` leaves the step `unchanged`. Doubtful because strict equality would report every real ruleset as changed, since GitHub returns default keys that the file does not carry; the engineer decides whether the note's "equals" becomes "holds", or the files pin GitHub's defaults so that the code can compare strictly.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `go test -run TestAcceptance/03 ./test/acceptance` on 70ff719, the last commit with tests: all 8 tests of scenarios 03.1 to 03.8 fail, so no test passes without the code
- `factory check-pr --base origin/main`: `check-pr: ok 21 checked`, after `check-pr: tests end at 70ff719b35090bd6503956ac296b4062cb5ff330`
- `factory check-tests --base origin/main`: `check-tests: ok 30 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design`: does not apply; the pull request changes no design note
- Commits after 70ff719 touch no test, CI, `.factory/` or design note; `internal/setup/` reads no environment variable and has no branch for a test run

## Demonstrated scenarios

None; the note has no demonstrated scenario. Scenario 03.9 is manual.
