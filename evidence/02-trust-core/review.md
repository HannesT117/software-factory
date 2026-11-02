# Review of increment 02: trust core

Reviewed commit: 6dcdf2858c229a2e1908d3e3b9ffe8efbc457ec6
Pull request: https://github.com/HannesT117/software-factory/pull/35
Result: 4 findings

## Findings

- **Special cases**, `internal/paths/paths.go:140`: `Code.Is` never counts a path under `evidence/` as code. The stop hook and `cross-review` use it. The note on `main` says, under Code paths: "Without the file, every path outside the groups counts", and `evidence/` is outside the groups. So the code follows the open design pull request #34, not the note on `main`, which is the measure. If #34 changes or does not merge, code and note differ. If #34 merges unchanged, this finding no longer applies.
- **Protected paths, early feedback**, `internal/guard/guard.go:26`: the guard compares paths with their exact case. A `Write` to `.GitHub/workflows/x.yml`, `Test/Acceptance/02-x/scenario-02.1.txtar` or `GO.MOD` is allowed; the binary of this commit confirms it. On a case-insensitive file system, such as the default on macOS, these writes change the protected files. Doubtful as a defect: the note names the paths in lower case, and the guard is a soft part (note, Limits). `check-pr` in CI still sees the real path. But local sessions on macOS are in scope (scenario 02.12), and there the guard gives no early feedback for these writes.
- **Tests that miss their scenario**, `test/acceptance/02-trust-core/scenario-02.26.txtar:26`: each blocking case runs `sh -c 'factory cross-review --base main; exit 0'`, so the test never checks the exit code of a block. The check `cross-review` in CI goes red only through the exit code. Code that prints `cross-review: blocking` and exits with 0 passes this test. Doubtful: the scenario names the output, not exit codes, and the Interface gives no exit code for `cross-review: blocking`. The code exits with 1 (`internal/crossreview/crossreview.go:98`).
- **Work outside the scenarios**, `docs/principles.md:54`: rule 37 drops "(72 at most)" from the limit of a title. No goal or scenario of the note asks for this, and the description of the pull request does not name it. The first comment names "the title rule" from the review of #30, but not this change to the rule.

## Checks run

- `go vet ./...`: no output
- `go test ./...`: `ok github.com/HannesT117/software-factory/test/acceptance 2.774s`
- `factory check-pr --base main`: `check-pr: ok 332 checked`, after `check-pr: tests end at d47ed7b7c097c9365f4cf18576463451a5d04b6d`
- `factory check-tests --base main`: `check-tests: ok 23 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory check-design docs/increments/02-trust-core`: `check-design: ok 32 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `go test ./...` on d47ed7b, the last commit with tests: all 23 scenario tests fail, and there are no other tests

## Demonstrated scenarios

- **scenario 02.15**: not run. The evals of the six skills are billed model runs, and `adapters/claude-code/authoring-skills/README.md` says that the agent offers them to the engineer instead of starting them unasked. No recording exists.

## Comprehension questions

- **question 1** This pull request changes `AGENTS.md` and `docs/principles.md`, which the engineer owns as code owner. Why does `factory check-pr` accept it, and what keeps it from merging without the engineer?
- **question 2** The stop hook keeps blocking an implementing session, although `evidence/02-trust-core/review.md` exists and names a commit. Where does the cause lie, and which commands show it?
- **question 3** `factory cross-review` reads `.factory/cross-review.json` from the base branch, not from the head of the pull request. Which misuse does that stop, and what happens on the pull request that adds the file for the first time?
