# Recording of scenario 06.13

`go test -count=1 -v ./test/acceptance/` at 16e4808119f0ab47fd9adabffe583e1f021d587b, the latest code commit, with Go 1.25.0 on linux/amd64. The run builds `factory` with `-cover -coverpkg=./...` and gives every script the same `GOCOVERDIR`.

## The lines of the acceptance run

```
acceptance: coverage 85.4% of statements
acceptance: uncovered internal/report/report.go:44 Warning
acceptance: uncovered internal/setup/setup.go:62 Error
acceptance: uncovered internal/skills/skills.go:140 lastLine
```

## The result of the tests

47 scripts pass and none fails: 20 of increment 02, 8 of increment 03, 6 of increment 05 and 13 of increment 06. The last lines:

```
PASS
ok  	github.com/HannesT117/software-factory/test/acceptance	3.701s
```

The patched job `test` runs the same command and prints the lines that start with `acceptance: `, `--- FAIL`, `FAIL` or `ok `.
