#!/usr/bin/env python3
"""
Check a skill against the documents of the repository it describes.

The script compares the skill with the repository as it is on a git ref
(by default the remote's default branch), not with the working tree, because
a working branch can be behind. It reports:

1. Paths that the skill names but that do not exist on the ref (paths with a
   folder, and top-level documents such as AGENTS.md).
2. Text that the skill shares verbatim with the repository's documents
   (a copy that will drift away from its source).
3. Declared copies ("copied from <path>@<commit>") whose source changed since.
4. Headings that differ between a template in the skill and its source,
   for each --template pair.

It cannot judge whether a rule is stated correctly; that needs the trace
table and the person.

Usage (from the authoring-skills directory):
    python -m scripts.check_sources <skill-dir> --repo <repo-root> [--ref origin/main]
        [--ignore <glob>] ... [--template <file-in-skill>[:N]=<path-in-repo>[:N]] ...

":N" selects the N-th fenced code block (counting from 1) of a file, for a
template that sits inside a code block of a larger document.

Exit codes: 0 no findings, 1 findings, 2 the check could not run.
"""

import argparse
import fnmatch
import re
import subprocess
import sys
from collections import defaultdict
from pathlib import Path

# Suffixes of files that count as documents in the repository
DOC_SUFFIXES = {'.md', '.mdx', '.markdown', '.txt', '.rst', '.adoc', '.yaml', '.yml'}
# Larger files are generated or data, not documents someone copies rules from
MAX_DOC_BYTES = 1_000_000
# Shortest line that counts as a copy. Prose lines must be long enough to be
# unlikely by chance; lines in code blocks (regexes, format strings) are short
# but specific, so they have a lower limit.
MIN_PROSE_CHARS = 40
MIN_CODE_CHARS = 12
# Show at most this many examples per finding
MAX_EXAMPLES = 5
# Folders of the skill that the agent does not read when the skill runs
SKIPPED_DIRS = {'evals', '__pycache__', 'node_modules'}

FENCE_RE = re.compile(r'^\s*(```|~~~)')
HEADING_RE = re.compile(r'^(#{1,6})\s+(.*?)\s*#*\s*$')
BACKTICK_RE = re.compile(r'`([^`\s]+)`')
MD_LINK_RE = re.compile(r'\[[^\]]*\]\(([^)\s]+)\)')
DECLARED_COPY_RE = re.compile(r'([\w][\w./-]*\.\w+)@([0-9a-f]{7,40})\b')
PATH_LIKE_RE = re.compile(r'^[\w.-]+(/[\w.-]+)*(\.\w{1,6}|/)$')
# A bare file name without a folder is usually a file created at run time
# (feedback.json), except for the conventional top-level documents of a repository
TOP_LEVEL_DOC_RE = re.compile(r'^[A-Z][A-Z0-9_-]*\.(md|txt|rst)$')
LIST_MARKER_RE = re.compile(r'^(\s*([-*+>]|\d+[.)]|#{1,6})\s+)+')


class CheckError(Exception):
    """The check cannot run; the message says how to fix it."""


def git(repo, *args, check=True):
    result = subprocess.run(['git', '-C', str(repo), *args], capture_output=True, text=True)
    if check and result.returncode != 0:
        raise CheckError(f"git {' '.join(args)} failed: {result.stderr.strip()}")
    return result


def resolve_ref(repo, ref):
    if git(repo, 'rev-parse', '--git-dir', check=False).returncode != 0:
        raise CheckError(f"{repo} is not a git repository. Pass the repository root with --repo.")
    candidates = [ref] if ref else []
    if not ref:
        head = git(repo, 'symbolic-ref', '--short', 'refs/remotes/origin/HEAD', check=False)
        if head.returncode == 0:
            candidates.append(head.stdout.strip())
        candidates += ['origin/main', 'origin/master']
    for candidate in candidates:
        result = git(repo, 'rev-parse', '--verify', '--quiet', f'{candidate}^{{commit}}', check=False)
        if result.returncode == 0:
            return candidate, result.stdout.strip()
    tried = ', '.join(candidates)
    raise CheckError(
        f"None of these refs exists in {repo}: {tried}. "
        "Run 'git fetch origin' first, or pass the default branch with --ref."
    )


def repo_files(repo, commit):
    out = git(repo, 'ls-tree', '-r', '-l', '--full-tree', commit).stdout
    files = {}
    for line in out.splitlines():
        meta, path = line.split('\t', 1)
        parts = meta.split()
        size = int(parts[3]) if parts[3].isdigit() else 0
        files[path] = size
    return files


