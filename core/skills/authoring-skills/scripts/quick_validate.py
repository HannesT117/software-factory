#!/usr/bin/env python3
"""
Validation for skills.

Errors are problems that break loading or upload: missing frontmatter, or a
field outside the rules of the Agent Skills specification. Warnings are
deviations from the skill authoring best practices and from the declarative
third-person style that authoring-skills asks for. Errors make the script exit
with 1; warnings do so only with --strict.

This file is the one source for the rules of skills. `factory validate-skills`
runs it with --codes and --forbid-tool-names.
"""

import re
import sys
from dataclasses import dataclass
from pathlib import Path

import yaml

PREFIX = 'quick_validate'

# Directories whose contents are not packaged as part of the skill, so any
# SKILL.md inside them does not count toward the single-SKILL.md check.
# Mirrors package_skill.py: __pycache__ and node_modules are excluded at any
# depth, evals only at the skill root.
EXCLUDED_DIR_PARTS = {'__pycache__', 'node_modules'}
ROOT_EXCLUDED_DIR_PARTS = {'evals'}

ALLOWED_FIELDS = {'name', 'description', 'license', 'allowed-tools', 'metadata', 'compatibility'}
RESERVED_NAME_WORDS = ('anthropic', 'claude')
TOOL_NAMES_RE = re.compile(r'\b(claude|cursor|codex|copilot|gemini)\b', re.IGNORECASE)
NAME_RE = re.compile(r'^[a-z0-9]+(-[a-z0-9]+)*$')
MAX_NAME = 64
MAX_DESCRIPTION = 1024
MAX_COMPATIBILITY = 500

# The best practices name 500 lines as the limit of the SKILL.md body
MAX_BODY_LINES = 500
# Longer Markdown files need a table of contents near the top
TOC_THRESHOLD_LINES = 100
TOC_SEARCH_LINES = 30
# At most this many line numbers per finding
MAX_LINES_SHOWN = 8

# Prompts that the agent hands to subagents address them as "you", so the
# voice checks skip them.
SUBAGENT_PROMPT_DIRS = {'agents'}
# Not part of the skill that the agent reads
SKIPPED_DIRS = {'evals', '__pycache__', 'node_modules'}

FRONTMATTER_RE = re.compile(r'^---\n(.*?)\n---(\n|$)', re.DOTALL)
FENCE_RE = re.compile(r'^\s*(```|~~~)')
INLINE_CODE_RE = re.compile(r'`[^`]*`')
QUOTED_RE = re.compile(r'"[^"\n]*"|“[^”\n]*”')
SECOND_PERSON_RE = re.compile(r"\b(you|your|yours|yourself|you're|you'll|you've)\b", re.IGNORECASE)
NAMED_HUMAN_RE = re.compile(r'\bthe (user|engineer|developer)s?\b', re.IGNORECASE)
FIRST_PERSON_RE = re.compile(r"\b(I|I'm|I'll|I've|me|my|we|our|us)\b")
MD_LINK_RE = re.compile(r'\[[^\]]*\]\(([^)\s]+)\)')
BACKSLASH_PATH_RE = re.compile(r'\b[\w.-]+\\[\w.-]+\.(md|py|sh|js|ts|json|yaml|yml|txt|html)\b')
ABSOLUTE_PATH_RE = re.compile(r'(?<![\w/])(/home/\w|/Users/\w|/root/|[A-Z]:\\\w)')
MONTHS = 'January|February|March|April|May|June|July|August|September|October|November|December'
DATE_PHRASE_RE = re.compile(
    rf'\b(before|after|until|since|as of|by|starting)\s+((early|mid|late)\s+)?(({MONTHS})\s+)?(19|20)\d\d\b',
    re.IGNORECASE,
)
TOC_HEADING_RE = re.compile(r'^#{1,4}\s+(table of )?contents\b', re.IGNORECASE)
# Other names for the agent that runs the skill. A product name, such as
# Claude Code or a Claude artifact, is not a name for the agent.
AGENT_TERM_RE = re.compile(
    r"\bthe assistant\b|\bthe AI\b|"
    r"\bClaude\b(?!\s+(Code|Haiku|Sonnet|Opus|Fable|desktop|API|app|apps|artifacts?|in Chrome|Agent SDK)\b)(?!\.ai)(?!-)",
)
AGENT_TERM_CASE_RE = re.compile(r'\bthe assistant\b', re.IGNORECASE)


