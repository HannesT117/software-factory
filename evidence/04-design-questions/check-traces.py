"""Checks the graders on subagent calls of two eval cases of designing-increments.

The judge of `claude plugin eval` sees a shortened trace without the subagent
calls. This script reads the full traces instead. Run the eval with
`--keep-temp`, then pass the `result.json` files:

    python3 check-traces.py <result.json> [<result.json> ...]

It prints one line per run and check. Grader 02 of asks-comprehension-questions
asks for the same content, not the same words, so the script prints both sets
of questions when the words differ, for a person to compare.
"""

import json
import re
import sys

NOTE = "docs/increments/02-shout/design.md"
CONVERSATION = {
    "asks-comprehension-questions": ["--loud", "birgit"],
    "explains-gaps-after-answers": ["--loud", "german"],
}
ANSWERS = ["exits with 1", "I don't know", "HELLO, ADA!"]


def text_of(content):
    if isinstance(content, str):
        return content
    return " ".join(part.get("text", "") for part in content or [])


def read_trace(path):
    calls, results, final = [], {}, ""
    for line in open(path):
        try:
            event = json.loads(line)
        except ValueError:
            continue
        message = event.get("message") or {}
        content = message.get("content")
        if event.get("type") == "assistant" and isinstance(content, list):
            calls += [b for b in content if b.get("type") == "tool_use"]
        if event.get("type") == "user" and isinstance(content, list):
            for b in content:
                if b.get("type") == "tool_result":
                    results[b["tool_use_id"]] = text_of(b.get("content"))
        if event.get("type") == "result":
            final = event.get("result") or ""
    return calls, results, final


def questions(text):
    return [q.strip() for q in re.findall(r"^\s*\d\.\s+(.*)$", text, re.M)]


def check(case, run_no, trace):
    calls, results, final = read_trace(trace)
    agents = [c for c in calls if c["name"] in ("Agent", "Task")]
    say = lambda name, ok: print(f"{case} run {run_no}: {name}: {'pass' if ok else 'FAIL'}")
    if not agents:
        say("a subagent is started", False)
        return
    agent = agents[-1]
    prompt = agent["input"].get("prompt", "")
    result = results.get(agent["id"], "")
    say("the subagent prompt names the note", NOTE in prompt)
    say("the subagent prompt holds nothing of the conversation",
        not any(w in prompt.lower() for w in CONVERSATION[case]))
    if case == "asks-comprehension-questions":
        reads = [c for c in calls if c["name"] == "Read" and c["input"].get("file_path", "").endswith(NOTE)]
        say("the note is read whole", any("offset" not in r["input"] and "limit" not in r["input"] for r in reads))
        asked, shown = questions(result), questions(final)
        say("2 or 3 questions are shown", 2 <= len(shown) <= 3 and len(shown) == len(asked))
        if asked != shown:
            print("  subagent:", *asked, sep="\n    ")
            print("  shown:", *shown, sep="\n    ")
    else:
        say("the subagent prompt holds the answers", all(a in prompt for a in ANSWERS))
        paragraphs = [p.strip() for p in result.split("\n\n") if len(p.strip()) > 80]
        say("the explanation is shown", bool(paragraphs) and all(p[:60] in final for p in paragraphs))


for path in sys.argv[1:]:
    for case in json.load(open(path))["cases"]:
        if case["name"] not in CONVERSATION:
            continue
        for arm in case["arms"].values():
            for n, run in enumerate(arm):
                check(case["name"], n, run["tracePath"])
