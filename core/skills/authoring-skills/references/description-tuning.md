# Description tuning

## Contents

- How skill triggering works
- Step 1: write trigger eval queries
- Step 2: review the queries with the person
- Step 3: run the optimization loop
- Step 4: apply the result

The `description` in the frontmatter is the main thing that decides whether the agent loads a skill. After the skill is stable and the person agrees, the agent offers to tune it. The loop needs a runner that starts the agent headless with the skill installed and records whether the agent loads it. Such a runner depends on the harness: the factory's adapter for a harness provides it, in `adapters/<harness>/authoring-skills/`, with the command in its README. Without a runner, tuning is not available, and the agent says so.

## How skill triggering works

Skills appear in the agent's list of available skills with their name and description, and the agent decides from these alone whether to read a skill. The agent consults skills only for tasks it cannot easily do by itself. A simple one-step request such as "read this PDF" may not trigger a skill even when the description matches perfectly, because the agent can do it directly with its basic tools. Complex, multi-step or specialized requests trigger skills reliably when the description matches. Trigger queries therefore have to be substantial enough that the agent would gain something from the skill.

Tuning makes a skill more likely to load; it cannot make loading certain. When a skill has to load every time, a one-line pointer in the project's agent instruction file, such as `AGENTS.md`, does that, and the agent offers it to the person.

## Step 1: Write trigger eval queries

The agent writes 20 queries, some that should trigger the skill and some that should not, and saves them as JSON:

```json
[
  {"query": "the person's prompt", "should_trigger": true},
  {"query": "another prompt", "should_trigger": false}
]
```

The queries are realistic: something a person would actually type to the agent. They are concrete and detailed, with file paths, context about the person's job or situation, column names and values, company names, URLs and a little backstory. Some are in lowercase, with abbreviations, typos or casual speech, and their lengths vary. They focus on edge cases rather than clear cases; the person reviews them in Step 2.

Weak: `"Format this data"`, `"Extract text from PDF"`, `"Create a chart"`

Strong: `"ok so my boss just sent me this xlsx file (its in my downloads, called something like 'Q4 sales final FINAL v2.xlsx') and she wants me to add a column that shows the profit margin as a percentage. The revenue is in column C and costs are in column D i think"`

**Should trigger (8 to 10 queries).** Different wordings of the same intent, some formal and some casual. Cases where the person does not name the skill or the file type but clearly needs it. Some uncommon uses, and cases where this skill competes with another one and should win.

**Should not trigger (8 to 10 queries).** The useful ones are near misses: queries that share words or ideas with the skill but need something else. Neighboring domains, unclear wording where a plain keyword match would trigger wrongly, and cases where the query touches what the skill does but another tool fits better. A query that is obviously unrelated ("write a fibonacci function" for a PDF skill) tests nothing.

## Step 2: Review the queries with the person

Bad queries lead to bad descriptions, so the person reviews them in an HTML page:

1. The agent reads the template `assets/eval_review.html`.
2. It replaces the placeholders:
   - `__EVAL_DATA_PLACEHOLDER__` with the JSON array of queries (without quotes around it; it is a JavaScript variable assignment)
   - `__SKILL_NAME_PLACEHOLDER__` with the skill's name
   - `__SKILL_DESCRIPTION_PLACEHOLDER__` with the current description
3. It writes the result to a temporary file (for example `eval_review_<skill-name>.html` in its scratch directory) and opens it, or sends it to the person where it cannot open a browser.
4. The person edits queries, switches `should_trigger`, adds or removes entries, and clicks "Export Eval Set".
5. The file is downloaded as `eval_set.json`, usually into the person's Downloads folder. When there are several (`eval_set (1).json`), the newest one counts.

## Step 3: Run the optimization loop

The agent tells the person that this takes a while and that it runs in the background. It saves the eval set in the workspace and starts:

```bash
python -m scripts.run_loop \
  --eval-set <path-to-trigger-eval.json> \
  --skill-path <path-to-skill> \
  --model <model-id-powering-this-session> \
  --max-iterations 5 \
  --verbose
```

The model is the one powering the current session, so that the test matches what the person will experience. While the loop runs, the agent checks the output now and then and tells the person which iteration it is on and how the scores look.

The loop splits the eval set into 60% for training and 40% held back for testing. It evaluates the current description (each query three times, for a reliable trigger rate), asks a model for an improvement based on the failures, and evaluates each new description on both sets, for up to five iterations. At the end it opens an HTML report per iteration and returns JSON with `best_description`, chosen by the test score rather than the training score to avoid overfitting.

The improver writes descriptions in the third person, so the result already follows the rule in `SKILL.md`. The agent still checks that the chosen description is in the third person and under 1,024 characters before applying it.

## Step 4: Apply the result

The agent takes `best_description` from the JSON output, writes it into the frontmatter of `SKILL.md`, shows the person the old and the new description, and reports the scores. Then it runs the validator again.
