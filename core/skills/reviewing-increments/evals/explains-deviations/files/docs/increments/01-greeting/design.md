# Increment 01: greeting

## Goal

- The command `greet` greets a person by name.
- `greet --known` lists the people that the project knows.

## Out of scope

- greetings in other languages

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Arguments | Go `os.Args`, Go `flag` | `os.Args` | One positional argument needs no parser. |

## Interface

`greet <name>` prints `greet: Hello, <name>!` and exits with 0. `greet --known` reads `names.txt`, one name per line, and prints one line `greet: known <name>` for each name, then exits with 0. Without a name, it prints `greet: error args: missing-name` and exits with 2. All output goes to stdout.

## Scenarios

- **scenario 01.1** (automated) **A name is greeted.** `greet Ada` greets Ada by her name.
- **scenario 01.2** (automated) **A missing name is an error.** `greet` without a name fails and says that the name is missing.
- **scenario 01.3** (manual) **The greeting reads well.** The owner reads three greetings and finds each one friendly.
- **scenario 01.4** (manual) **Known people are listed.** The owner runs `greet --known` and sees each name of `names.txt`.

## Evidence

- **Approved scenarios:** 01.1, 01.2.
- **Manual:** 01.3, 01.4.

## Rollout

- **Off switch:** none; the command has no users yet.
- **Failure signal:** a user reports a wrong greeting.