def read_blobs(repo, commit, paths):
    """Read many files from the ref in one git process."""
    if not paths:
        return {}
    proc = subprocess.run(
        ['git', '-C', str(repo), 'cat-file', '--batch'],
        input=''.join(f'{commit}:{p}\n' for p in paths).encode(),
        capture_output=True,
    )
    if proc.returncode != 0:
        raise CheckError(f"git cat-file failed: {proc.stderr.decode(errors='replace').strip()}")
    data, pos, blobs = proc.stdout, 0, {}
    for path in paths:
        end = data.index(b'\n', pos)
        header = data[pos:end].decode()
        pos = end + 1
        if header.endswith('missing'):
            continue
        size = int(header.split()[2])
        blobs[path] = data[pos:pos + size].decode('utf-8', errors='replace')
        pos += size + 1
    return blobs


def normalize(line):
    line = LIST_MARKER_RE.sub('', line.strip())
    return re.sub(r'\s+', ' ', line).strip()


def classified_lines(text):
    """Yield (line_number, normalized_line, in_code_block)."""
    in_fence = False
    for number, line in enumerate(text.splitlines(), start=1):
        if FENCE_RE.match(line):
            in_fence = not in_fence
            continue
        norm = normalize(line)
        if norm:
            yield number, norm, in_fence


def code_block(text, index):
    """Return the index-th fenced code block (1-based) of text."""
    blocks, current, in_fence = [], [], False
    for line in text.splitlines():
        if FENCE_RE.match(line):
            if in_fence:
                blocks.append('\n'.join(current))
                current = []
            in_fence = not in_fence
            continue
        if in_fence:
            current.append(line)
    if not 1 <= index <= len(blocks):
        raise CheckError(f"asked for code block {index}, but the file has {len(blocks)} fenced code block(s)")
    return blocks[index - 1]


def headings(text, inside_code_too):
    result, in_fence = [], False
    for line in text.splitlines():
        if FENCE_RE.match(line):
            in_fence = not in_fence
            continue
        if in_fence and not inside_code_too:
            continue
        match = HEADING_RE.match(line)
        if match and not line.startswith('#!'):
            result.append((len(match.group(1)), match.group(2)))
    return result


def skill_docs(skill):
    for path in sorted(skill.rglob('*.md')):
        rel = path.relative_to(skill)
        if any(part in SKIPPED_DIRS for part in rel.parts[:-1]):
            continue
        yield path, rel


def check_paths(skill, files, ignore):
    directories = {str(Path(p).parent) for p in files} | {'.'}
    missing = defaultdict(list)
    for path, rel in skill_docs(skill):
        for number, line in enumerate(path.read_text(errors='replace').splitlines(), start=1):
            tokens = BACKTICK_RE.findall(line) + MD_LINK_RE.findall(line)
            for token in tokens:
                token = re.sub(r'(:\d+(-\d+)?|#.*)$', '', token)
                if '://' in token or token.startswith(('-', '/', '~', '$')) or not PATH_LIKE_RE.match(token):
                    continue
                if '/' not in token.rstrip('/') and not TOP_LEVEL_DOC_RE.match(token):
                    continue
                if any(fnmatch.fnmatch(token, pattern) for pattern in ignore):
                    continue
                if (skill / token).exists() or (path.parent / token).exists():
                    continue
                stripped = token.rstrip('/')
                if stripped in files or stripped in directories:
                    continue
                missing[token].append(f"{rel.as_posix()}:{number}")
    return missing


def check_copies(skill, docs):
    index = {}
    for doc_path, text in docs.items():
        for number, norm, _ in classified_lines(text):
            index.setdefault(norm, (doc_path, number))
    copies = defaultdict(list)
    for path, rel in skill_docs(skill):
        for number, norm, in_code in classified_lines(path.read_text(errors='replace')):
            limit = MIN_CODE_CHARS if in_code else MIN_PROSE_CHARS
            if len(norm) >= limit and norm in index:
                source, source_line = index[norm]
                copies[(rel.as_posix(), source)].append((number, source_line, norm))
    return copies


def check_declared_copies(skill, repo, commit):
    findings = []
    for path, rel in skill_docs(skill):
        for number, line in enumerate(path.read_text(errors='replace').splitlines(), start=1):
            for source, copied_at in DECLARED_COPY_RE.findall(line):
                if git(repo, 'cat-file', '-e', f'{copied_at}^{{commit}}', check=False).returncode != 0:
                    findings.append(f"{rel.as_posix()}:{number}: commit {copied_at} of {source} is not in the repository (fetch it, or fix the hash)")
                    continue
                changed = git(repo, 'diff', '--quiet', copied_at, commit, '--', source, check=False).returncode
                if changed == 1:
                    findings.append(f"{rel.as_posix()}:{number}: {source} changed between {copied_at} and the ref; update the copy and its commit")
    return findings


def split_spec(spec):
    match = re.match(r'^(.*?)(?::(\d+))?$', spec)
    return match.group(1), int(match.group(2)) if match.group(2) else None


