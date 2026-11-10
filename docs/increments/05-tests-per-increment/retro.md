# Retro of increment 05

The code of increment 05 came in the pull request of increment 06; its retro covers what differs from there.

## Metrics

- Pull requests: the design note (#60), the cross-review per increment (#64), the list of unchanged behaviour (#67), and the tests and code together with increment 06.
- Engineer decisions in the conversation: two, that tests which pass before the code are a general problem, and the mark "unchanged behaviour" in the note.
- Tests: 7 test files from a separate session; 05.6 was replaced by 06.9 before the code.

## Review sources

| Source | Findings | Findings that changed code or a note | Cost |
|---|---|---|---|
| Other model on #67 | 1 note | 0: the Interface already fixed the point | one call |
| Test writer | 1: the tests of 05.4 and 05.5 pass before the code | 1: the section Unchanged behaviour | one session |

## Agent version

- Claude Code 2.1.322; model `claude-opus-5-5`.

## Corrections

- The rule that every added test fails before the code did not fit a scenario that fixes what a narrower rule keeps. The note now lists such scenarios, and `tests-start-red` reads the list from the approved note.

## Points for the factory

- While the tests of one increment waited, `checks` was red for every pull request with code; two pull requests merged past it without a record. Increment 06 adds the record.

## Summary

`factory check-tests --base` now asks for missing tests only in the pull requests of their increment and names the increments that wait. `cross-review` sends only the notes of the pull request's own increments. A note can list the scenarios whose tests pass before the code. Open: the demonstration of 05.8 on the pull request.