@dataclass
class Finding:
    severity: str  # 'error' or 'warning'
    path: str
    line: int
    code: str
    detail: str

    def where(self):
        return f'{self.path}:{self.line}' if self.line else self.path


def _counts_as_skill_md(rel_path):
    """True if a SKILL.md at rel_path, relative to the skill root, would be packaged."""
    dir_parts = rel_path.parts[:-1]
    if any(part in EXCLUDED_DIR_PARTS for part in dir_parts):
        return False
    if dir_parts and dir_parts[0] in ROOT_EXCLUDED_DIR_PARTS:
        return False
    return True


def _field_line(content, key):
    """The line number of a top-level frontmatter key, or 0."""
    for number, line in enumerate(content.splitlines(), start=1):
        if number > 1 and line == '---':
            return 0
        if line.startswith(f'{key}:'):
            return number
    return 0


def _read(path):
    return path.read_text(errors='replace').replace('\r\n', '\n')


def _lines_note(numbers):
    shown = ', '.join(str(n) for n in numbers[:MAX_LINES_SHOWN])
    more = f' and {len(numbers) - MAX_LINES_SHOWN} more' if len(numbers) > MAX_LINES_SHOWN else ''
    return f'line {shown}{more}' if len(numbers) == 1 else f'lines {shown}{more}'


def check_errors(skill_dir, forbid_tool_names=False):
    """Return the errors of the skill in skill_dir, and its frontmatter when it could be read."""
    skill_path = Path(skill_dir)
    skill_md = skill_path / 'SKILL.md'
    where = str(skill_md)
    errors = []

    def error(code, detail, line=0, path=where):
        errors.append(Finding('error', path, line, code, detail))

    if not skill_md.is_file():
        error('no-frontmatter', 'SKILL.md not found')
        return errors, None

    for extra in sorted(skill_path.rglob('SKILL.md')):
        rel = extra.relative_to(skill_path)
        if rel.as_posix() != 'SKILL.md' and _counts_as_skill_md(rel):
            error('extra-skill-file', 'a skill has exactly one SKILL.md, at its root; '
                  'rename a supporting file, or make it a skill of its own', path=str(extra))

    content = _read(skill_md)
    match = FRONTMATTER_RE.match(content)
    if not match:
        error('no-frontmatter', 'SKILL.md must start with a YAML frontmatter between two lines ---', line=1)
        return errors, None
    try:
        frontmatter = yaml.safe_load(match.group(1))
    except yaml.YAMLError as e:
        error('invalid-yaml', ' '.join(str(e).split()), line=1)
        return errors, None
    if not isinstance(frontmatter, dict):
        error('invalid-yaml', 'the frontmatter must be a YAML map', line=1)
        return errors, None

    for key in sorted(set(frontmatter) - ALLOWED_FIELDS, key=str):
        error('unknown-field', f'{key}; allowed are {", ".join(sorted(ALLOWED_FIELDS))}',
              line=_field_line(content, str(key)))

    name = frontmatter.get('name')
    line = _field_line(content, 'name')
    if not isinstance(name, str) or not name.strip():
        error('name-format', 'name is missing or not a string', line=line)
        name = ''
    else:
        name = name.strip()
        if not NAME_RE.match(name):
            error('name-format', f"'{name}': lowercase letters, digits and single hyphens, "
                  'not at the start or end', line=line)
        if len(name) > MAX_NAME:
            error('name-too-long', f'{len(name)} characters, at most {MAX_NAME}', line=line)
        folder = skill_path.resolve().name
        if name != folder:
            error('name-mismatch', f"'{name}' differs from the folder name '{folder}'", line=line)
        for word in RESERVED_NAME_WORDS:
            if word in name:
                error('reserved-word', f"'{name}' contains '{word}'", line=line)

    description = frontmatter.get('description')
    line = _field_line(content, 'description')
    if not isinstance(description, str) or not description.strip():
        error('description-missing', 'the description says what the skill does and when the agent uses it',
              line=line)
        description = ''
    else:
        description = description.strip()
        if '<' in description or '>' in description:
            error('description-brackets', 'no < or > in the description', line=line)
        if len(description) > MAX_DESCRIPTION:
            error('description-too-long', f'{len(description)} characters, at most {MAX_DESCRIPTION}', line=line)

    if 'compatibility' in frontmatter:
        compatibility = frontmatter['compatibility']
        line = _field_line(content, 'compatibility')
        if not isinstance(compatibility, str) or not compatibility.strip():
            error('compatibility-too-long', 'compatibility is a string of 1 to 500 characters', line=line)
        elif len(compatibility) > MAX_COMPATIBILITY:
            error('compatibility-too-long', f'{len(compatibility)} characters, at most {MAX_COMPATIBILITY}',
                  line=line)

    if 'metadata' in frontmatter and not isinstance(frontmatter['metadata'], dict):
        error('metadata-not-map', 'metadata is a map of keys to values', line=_field_line(content, 'metadata'))

    if forbid_tool_names:
        for key, value in (('name', name), ('description', description)):
            found = TOOL_NAMES_RE.search(value)
            if found:
                error('tool-name', f'{key} names {found.group(0)}', line=_field_line(content, key))
        offset = content.count('\n', 0, match.end())
        for number, line_text in enumerate(content[match.end():].splitlines(), start=offset + 1):
            found = TOOL_NAMES_RE.search(line_text)
            if found:
                error('tool-name', f'the body names {found.group(0)}', line=number)

    return errors, frontmatter


