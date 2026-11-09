#!/usr/bin/env bash
set -eu
cp -R "$(dirname "$0")/files/." .
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
git init -q -b main
git add README.md go.mod docs/domain docs/increments/01-greeting evidence
git commit -q -m "Greet a person by name"
git checkout -q -b design-shout
git add docs/increments/02-shout
git commit -q -m "Plan a shouted greeting"
