# Checklist before a skill is finished

## Contents

- Sources and domain rules
- Frontmatter and triggering
- Style
- Structure
- Scripts and code
- Testing

This checklist follows Anthropic's "Skill authoring best practices" (platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices) and adds checks for problems the best practices do not cover: wrong or outdated domain content, copied rules, and the wrong voice. The validator (`scripts/quick_validate.py`) and the source check (`scripts/check_sources.py`) cover the items marked *(script)*; the agent checks the others by reading.

## Sources and domain rules

- [ ] Every rule the skill states about the domain is in the trace table, with the file and section that hold it.
- [ ] The agent read each source on the remote's default branch after fetching, not in the working tree or an open pull request.
- [ ] The person has seen the trace table. Rules without a source were asked about or removed.
- [ ] The skill points to sources instead of copying regexes, formats or section lists. *(script: verbatim overlap with repository documents)*
- [ ] Where a copy is unavoidable, it is verbatim and names its source path and commit.
- [ ] Every repository path the skill names exists on the default branch. *(script)*
- [ ] Templates match their source templates in heading names and order. *(script, with `--template`)*
- [ ] Template blocks contain only lines that belong in the output; notes to the author are outside the block.
- [ ] Each rule appears once and in one form, with no second version of a limit in parentheses.

## Frontmatter and triggering

- [ ] `name` keeps the limits of the Agent Skills specification. *(script)*
- [ ] `name` is a gerund, or follows the pattern of the skills around it. *(script: warns when it is not a gerund)*
- [ ] `description` keeps the limits of the Agent Skills specification. *(script)*
- [ ] `description` is in the third person: no "I", "you" or "Use this skill". *(script)*
- [ ] `description` says what the skill does and when the agent uses it, with the key terms a person would type, and what it is not for when a neighboring skill exists.
- [ ] When the skill must load every time, the person knows that a pointer in the project's agent instruction file, such as `AGENTS.md`, is the reliable way.

## Style

- [ ] Declarative, third person: "The agent does X when Y, because Z". *(script: warns about "you", "the user", "the engineer", "the developer" outside code blocks)*
- [ ] One fixed term for the agent throughout: "the agent", never a product name; helpers are "subagents". The human is "the person".
- [ ] Project style rules (AGENTS.md, CONTRIBUTING.md) were checked, and unclear cases were asked about.
- [ ] Ordered steps only where the order matters; each step says what it achieves.
- [ ] Concise: nothing the agent already knows.
- [ ] Consistent terminology throughout.
- [ ] One default per decision instead of a list of options.
- [ ] No time-sensitive statements; superseded methods are in an "old patterns" section. *(script: warns about date phrases)*
- [ ] Strong words (MUST, NEVER) are rare and have their reason next to them.
- [ ] Examples are concrete, not abstract.
- [ ] MCP tools are named in full (`ServerName:tool_name`).

## Structure

- [ ] The `SKILL.md` body keeps the line limit of the best practices. *(script)*
- [ ] Details are in separate files where needed, and each is linked from `SKILL.md` with a word on when to read it.
- [ ] References are one level deep: `SKILL.md` links every reference file, and no file is reachable only through another reference file. *(script)*
- [ ] Long reference files start with a table of contents. *(script)*
- [ ] Descriptive file names; directories ordered by domain or feature.
- [ ] Forward slashes only; no absolute paths of one machine. *(script)*
- [ ] Workflows have clear steps; complex ones have a checklist the agent can copy into its answer.

## Scripts and code

- [ ] Scripts handle errors instead of passing them to the agent, with error messages that say how to fix the problem.
- [ ] No unexplained constants: each value has a comment with its reason.
- [ ] Required packages are listed, and their availability on the target platform is checked (some platforms install packages at run time, others have no network access).
- [ ] The skill says for each script whether the agent runs it or reads it.
- [ ] Critical operations have a validation step and a feedback loop (check, fix, check again).

## Testing

- [ ] At least three test cases exist.
- [ ] At least one test checks a domain rule against its source, not against the skill.
- [ ] Tests ran with the skill and against a baseline.
- [ ] Tests ran on every model the skill will be used with, or the person knows which models remain untested.
- [ ] Tests used real usage scenarios, and the person reviewed the outputs in the viewer.
- [ ] Feedback from the team was taken into account, where there is a team.
