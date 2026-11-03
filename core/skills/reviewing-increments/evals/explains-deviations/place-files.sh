#!/usr/bin/env bash
set -eu
cp -R "$(dirname "$0")/files/." .
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
git init -q -b main
git add README.md go.mod docs names.txt
git commit -q -m "Add the design note for greetings"
git checkout -q -b greeting
git add test
git commit -q -m "Add the acceptance tests for greetings"
git add cmd internal
git commit -q -m "Greet a person by name"
