# Recording of scenario 05.8

The Python of the step "Each added test fails there, unless its approved note lists it as unchanged, and every other test passes" of the job `tests-start-red`, as `.github/workflows/test.yml` holds it after the patch in the first comment of pull request #69 (patch SHA-256 `13cea05f577ac2374af4c2499b449f6e46f0cc2db9f3036479f6e7c7758b6f25`). The script is the heredoc between `<<'PY'` and `PY`, without its indentation of 10 spaces (SHA-256 `131b92e2dd384905d5da9936efd8aabbf719c4a7b90bf3231c82c9b723b91834`). It ran as the job runs it, with Python 3 and Go 1.25.0.

## Where the tests end

- Base: `BASE=simbase`, a27bd133bd18135d77c99f73db9d8ae9b1322c22, which is `origin/main` (2461855c4630246532a944a6f53869eb3e38051c) with the design notes of #67 and #68 merged.
- End of the tests: dbe5e8db3403affb12c89cee5d6388c6b805321f, the last commit before the first code commit c1f93d8; the merge 468c197 between them brings in only the two design notes.
- `factory check-pr --base simbase` on the branch prints no line `tests end at`: the pull request changes `.github/CODEOWNERS`, so it mixes path groups, and the command stops before it looks for the end of the tests. In a scratch clone, a commit with the tree of `simbase` on top of `origin/main`, as #67 and #68 will land, merged into c1f93d8, gives `check-pr: tests end at dbe5e8db3403affb12c89cee5d6388c6b805321f` and `check-pr: ok 25 checked`. At the head of the branch, merged the same way, it again gives no such line.

## Inputs

- Added tests: `git diff --name-only --diff-filter=A "$BASE...$END" -- 'test/acceptance/*.txtar'` at the end commit lists 18 files: 05.1 to 05.5, 05.7, 06.1 to 06.12. `scenario-06.15.txtar` is missing from it: git pairs it with the deleted `scenario-02.26.txtar` as a rename (`R051`).
- Results: `go test -count=1 -json ./test/acceptance/` at the end commit, 47 scripts. Of the added tests, 05.4, 05.5 and 06.15 pass, the 16 others fail; the 28 tests of increments 02 and 03 pass.
- Unchanged behaviour on `simbase`: 05.4 and 05.5 in the note of 05, 06.15 in the note of 06.

## Run 1: the results at the end of the tests

```
$ BASE=simbase python3 -I tests-start-red.py added.txt tests.json
tests-start-red: ok 47 checked
exit 0
```

## Run 2: the same results, with the test of 05.4 set to failed

The one event `{"Action": "pass", "Test": "TestAcceptance/05-tests-per-increment/scenario-05.4"}` of `tests.json` became `fail`; nothing else changed.

```
$ BASE=simbase python3 -I tests-start-red.py added.txt tests-05.4-failed.json
tests-start-red: error TestAcceptance/05-tests-per-increment/scenario-05.4: fail, want pass
tests-start-red: 1 errors
exit 1
```

## Run 3: the results of run 1 against a base without the lists

```
$ BASE=origin/main python3 -I tests-start-red.py added.txt tests.json
tests-start-red: error TestAcceptance/05-tests-per-increment/scenario-05.4: pass, want fail
tests-start-red: error TestAcceptance/05-tests-per-increment/scenario-05.5: pass, want fail
tests-start-red: 2 errors
exit 1
```

This is the order that the Rollout of the note names: without the note on the base branch, the tests of 05.4 and 05.5 get `pass, want fail`.

## In CI

On this pull request the job finds no end commit, because `check-pr` names none for a pull request that mixes path groups; the step above is then skipped and the job passes without a check. These runs are the only ones of the new rule on the tests of increments 05 and 06.