def _skill_files(skill_path, suffix):
    for p in sorted(skill_path.rglob(f'*{suffix}')):
        rel = p.relative_to(skill_path)
        if any(part in SKIPPED_DIRS for part in rel.parts[:-1]):
            continue
        yield p, rel


def _prose_lines(text, first=1):
    """Yield (number, line) outside fenced code blocks, from line first on, without inline code and quotes."""
    in_fence = False
    for number, line in enumerate(text.splitlines(), start=1):
        if FENCE_RE.match(line):
            in_fence = not in_fence
            continue
        if not in_fence and number >= first:
            yield number, QUOTED_RE.sub('', INLINE_CODE_RE.sub('', line))


def _links(path, skill_path, text, first=1):
    """Yield (number, skill-relative path) for each Markdown file of the skill that a link in text names."""
    root = skill_path.resolve()
    for number, line in _prose_lines(text, first):
        for target in MD_LINK_RE.findall(line):
            target = target.split('#')[0]
            if not target.endswith('.md') or re.match(r'^[a-z]+:', target):
                continue
            candidate = (path.parent / target).resolve()
            if candidate.is_file() and root in candidate.parents:
                yield number, candidate.relative_to(root).as_posix()


def check_warnings(skill_dir, frontmatter):
    """Return the warnings of a skill whose frontmatter could be read."""
    skill_path = Path(skill_dir)
    skill_md = skill_path / 'SKILL.md'
    content = _read(skill_md)
    body_start = content.count('\n', 0, FRONTMATTER_RE.match(content).end()) + 1
    warnings = []

    def warn(path, lines, code, detail):
        lines = sorted({n for n in lines if n})
        if len(lines) > 1:
            detail = f'{detail} ({_lines_note(lines)})'
        warnings.append(Finding('warning', str(path), lines[0] if lines else 0, code, detail))

    name = str(frontmatter.get('name') or '').strip()
    first_word = name.split('-')[0]
    if first_word and not first_word.endswith('ing'):
        warn(skill_md, [_field_line(content, 'name')], 'name-not-gerund',
             f"'{name}': the best practices prefer the -ing form, such as processing-pdfs")

    description = QUOTED_RE.sub('', str(frontmatter.get('description') or '').strip())
    found = sorted({m.group(0) for regex in (SECOND_PERSON_RE, FIRST_PERSON_RE, NAMED_HUMAN_RE)
                    for m in regex.finditer(description)})
    if found or re.match(r'use this skill\b', description, re.IGNORECASE):
        said = ', '.join(found) if found else 'Use this skill'
        warn(skill_md, [_field_line(content, 'description')], 'description-voice',
             f'{said}: the description says in the third person what the skill does, and calls the human the person')

    body_lines = content[FRONTMATTER_RE.match(content).end():].count('\n') + 1
    if body_lines > MAX_BODY_LINES:
        warn(skill_md, [], 'body-too-long', f'{body_lines} lines, at most {MAX_BODY_LINES}; '
             'move details into reference files that SKILL.md links')

    linked = {t for _, t in _links(skill_md, skill_path, content, body_start)}
    for path, rel in _skill_files(skill_path, '.md'):
        is_skill_md = rel.as_posix() == 'SKILL.md'
        text = content if is_skill_md else _read(path)
        first = body_start if is_skill_md else 1
        lines = text.splitlines()

        if not is_skill_md:
            nested = [n for n, t in _links(path, skill_path, text) if t not in linked and t != 'SKILL.md']
            if nested:
                warn(path, nested, 'unlinked-reference',
                     'links a Markdown file of the skill that SKILL.md does not link')
            if len(lines) > TOC_THRESHOLD_LINES and not any(TOC_HEADING_RE.match(l) for l in lines[:TOC_SEARCH_LINES]):
                warn(path, [], 'missing-contents', f'{len(lines)} lines and no heading Contents '
                     f'in the first {TOC_SEARCH_LINES} lines')

        prose = list(_prose_lines(text, first))
        if rel.parts[0] not in SUBAGENT_PROMPT_DIRS:
            voice = [n for n, l in prose if SECOND_PERSON_RE.search(l) or NAMED_HUMAN_RE.search(l)]
            if voice:
                warn(path, voice, 'reader-voice', 'says you, or names the user, engineer or developer; '
                     'write the agent does X, and call the human the person')
            terms = [n for n, l in prose if AGENT_TERM_RE.search(l) or AGENT_TERM_CASE_RE.search(l)]
            if terms:
                warn(path, terms, 'agent-term', 'calls the agent the assistant, the AI or Claude; '
                     'the fixed term is the agent')

        for regex, code, detail in (
            (BACKSLASH_PATH_RE, 'path-style', 'a path with backslashes; use forward slashes'),
            (ABSOLUTE_PATH_RE, 'path-style', 'an absolute path of one machine; use a path relative to the skill'),
            (DATE_PHRASE_RE, 'time-phrase', 'a phrase that goes stale; describe the current method'),
        ):
            hits = [n for n, l in prose if regex.search(l)]
            if hits:
                warn(path, hits, code, detail)

    return warnings


