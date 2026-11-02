# Skill tests in Claude Code

The skill `authoring-skills` in `core/skills/` names no harness. Two parts of its test loop depend on Claude Code and live here.

## Description tuning

`scripts/run_loop.py` installs the skill for `claude -p`, sends each trigger query several times, records whether Claude loads the skill, and asks a model for a better description. `references/description-tuning.md` in the skill describes the loop. From this folder:

```bash
python -m scripts.run_loop \
  --eval-set <path-to-trigger-queries.json> \
  --skill-path <path-to-skill> \
  --model <model-id-of-the-current-session> \
  --max-iterations 5 \
  --verbose
```

`run_eval.py`, `improve_description.py`, `generate_report.py` and `utils.py` serve `run_loop.py`. The scripts need Python 3.10 or later and the `claude` command.

## Test cases from the command line

Claude Code 2.1.269 or later runs the case folders of a skill with and without the skill, from the skill's folder: `claude plugin eval . --allow-tools Write,Edit,Bash`, plus `--scaffold` when a case has input files. It does not read `evals.json` and writes its report under `evals/results/`. These are real, billed model runs, so the agent offers them to the person instead of starting them unasked.

## License

The scripts derive from Anthropic's `skill-creator` under the Apache License 2.0: `LICENSE.txt`, and the changes in `NOTICE.txt`.
