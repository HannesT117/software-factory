#!/usr/bin/env bash
set -eu
cp -R "$(dirname "$0")/files/." .
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
# A fixed date and no git config of the host give fixed commit hashes, so that
# a grader can name the latest one.
export GIT_AUTHOR_DATE=2026-01-05T10:00:00Z GIT_COMMITTER_DATE=2026-01-05T10:00:00Z
export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_NOSYSTEM=1
git init -q -b main
git add README.md go.mod docs
git commit -q -m "Add the design note for greetings"
git checkout -q -b greeting
git add test
git commit -q -m "Add the acceptance tests for greetings"
git add cmd internal
git commit -q -m "Greet a person by name"
