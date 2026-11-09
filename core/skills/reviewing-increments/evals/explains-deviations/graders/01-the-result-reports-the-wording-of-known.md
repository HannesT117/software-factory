---
type: tool_used
tool: Write
input_match: '(?=.*"file_path":\s*"[^"]*evidence/01-greeting/review\.md")(?=.*## Findings(?:(?!\\n## ).)*?\\n- \*\*(?=(?:(?!\\n).)*cmd/greet/main\.go)(?=(?:(?!\\n).)*is known)(?=(?:(?!\\n).)*greet: known [<A]))'
---

The review result that the run writes, evidence/01-greeting/review.md, has under `## Findings` a finding line (`- **<kind>**, ...`) that names cmd/greet/main.go, the printed wording "is known" and the wording "greet: known <name>" that the note fixes.
