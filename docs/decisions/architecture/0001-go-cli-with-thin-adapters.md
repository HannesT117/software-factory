---
date: 2026-11-01
---

# Factory tooling as a Go CLI with thin adapters per agent tool

## Context and Problem Statement

The factory has checks, hooks and templates. The checks run in three places: in agent hooks, on the engineer's machine, and in CI. Agent tools call hooks in different ways, and they change faster than the projects that use them. In which language and form is the factory tooling built, so that the agent tool stays replaceable?

## Decision Drivers

* Hooks have no cross-tool standard. Each agent tool has its own hook input and output.
* A compiled binary needs no language runtime in the project. A script needs its runtime in every project, cloud session and CI job.
* Startup time, measured in a Linux container: a Go binary about 2 ms, a Node script about 30 ms. Model latency is seconds, so this difference alone does not decide.
* The CLI also decides which pull requests may merge: it checks path groups and test files. A checker with that role must stay small enough to read line by line in a review, and it must not depend on third-party code at runtime.
* The skill `authoring-skills`, which agents use to write skills, has its own validator: a Python script with PyYAML, derived from Anthropic's `skill-creator`. The skill also works without the factory. On 2026-11-01, `factory validate-skills` and that validator differed in five rules.

## Considered Options

1. Go CLI `factory`
2. Node or TypeScript scripts
3. TypeScript compiled into one binary (for example `deno compile`, as Swamp does)
4. Rust CLI
5. Existing tools only: GitHub settings and shell scripts

## Decision Outcome

Chosen option 1, because its checks need no runtime in the project, gives the checks one tested home, and its standard library covers what the checks need without third-party code.

* **CLI:** one Go binary `factory`. It uses only the Go standard library at runtime. It calls `git`, and for `validate-skills` it calls `python3` with the validator of the skill `authoring-skills`, which the binary carries. That validator is the only place for the rules of skills. Test-only dependencies are allowed. The policy logic lives in tool-neutral packages.
* **Adapters:** `adapters/<tool>/` translates the tool's hook input and output and calls `factory`. Only adapters name a specific agent tool.

### Consequences

* Good, because the project needs no runtime for the checks of path groups, test files, design notes and decision records.
* Good, because the skill `authoring-skills` and CI apply the same rules to a skill.
* Good, because one tool serves the engineer, the hooks and CI.
* Good, because replacing the agent tool means adding one adapter.
* Bad, because a compiled tool needs releases and an install step in every environment ([ADR 0002](0002-delivery-as-pinned-copies.md)).
* Bad, because `factory validate-skills` needs Python 3.10 or later with PyYAML wherever it runs.

### Confirmation

* `factory validate-skills` reports `tool-name` when a skill names a specific agent tool (scenario 02.14).
* For the same skill, `factory validate-skills` and the validator of `authoring-skills` report the same findings (scenario 02.32).
* The review checks that the tool-neutral packages name no agent tool.

## Pros and Cons of the Options

### Go CLI

* Good, because one static binary runs on macOS, Linux and Windows without a runtime.
* Good, because the standard library covers files, JSON, regular expressions and processes, so the checker stays free of runtime dependencies.
* Bad, because releases and an install step are needed.

### Node or TypeScript scripts

* Good, because a TypeScript project already has Node, so nothing needs to be released.
* Bad, because every project, cloud session and CI job needs Node.

### TypeScript compiled into one binary

* Good, because the language matches projects written in TypeScript.
* Bad, because each binary carries the whole Deno runtime.

### Rust CLI

* Good, because it has the same distribution model and stronger types.
* Bad, because Rust's standard library has no JSON parser and no regular expressions, so the checks would depend on crates such as `serde_json` and `regex`.

### Existing tools only

* Good, because there is no own binary to build, release and install.
* Bad, because shell scripts behave differently on macOS, Linux and Windows and are hard to test.
* Bad, because each agent tool needs its own hook scripts, with no shared logic.

## More Information

* **Confidence:** medium. The startup times and the differences between the two skill checkers are measured. The factory has no Go code on `main` yet.
* **Skill rules:** two other places were considered. With its own rules in Go, the skill could check a draft only where `factory` is installed, and its style warnings would need a rewrite in Go. A rule file shared by both checkers would leave two implementations of the logic around it.
* **Revisit when** a cross-tool hook standard appears, the checks grow beyond what the engineer can read line by line, or Python is missing in an environment that must check skills.
* **Sources:** Claude Code docs on hooks; Swamp source repository; startup times measured on 2026-10-26; Anthropic `skill-creator` (Apache License 2.0); comparison of the two skill checkers on 2026-11-01. Details are in the [research notes](../../research/2026-10-factory-basis.md).
