# Environments

## Contents

- With subagents (the default)
- Without subagents
- Without a display
- Updating an installed skill

## With subagents

`SKILL.md` and `references/testing.md` describe this case: subagents run the tests in parallel, and the viewer runs as a local server or as a static file. Description tuning works where the factory's adapter for the harness provides a runner.

## Without subagents

The loop is the same (draft, test, review, improve, repeat), but without subagents some parts change:

- **Test cases.** Without subagents there is no parallel run. For each test, the agent reads the skill's `SKILL.md` and follows it to do the task itself, one test after the other. This is less strict than independent subagents, because the agent wrote the skill and knows everything in it, but it is a useful check, and the person's review makes up for part of the difference. The agent skips the baseline runs.
- **Review.** When the agent cannot open a browser, it skips the viewer and shows each prompt and its output in the conversation. Files the person needs to see (such as a .docx or .xlsx) are saved to the file system, and the agent says where. It asks for feedback in the conversation.
- **Benchmark.** The agent skips it, because without baselines it means nothing, and relies on the person's qualitative feedback.
- **Iterations.** The agent can still keep the results in iteration directories when it has a file system.
- **Description tuning** needs a runner that starts the agent headless, and is not available.
- **Blind comparison** needs subagents and is not available.
- **Packaging** works anywhere with Python and a file system; the person downloads the `.skill` file.

## Without a display

- Subagents are available, so the main workflow works: parallel test runs, baselines, grading. When timeouts become a serious problem, the agent runs the tests one after the other.
- There is no browser or display. The agent creates the viewer with `--static <output_path>` and gives the person a link to the HTML file.
- The "Submit All Reviews" button downloads `feedback.json`. The agent reads it from there, and may first have to ask for access.
- The agent still opens the viewer before it evaluates the outputs itself. The agent tends to skip this step without a display, and the person's early look at the results is what the loop depends on.
- Packaging works, because `package_skill.py` needs only Python and a file system.
- Description tuning works where the harness's adapter provides a runner. The agent saves it for the end, when the skill is finished and the person agrees.

## Updating an installed skill

When the person wants to change an existing skill instead of creating a new one:

- **The name stays.** The agent keeps the directory name and the `name` in the frontmatter. An installed `research-helper` becomes `research-helper.skill`, not `research-helper-v2`.
- **The agent edits a copy.** The installed path may be read-only, so the agent copies the skill to a writable location (its scratch directory), edits it there, and packages from the copy.
- **Packaging by hand** also starts in a writable temporary location, and the agent copies the result to the output directory afterwards, because writing directly there may fail.
