# Review of increment 03: setup script

Reviewed commit: 865d49acba33b428207b91548f43bf2f0d881058
Pull request: none yet, branch `claude/03-setup-script-3`
Result: no findings

## Findings

None.

The tests come first: the test commits run from 7c572ce to 8e30899, the first code commit is 8dad2f8, and no commit after it changes a file under `test/`. At 8e30899 every test of 03.1 to 03.8 fails, because the command `setup` does not exist yet.

8e30899 runs the second run of 03.2, the dry run of 03.4 and the second run of 03.8 through `stdin-probe` (`test/acceptance/stdinprobe_test.go:23`). It gives the command a closed pipe with one line and fails with exit 99 if any byte is gone after the exit. This checks the note's "does not read standard input" (scenarios 03.2, 03.4, 03.8), not only the missing announcement. Three changes to `internal/setup/setup.go`, each reverted afterwards, make the tests fail: one byte read at the start of every run fails 03.1, 03.2, 03.4 and 03.8; a read in the step `secret` when the secret exists fails 03.2 and 03.8 with `stdin-probe: factory read 41 bytes of standard input`; a read in the dry run fails 03.4 the same way. Without a change, all pass.

The code in `cmd/` and `internal/` equals that of the old branch `claude/03-setup-script`, which round 3 reviewed at its paging commit; `ruleset main` compares by containment, as the note on `main` says since #56.

The cherry-picks keep `main`'s content where `main` changed it. `docs/setup/README.md:13`, step 4, holds `main`'s rule that a workflow's secrets are environment secrets, never repo secrets, now for the environments other than `cross-review`, which step 3 creates. Step 2 takes over the copy of the workflows from `main`'s old step 4. `README.md` changes only the version line and the row of increment 03; the rows of 02, 04 and 05 equal `main`. No doc names `docs/setup/ruleset.json` or `ruleset-branches.json` any more.

## Checks run

- `go vet ./...`: no output
- `go test -count=1 ./...`: `ok github.com/HannesT117/software-factory/test/acceptance`
- `go test -count=1 ./test/acceptance/` at 8e30899: 03.1 to 03.8 fail
- `go test -count=1 ./test/acceptance/` with each of the three stdin reads above: fails as named; restored: ok
- `factory check-pr --base origin/main`: `check-pr: ok 25 checked`, tests end at 8e30899
- `factory check-tests --base origin/main`: `check-tests: 6 errors`, the replaced test of 02.7 and the missing tests of 05.1 to 05.5; `main` reports the same 6 and also the 8 missing tests of 03, 14 in all, so the branch adds none
- `factory check-design docs/increments/03-setup-script`: `check-design: ok 9 checked`
- `factory check-adr`: `check-adr: ok 4 checked`
- `factory validate-skills --strict core/skills`: `validate-skills: ok 6 checked`
- `factory hook stop` before this result: it asked for a result for 865d49acba33b428207b91548f43bf2f0d881058

## Demonstrated scenarios

None; the note has no demonstrated scenario. Scenario 03.9 is manual.
