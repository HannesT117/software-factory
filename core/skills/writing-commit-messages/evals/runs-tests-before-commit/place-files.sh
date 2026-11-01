#!/usr/bin/env bash
set -eu
cp -R "$(dirname "$0")/files/base/." .
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
git init -q -b main
git add -A
git commit -q -m "Greet a person by name"
git checkout -q -b languages
cp -R "$(dirname "$0")/files/change/." .
