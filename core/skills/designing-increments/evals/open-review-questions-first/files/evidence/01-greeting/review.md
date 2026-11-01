# Review of increment 01: greeting

Reviewed commit: 3f1c2a9b8d7e6f5a4b3c2d1e0f9a8b7c6d5e4f3a
Result: no findings

## Findings

None.

## Comprehension questions

- **question 1** Why does `greet` print its error to stdout and not to stderr?
  - **answer** All output goes to stdout, so that one pattern in a test matches every line.
- **question 2** A user reports that `greet` prints nothing at all. Where do you look first, and what do you run?
