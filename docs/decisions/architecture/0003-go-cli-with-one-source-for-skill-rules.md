---
date: 2026-11-01
---

# Factory tooling as a Go CLI with thin adapters, and one source for skill rules

## Context and Problem Statement

The factory has checks, hooks and templates. The checks run in three places: in agent hooks, on the engineer's machine, and in CI. Agent tools call hooks in different ways, and they change faster than the projects that use them. The skill `authoring-skills`, which agents use to write skills, brings its own validator for skills. In which language and form is the factory tooling built, so that the agent tool stays replaceable, and where do the rules for skills live?

## Decision Drivers

* Hooks have no cross-tool standard. Each agent tool has its own hook input and output.
* A compiled binary needs no language runtime in the project. A script needs its runtime in every project, cloud session and CI job.
* Startup time, measured in a Linux container: a Go binary about 2 ms, a Node script about 30 ms. Model latency is seconds, so this difference alone does not decide.
* The CLI also decides which pull requests may merge: it checks path groups and test files. A checker with that role must stay small enough to read line by line in a review, and it must not depend on third-party code at runtime.
* The validator of `authoring-skills` is a Python script with PyYAML, derived from Anthropic's `skill-creator`. It checks the fields of the Agent Skills specification and warns about deviations from the skill authoring best practices. The skill also works without the factory.
* On 2026-11-01, `factory validate-skills` and that validator differed in five rules: the folder name, agent tool names, reserved words, angle brackets in the description, and one `SKILL.md` per skill.

## Considered Options

1. Go CLI `factory` that runs the validator of `authoring-skills` for skill rules
2. Go CLI `factory` with its own skill rules; `authoring-skills` calls `factory validate-skills`
3. Go CLI `factory` and the validator, both reading one shared rule file
4. Node or TypeScript scripts
5. TypeScript compiled into one binary (for example `deno compile`, as Swamp does)
6. Rust CLI
7. Existing tools only: GitHub settings and shell scripts

## Decision Outcome

Chosen option 1, because the rules for skills then have one source that works with and without the factory, while the checks that decide merges stay in Go without third-party code.

* **CLI:** one Go binary `factory`. It uses only the Go standard library at runtime. It calls `git`, and for `validate-skills` it calls `python3`. Test-only dependencies are allowed. The policy logic lives in tool-neutral packages.
* **Skill rules:** only in `core/skills/authoring-skills/scripts/quick_validate.py`. `factory validate-skills` carries that script, runs it, and reports its findings in the factory's output format.
* **Adapters:** `adapters/<tool>/` translates the tool's hook input and output and calls `factory`. Only adapters name a specific agent tool.

### Consequences

* Good, because a skill gets the same findings from the skill `authoring-skills` and from CI.
* Good, because the checks of path groups and test files need no runtime in the project.
* Good, because one tool serves the engineer, the hooks and CI.
* Good, because replacing the agent tool means adding one adapter.
* Bad, because `factory validate-skills` needs Python 3.10 or later with PyYAML wherever it runs.
* Bad, because a compiled tool needs releases and an install step in every environment ([ADR 0002](0002-delivery-as-pinned-copies.md)).

### Confirmation

* For the same skill, `factory validate-skills` and `quick_validate.py` report the same findings (design note of increment 02).
* `factory validate-skills` reports `tool-name` when a skill names a specific agent tool (scenario 02.14).
* The review checks that the tool-neutral packages name no agent tool.

## Pros and Cons of the Options

### Go CLI that runs the validator

* Good, because each skill rule exists once.
* Good, because the skill keeps its own validator and stays usable without the factory.
* Bad, because the skill check needs Python and PyYAML.

### Go CLI with its own skill rules

* Good, because the factory needs no runtime in the project.
* Bad, because the skill can check a draft only where `factory` is installed.
* Bad, because the style warnings of the best practices would have to be rewritten in Go.

### Go CLI and validator with a shared rule file

* Good, because limits, fields and patterns exist once.
* Bad, because two implementations of the logic around that data can still drift apart.

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

* **Confidence:** medium. The startup times and the differences between the two skill checkers are measured.
* **Revisit when** a cross-tool hook standard appears, the checks grow beyond what the engineer can read line by line, or Python is missing in an environment that must check skills.
* **Sources:** Claude Code docs on hooks; Swamp source repository; startup times measured on 2026-10-26; Anthropic `skill-creator` (Apache License 2.0); comparison of the two skill checkers on 2026-11-01. Details are in the [research notes](../../research/2026-10-factory-basis.md).