def check(skill_dir, forbid_tool_names=False):
    """Return all findings of the skill in skill_dir: errors first, then warnings."""
    errors, frontmatter = check_errors(skill_dir, forbid_tool_names)
    if frontmatter is None:
        return errors
    return errors + check_warnings(skill_dir, frontmatter)


def validate_skill(skill_dir):
    """Return (valid, message), as package_skill.py expects. Warnings do not count."""
    errors, _ = check_errors(skill_dir)
    if not errors:
        return True, 'Skill is valid!'
    return False, '; '.join(f'{e.where()}: {e.code} {e.detail}' for e in errors)


def main(argv):
    import argparse
    parser = argparse.ArgumentParser(description='Validate a skill and report deviations from the best practices.')
    parser.add_argument('skill_directory')
    parser.add_argument('--codes', action='store_true',
                        help='print one grep-able line per finding, with its code')
    parser.add_argument('--forbid-tool-names', action='store_true',
                        help='report the name of an agent tool in name, description or body as an error')
    parser.add_argument('--strict', action='store_true', help='count warnings as errors')
    args = parser.parse_args(argv)

    findings = check(args.skill_directory, args.forbid_tool_names)
    failing = [f for f in findings if f.severity == 'error' or args.strict]

    if args.codes:
        for f in findings:
            severity = 'error' if f in failing else f.severity
            print(f'{PREFIX}: {severity} {f.where()}: {f.code} {f.detail}')
        print(f'{PREFIX}: {len(failing)} errors' if failing else f'{PREFIX}: ok 1 checked')
    else:
        for f in findings:
            print(f'{f.severity}: {f.where()}: {f.detail} [{f.code}]')
        print(f'{len(failing)} errors' if failing else 'Skill is valid.')
    return 1 if failing else 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
