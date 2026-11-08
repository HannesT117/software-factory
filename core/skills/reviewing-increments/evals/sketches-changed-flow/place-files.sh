#!/usr/bin/env bash
set -eu
cp -R "$(dirname "$0")/files/." .
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
git init -q -b main
git add README.md go.mod docs
git commit -q -m "Add the design note for push"
git checkout -q -b push
git add test
git commit -q -m "Add the acceptance test for push"
git add cmd
git commit -q -m "Push the queued messages"
