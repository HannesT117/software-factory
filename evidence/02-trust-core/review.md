# Review of increment 02: trust core

Reviewed commit: 4035b49c2d10ccc78325e4ebf5451e30b2cb002f
Pull request: https://github.com/HannesT117/software-factory/pull/39
Result: 2 findings

## Findings

- **Special cases**, `adapters/claude-code/hooks.go:53`: the guard hook takes the project root from `cwd`, not from the file. When `cwd` lies outside the repo, the path relative to it does not start with `test/`, `.github/` or `.factory/`, so the hook allows the write. With `cwd` `/tmp` and `file_path` `<repo>/test/acceptance/02-trust-core/x.txtar`, `factory hook pre-tool-use` prints nothing and exits 0; with `cwd` `<repo>/internal` it denies the same write. Scenario 02.13 says the hook refuses edits to test files, without a condition on the working folder. The hook is a soft guard (note, Limits), and `check-pr` still catches the change, so this weakens only the early feedback.
- **Work outside the scenarios**, `docs/principles.md:54`: rule 37 changes the title limit from "at most 50 characters (72 at most)" to "at most 50 characters". No goal or scenario of the note asks for this, and the commit message of 72462be does not name it. Doubtful, because the skill `writing-pull-requests` uses the same limit, so the change can be an alignment; the engineer decides whether the stricter rule is wanted.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `go test ./test/acceptance/` on bc39946, the last commit with tests: all 22 tests fail, so no test passes without the code
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory check-design docs/increments/02-trust-core`: `check-design: ok 31 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-tests --base origin/main`: `check-tests: ok 22 checked`
- `factory check-pr --base origin/main`: `check-pr: ok 320 checked`, after `check-pr: tests end at bc399465790e3f6c91eb9ffa44b9080749321a08`

## Demonstrated scenarios

- **scenario 02.15**: not run yet; the reviewer offered the run to the engineer, because it starts agent sessions for each skill with and without the skill.
