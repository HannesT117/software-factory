# Increment 02: shout

## Goal

- `greet --shout <name>` greets a person in capital letters, for noisy rooms.

## Out of scope

- other languages

## Existing tools

| Part | Checked | Decision | Reason |
|---|---|---|---|
| Capital letters | Go `strings.ToUpper` | used | It handles names with letters outside ASCII. |

## Interface

`greet --shout <name>` prints `greet: HELLO, <NAME>!`, with the name in capital letters, and exits with 0. `greet --shout` without a name prints `greet: error args: missing-name` and exits with 2. An unknown flag prints `greet: error args: unknown-flag <flag>` and exits with 2. All output goes to stdout.

## Scenarios

- **scenario 02.1** (automated) **A shouted greeting.** `greet --shout Ada` greets Ada in capital letters.
- **scenario 02.2** (automated) **A shout needs a name.** `greet --shout` without a name fails and says that the name is missing.
- **scenario 02.3** (automated) **An unknown flag is an error.** `greet --whisper Ada` fails and names the flag.

## Evidence

- **Approved scenarios:** 02.1 to 02.3.

## Rollout

- **Off switch:** none; the flag is new and optional.
- **Failure signal:** a user reports a shouted greeting in small letters, or the error line `unknown-flag --shout`.
