---
name: authoring-skills
description: Writes, revises and tests Agent Skills, from frontmatter and body to reference files, scripts and evals. The agent uses it whenever it creates or changes a skill, turns a workflow into a skill, runs evals on a skill or tunes a skill's description. Not for plugins, hooks, MCP servers or agent instruction files such as AGENTS.md.
license: Apache-2.0. LICENSE.txt has the complete terms; NOTICE.txt lists the changes to the original skill-creator.
---

# Authoring skills

Where these rules differ from skill-creator, these rules apply.

## Sources

- The agent reads the project's rules on the remote's default branch, after a fetch, with `git show origin/<default-branch>:<path>`; a working branch or an open pull request can be out of date.
- The skill points to the file and section that hold a rule instead of copying it. A copy is only for a skill that runs without the repository; it is verbatim and names its source path and commit.
- Each rule appears once, in the wording of its source.
- Templates contain only lines of the output. Their headings match the source template in name and order.
- The agent keeps a trace table, each domain rule of the skill next to the file and section that state it, and shows it to the person with the draft. A rule without a source is a question for the person.

## Voice

- The agent is "the agent", never a product name; helpers are "subagents"; the human is "the person". Prompts under `agents/` address their subagent as "you".
- The description says what the skill does, when the agent uses it, and what it is not for when a neighboring skill exists, in the third person.
- A skill that has to load every time also gets a one-line pointer in the project's agent instruction file, because a description cannot guarantee loading.
- Where the project's own style rules may or may not apply to skills, the agent asks the person.

## Checks

The agent runs both from this skill's folder after each draft, and fixes each error and each warning it cannot justify:

```bash
python -m scripts.quick_validate <path-to-skill>
python -m scripts.check_sources <path-to-skill> --repo <repo-root> --ref origin/<default-branch>
```

`check_sources` also takes `--template <file-in-skill>=<path-in-repo>` and `--ignore <glob>`. Then the agent goes through [references/best-practices.md](references/best-practices.md). The scripts need Python 3.10 or later with PyYAML, and git.

## Tests

Before the first run, the agent reads [references/testing.md](references/testing.md) completely; [references/schemas.md](references/schemas.md) holds the file formats. Each case runs with the skill and without it (or with the old version), all started together, on each model the skill will run on. The person sees the review viewer before the agent judges the outputs.

- [references/description-tuning.md](references/description-tuning.md): after the skill is stable and the person agrees.
- [references/environments.md](references/environments.md): without subagents, without a display, or when an installed skill changes.
- [agents/grader.md](agents/grader.md), [agents/comparator.md](agents/comparator.md), [agents/analyzer.md](agents/analyzer.md): prompts for the subagents that grade, compare blind and analyze.

## Packaging

`python -m scripts.package_skill <path-to-skill>` validates the skill and writes a `.skill` file without `evals/`.
