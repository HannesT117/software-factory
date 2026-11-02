# Testing a skill

## Contents

- Test cases: the `evals/` folder of the skill
- Workspace layout for results
- Step 1: start all runs, with the skill and the baseline, in one turn
- Step 2: draft assertions while the runs are in progress
- Step 3: save timing data as runs finish
- Step 4: grade, aggregate, and open the viewer
- What the person sees in the viewer
- Step 5: read the feedback
- Several models

Steps 1 to 5 are one continuous sequence; the agent does not stop partway through. The agent uses this sequence and no other testing skill, such as `/skill-test`.

## Test cases

After the draft, the agent writes two or three realistic test prompts and shows them to the person, for example: "Here are a few test cases I'd like to try. Do these look right, or do you want to add more?" Then it runs them.

Each test case is its own folder under `evals/` in the skill directory. At first a case holds only the prompt; the agent drafts the assertions in Step 2, while the runs are in progress.

```
my-skill/
├── SKILL.md
└── evals/
    └── flags-negative-quantity/   # the case name: lowercase-with-hyphens, says what is tested
        ├── prompt.md              # run settings, then the person's task prompt
        ├── case.yaml              # the name, and what a good result looks like
        ├── graders/               # one file per assertion; a placeholder until Step 2
        └── files/                 # input files, only if the test needs any
```

The format is strict, so before the first case the agent reads the section "Test cases" in `references/schemas.md` and copies its templates. These are plain files that the agent writes and reads itself, so this works everywhere without extra tooling. A harness may also ship a runner for this layout, so that the same tests run from a terminal or in CI.

For a skill that rests on project documents (see "Grounding the skill in its sources" in `SKILL.md`), at least one test checks a domain rule against its source: the prompt makes the agent produce output that the rule governs, and the assertion quotes the rule from the source on the default branch, not from the skill. A test that takes its expected answer from the skill only confirms that the skill agrees with itself.

**When the skill already has `evals/evals.json`** (the older format), its tests still count: the agent runs them together with the case folders as one set. The file belongs to the person, so the agent does not edit, convert or delete it unless the person asks; anything new or changed goes into a case folder. The section "evals.json" in `references/schemas.md` says how the two sets combine.

## Workspace layout

Results go into `<skill-name>-workspace/`, a sibling of the skill directory. Inside it, results are ordered by iteration (`iteration-1/`, `iteration-2/`, …), and inside an iteration each test case has a directory with its case name. The directory names stay the same from one iteration to the next, because the viewer matches a run to its previous output and feedback by that path. The agent creates directories as it needs them, not all at the start.

```
<skill-name>-workspace/
└── iteration-1/
    ├── flags-negative-quantity/  # one directory per test case
    │   ├── eval_metadata.json    # the viewer takes the prompt from here
    │   ├── with_skill/           # a configuration; with one run, also the run directory
    │   │   ├── outputs/          # the files the run produced
    │   │   ├── grading.json
    │   │   └── timing.json
    │   └── without_skill/        # the baseline (or old_skill/), with the same contents
    ├── handles-empty-file/
    ├── benchmark.json
    └── benchmark.md
```

The viewer and the aggregation script both read this layout. A run directory is the one that holds a run's `outputs/`, `grading.json` and `timing.json`. When a configuration runs more than once, each run gets its own directory inside it (`with_skill/run-1/`, `with_skill/run-2/`, …); `eval_metadata.json` stays where it is.

## Step 1: Start all runs in the same turn

For each test case, the agent starts two subagents in the same turn, one with the skill and one without. It does not start the with-skill runs first and come back for the baselines later, so that everything finishes at about the same time.

**With-skill run:**

```
Execute this task:
- Skill path: <path-to-skill>
- Task: <eval prompt>
- Input files: <full paths of the test's input files, or "none">
- Save outputs to: <workspace>/iteration-<N>/<case-name>/with_skill/outputs/
- Outputs to save: <what the person cares about, e.g. "the .docx file", "the final CSV">
```

**Baseline run** (same prompt; the baseline depends on the situation):

- **New skill**: no skill at all. Same prompt, no skill path, outputs to `without_skill/outputs/`.
- **Existing skill**: the old version. Before editing, the agent copies the skill (`cp -r <skill-path> <workspace>/skill-snapshot/`) and points the baseline subagent at the copy. Outputs go to `old_skill/outputs/`.

For each test case, the agent writes an `eval_metadata.json` directly in the case's directory, next to `with_skill/`; the assertions can be empty for now. The viewer and the benchmark never look into the skill's `evals/` folder. This file is where they get the prompt from, and without it the viewer shows "(No prompt found)". The agent writes it again for every test case in every iteration; it does not carry over.

```json
{
  "eval_name": "flags-negative-quantity",
  "prompt": "The person's task prompt",
  "assertions": []
}
```

`eval_name` is the case name, `prompt` is the test's prompt without the settings block, and `assertions` lists its assertions (never the placeholder grader).

## Step 2: Draft assertions while the runs are in progress

The agent uses the waiting time. It drafts quantitative assertions for each test case and explains them to the person. When a test already has assertions (files in its `graders/` folder other than the placeholder, or `expectations` / `assertions` in an old `evals/evals.json`), the agent reviews them and explains what they check.

Good assertions can be verified objectively and have descriptive names, so that a person who glances at the benchmark viewer understands what each one checks. Subjective skills (writing style, design quality) are evaluated qualitatively; the agent does not force assertions onto things that need human judgment.