def check_template(skill, repo, commit, pair):
    if '=' not in pair:
        raise CheckError(f"--template needs the form <file-in-skill>=<path-in-repo>, got '{pair}'")
    left, right = pair.split('=', 1)
    skill_file, skill_block = split_spec(left)
    repo_file, repo_block = split_spec(right)
    skill_path = skill / skill_file
    if not skill_path.exists():
        raise CheckError(f"--template: {skill_file} does not exist in the skill")
    skill_text = skill_path.read_text(errors='replace')
    result = git(repo, 'show', f'{commit}:{repo_file}', check=False)
    if result.returncode != 0:
        return [f"{repo_file} does not exist on the ref"]
    repo_text = result.stdout
    if skill_block:
        skill_text = code_block(skill_text, skill_block)
    if repo_block:
        repo_text = code_block(repo_text, repo_block)
    ours = headings(skill_text, inside_code_too=bool(skill_block))
    theirs = headings(repo_text, inside_code_too=bool(repo_block))
    our_names = [h[1] for h in ours]
    their_names = [h[1] for h in theirs]
    findings = []
    for level, name in theirs:
        if name not in our_names:
            findings.append(f"missing in the skill: {'#' * level} {name}")
    for level, name in ours:
        if name not in their_names:
            findings.append(f"not in the source: {'#' * level} {name}")
    shared_ours = [h for h in ours if h[1] in their_names]
    shared_theirs = [h for h in theirs if h[1] in our_names]
    if [h[1] for h in shared_ours] != [h[1] for h in shared_theirs]:
        findings.append("shared headings are in a different order: "
                        f"skill {[h[1] for h in shared_ours]}, source {[h[1] for h in shared_theirs]}")
    for (our_level, name), (their_level, _) in zip(
            sorted(shared_ours, key=lambda h: h[1]), sorted(shared_theirs, key=lambda h: h[1])):
        if our_level != their_level:
            findings.append(f"different level for '{name}': skill {'#' * our_level}, source {'#' * their_level}")
    return findings


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__.split('\n\n')[0])
    parser.add_argument('skill', help='path to the skill directory')
    parser.add_argument('--repo', required=True, help='root of the git repository the skill describes')
    parser.add_argument('--ref', help='ref to compare with (default: origin/HEAD, then origin/main, then origin/master)')
    parser.add_argument('--ignore', action='append', default=[],
                        help='glob of paths the skill names that are not repository paths (for example "evals/*"); repeatable')
    parser.add_argument('--template', action='append', default=[],
                        help='<file-in-skill>[:N]=<path-in-repo>[:N]; compare the headings of a template with its source')
    args = parser.parse_args(argv)

    skill, repo = Path(args.skill).resolve(), Path(args.repo).resolve()
    if not (skill / 'SKILL.md').exists():
        print(f"Error: {skill} has no SKILL.md. Pass the skill directory.", file=sys.stderr)
        return 2
    try:
        ref, commit = resolve_ref(repo, args.ref)
        files = repo_files(repo, commit)
        doc_paths = [p for p, size in files.items()
                     if Path(p).suffix.lower() in DOC_SUFFIXES and size <= MAX_DOC_BYTES]
        docs = read_blobs(repo, commit, doc_paths)

        print(f"Compared with {ref} ({commit[:12]}), {len(docs)} documents.")
        print("If the remote may have changed, run 'git fetch' and check again.\n")
        found = 0

        missing = check_paths(skill, files, args.ignore)
        if missing:
            found += len(missing)
            print(f"Paths not found on {ref} ({len(missing)}):")
            for token, places in sorted(missing.items()):
                more = f" and {len(places) - 1} more" if len(places) > 1 else ''
                print(f"  - {token}  ({places[0]}{more})")
            print("  Repository paths must exist on the default branch. Paths of files created at run time can be ignored.\n")

        copies = check_copies(skill, docs)
        if copies:
            found += len(copies)
            print(f"Text shared verbatim with repository documents ({len(copies)}):")
            for (skill_file, source), hits in sorted(copies.items()):
                print(f"  - {skill_file} shares {len(hits)} line(s) with {source}:")
                for number, source_line, text in hits[:MAX_EXAMPLES]:
                    print(f"      {skill_file}:{number} = {source}:{source_line}  {text[:80]}")
                if len(hits) > MAX_EXAMPLES:
                    print(f"      … and {len(hits) - MAX_EXAMPLES} more")
            print("  Point to the source instead of copying. Where a copy is unavoidable, mark it 'copied from <path>@<commit>'.\n")

        declared = check_declared_copies(skill, repo, commit)
        if declared:
            found += len(declared)
            print(f"Declared copies that need attention ({len(declared)}):")
            for line in declared:
                print(f"  - {line}")
            print()

        for pair in args.template:
            findings = check_template(skill, repo, commit, pair)
            if findings:
                found += len(findings)
                print(f"Template {pair}:")
                for line in findings:
                    print(f"  - {line}")
                print()
            else:
                print(f"Template {pair}: headings match.\n")
    except CheckError as error:
        print(f"Error: {error}", file=sys.stderr)
        return 2

    print("No findings." if not found else f"{found} finding(s).")
    return 1 if found else 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
