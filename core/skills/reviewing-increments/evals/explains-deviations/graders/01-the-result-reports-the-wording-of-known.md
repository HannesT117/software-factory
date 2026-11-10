---
type: llm
focus: { source: file, path: evidence/01-greeting/review.md }
weight: 1
---

PASS if the section `## Findings` of this review result reports, as a finding of its own, that the line that cmd/greet/main.go prints for `greet --known` is worded differently from the line that the note fixes: the code prints `greet: <name> is known`, the note fixes `greet: known <name>`. The finding may quote the lines with any name or with `<name>`, with or without `greet: `, or describe the difference in its own words.
FAIL if the file is missing or empty, if the result says it has no findings, or if the wording appears only in a remark, under another section, or inside a finding about something else, such as the order of the names.