The agent saves each assertion as its own file in the case's `graders/` folder, deletes the placeholder, and puts the same sentences into the `assertions` list of that test's `eval_metadata.json`. Each assertion is about the result itself ("summary.md states a total revenue of 693.5"), not about where this workflow saves it ("… in the outputs folder"), so that it also holds under another runner.

For a test left to human judgment, the placeholder stays. It exists only so that a runner can load the case, and the agent does not grade it.

The agent also tells the person what the viewer will show: the qualitative outputs and the quantitative benchmark.

## Step 3: Save timing data as runs finish

When a subagent task finishes, its notification contains `total_tokens` and `duration_ms`. The agent saves them at once to `timing.json` in the run directory:

```json
{
  "total_tokens": 84852,
  "duration_ms": 23332,
  "total_duration_seconds": 23.3
}
```

This may be the only chance to capture this data, because a harness may report it only once, when the run ends. The agent handles each finished run as it arrives instead of collecting them.

## Step 4: Grade, aggregate, and open the viewer

When all runs are done:

1. **Grade each run.** A grader subagent (or the agent itself) reads `agents/grader.md` and evaluates each assertion against the outputs. The result goes to `grading.json` in each run directory. Its `expectations` array uses the fields `text`, `passed` and `evidence` (not `name`/`met`/`details` or other variants), because the viewer depends on these exact names. Assertions that a program can check are checked by a script, not by eye: a script is faster, more reliable, and reusable in later iterations.

2. **Aggregate the benchmark.** The agent runs the aggregation script from this skill's directory:

   ```bash
   python -m scripts.aggregate_benchmark <workspace>/iteration-N --skill-name <name>
   ```

   It writes `benchmark.json` and `benchmark.md` with pass rate, time and tokens per configuration, with mean ± standard deviation and the difference. When the agent writes `benchmark.json` by hand, `references/schemas.md` has the exact schema the viewer expects. Each with-skill configuration comes before its baseline.

3. **Analyst pass.** The agent reads the benchmark data and looks for patterns that the averages hide; the section "Analyzing Benchmark Results" in `agents/analyzer.md` lists what to look for: assertions that pass with and without the skill (they do not discriminate), evals with high variance (possibly flaky), and trade-offs between time, tokens and quality. The agent adds its notes to the `notes` list in `benchmark.json`, after the notes the script already wrote there.

4. **Open the viewer** with the qualitative outputs and the quantitative data:

   ```bash
   nohup python <authoring-skills-path>/eval-viewer/generate_review.py \
     <workspace>/iteration-N \
     --skill-name "my-skill" \
     --benchmark <workspace>/iteration-N/benchmark.json \
     > /dev/null 2>&1 &
   VIEWER_PID=$!
   ```

   From iteration 2 on, the agent also passes `--previous-workspace <workspace>/iteration-<N-1>`.

   Without a display (a remote container, for example), `--static <output_path>` writes a standalone HTML file instead of starting a server. The "Submit All Reviews" button then downloads `feedback.json`, which the agent copies into the workspace for the next iteration.

   The agent always uses `generate_review.py` for the viewer and does not write its own HTML. It opens the viewer **before** it evaluates the outputs itself, because getting the results in front of the person early is the point of the loop.

5. **Tell the person**, for example: "I've opened the results in your browser. There are two tabs: 'Outputs' lets you click through each test case and leave feedback, 'Benchmark' shows the quantitative comparison. When you're done, come back here and let me know."

## What the person sees in the viewer

The "Outputs" tab shows one test case at a time:

- **Prompt**: the task that was given
- **Output**: the files the skill produced, rendered inline where possible
- **Previous Output** (iteration 2 and later): a collapsed section with the last iteration's output
- **Formal Grades** (if grading ran): a collapsed section with pass or fail per assertion
- **Feedback**: a text box that saves automatically while the person types
- **Previous Feedback** (iteration 2 and later): the person's comments from last time, below the text box

The "Benchmark" tab shows the summary: pass rates, time and token use per configuration, with details per eval and the analyst's notes.

The person moves with the previous and next buttons or the arrow keys. "Submit All Reviews" saves all feedback to `feedback.json`.

## Step 5: Read the feedback

When the person says they are done, the agent reads `feedback.json`:

```json
{
  "reviews": [
    {"run_id": "flags-negative-quantity-with_skill", "feedback": "the chart is missing axis labels", "timestamp": "..."},
    {"run_id": "handles-empty-file-with_skill", "feedback": "", "timestamp": "..."},
    {"run_id": "csv-summary-2-with_skill", "feedback": "perfect, love this", "timestamp": "..."}
  ],
  "status": "complete"
}
```

Empty feedback means the person found the output fine. The agent focuses its improvements on the test cases with specific complaints, and then stops the viewer server:

```bash
kill $VIEWER_PID 2>/dev/null
```

## Several models

When the skill will run on more than one model (for example a small, a medium and a large model), the agent runs the with-skill configuration once per model, in its own configuration directory (`with_skill-small/`, `with_skill-large/`, …), with the model set on the subagent. The questions differ per model: whether the small model gets enough guidance, whether the medium one finds the instructions clear and efficient, whether the large one is slowed down by explanations it does not need. When subagents cannot choose a model in the environment, the agent tells the person which models remain untested.
