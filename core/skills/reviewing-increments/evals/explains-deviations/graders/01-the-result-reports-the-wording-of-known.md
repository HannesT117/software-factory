---
type: llm
focus: { source: file, path: evidence/01-greeting/review.md }
weight: 1
---

PASS if the section `## Findings` of this review result reports, as a finding of its own, that cmd/greet/main.go prints `greet: <name> is known` where the note fixes `greet: known <name>`.
FAIL if the file is missing or empty, if the result says it has no findings, or if the wording appears only in a remark, under another section, or inside a finding about something else, such as the order of the names.
