# Review of increment 02: trust core

Reviewed commit: 858863105af6df0ef668e4ce63c880aed5e17b7d
Pull request: https://github.com/HannesT117/software-factory/pull/39
Result: no findings

## Findings

None.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `go test ./test/acceptance/` on bc39946, the last commit with tests: all 22 tests fail, so no test passes without the code
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/02-trust-core`: `check-design: ok 31 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-tests --base origin/main`: `check-tests: ok 22 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 321 checked`, after `check-pr: tests end at bc399465790e3f6c91eb9ffa44b9080749321a08`
- `factory hook pre-tool-use` with `cwd` `/tmp` and a `Write` to a new folder under `test/acceptance/` of the repo: denied
- `factory hook stop` in a sample repo where a merge from `main` brings a change to `b.go`: blocks and names the merge commit

## Demonstrated scenarios

- **scenario 02.15**: not run yet; the reviewer offered the run to the engineer, because it starts agent sessions for each skill with and without the skill.
