"""Checks what the graders of two eval cases of designing-increments cannot.

The `tool_used` graders of the cases check the subagent's prompt, but they
cannot tell whose Read a call is, and the judge of `claude plugin eval` sees
only the first and last 12 messages of a trace. This script reads the full
traces: it checks that the subagent itself reads the whole note, that its
report holds only questions, and that the design agent shows the subagent's
questions and explanation; it prints both when the words differ. Run the eval
with `--keep-temp`, then pass the `result.json` files:

    python3 check-traces.py <result.json> [<result.json> ...]

It prints one line per run and check. The skill asks the design agent to show
all questions and to show the explanation unchanged; the script reports a
question with other words as FAIL and prints both, for a person to compare.
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
    """Returns the tool calls, each with the id of the subagent call that made
    it (None for the design agent), the tool results by call id, and the final
    answer."""
    calls, results, final = [], {}, ""
    for line in open(path):
        try:
            event = json.loads(line)
        except ValueError:
            continue
        message = event.get("message") or {}
        content = message.get("content")
        if event.get("type") == "assistant" and isinstance(content, list):
            parent = event.get("parent_tool_use_id")
            calls += [dict(b, parent=parent) for b in content if b.get("type") == "tool_use"]
        if event.get("type") == "user" and isinstance(content, list):
            for b in content:
                if b.get("type") == "tool_result":
                    results[b["tool_use_id"]] = text_of(b.get("content"))
        if event.get("type") == "result":
            final = event.get("result") or ""
    return calls, results, final


def report(text):
    """The subagent's own words: the harness puts a frame before them and the
    usage after them, and indents each line."""
    if "The report follows:" in text:
        text = text.split("The report follows:", 1)[1]
    text = text.split("agentId:", 1)[0]
    return "\n".join(line.strip() for line in text.splitlines())


def questions(text):
    return [q.strip() for q in re.findall(r"^\s*\d\.\s+(.*)$", text, re.M)]


def check(case, run_no, trace):
    calls, results, final = read_trace(trace)
    agents = [c for c in calls if c["name"] in ("Agent", "Task") and c["parent"] is None]
    say = lambda name, ok: print(f"{case} run {run_no}: {name}: {'pass' if ok else 'FAIL'}")
    if not agents:
        say("a subagent is started", False)
        return
    agent = agents[-1]
    prompt = agent["input"].get("prompt", "")
    result = report(results.get(agent["id"], ""))
    say("the subagent prompt names the note", NOTE in prompt)
    say("the subagent prompt holds nothing of the conversation",
        not any(w in prompt.lower() for w in CONVERSATION[case]))
    if case == "asks-comprehension-questions":
        # A whole read is a Read without offset or limit, or a `cat` or
        # `git show` of the note whose output goes to no other command.
        mine = [c for c in calls if c["parent"] == agent["id"]]
        reads = [c for c in mine if c["name"] == "Read" and c["input"].get("file_path", "").endswith(NOTE)]
        shells = [c for c in mine if c["name"] == "Bash" and NOTE in c["input"].get("command", "")
                  and re.search(r"\b(cat|git\s+show)\b", c["input"]["command"])]
        whole = [r for r in reads if "offset" not in r["input"] and "limit" not in r["input"]]
        whole += [b for b in shells if "|" not in b["input"]["command"]]
        say("the subagent reads the note, each time whole",
            bool(whole) and len(whole) == len(reads) + len(shells))
        asked, shown = questions(result), questions(final)
        rest = [l for l in result.splitlines() if l.strip() and not re.match(r"\s*\d\.\s", l)]
        say("the subagent's report holds only questions", bool(asked) and not rest
            and all("?" in q for q in asked))
        if rest:
            print("  other lines:", *rest, sep="\n    ")
        say("2 or 3 questions are shown", 2 <= len(shown) <= 3 and len(shown) == len(asked))
        say("the shown questions are the subagent's, word for word", asked == shown)
        if asked != shown:
            print("  subagent:", *asked, sep="\n    ")
            print("  shown:", *shown, sep="\n    ")
    else:
        say("the subagent prompt holds the answers", all(a in prompt for a in ANSWERS))
        paragraphs = [p.strip() for p in result.split("\n\n") if len(p.strip()) > 80]
        flat = lambda t: " ".join(t.split())
        missing = [p for p in paragraphs if flat(p) not in flat(final)]
        say("the explanation is shown unchanged", bool(paragraphs) and not missing)
        for p in missing:
            print("  changed:", p[:200])


for path in sys.argv[1:]:
    for case in json.load(open(path))["cases"]:
        if case["name"] not in CONVERSATION:
            continue
        for arm in case["arms"].values():
            for n, run in enumerate(arm):
                check(case["name"], n, run["tracePath"])
