#!/usr/bin/env bash
set -eu
here="$(dirname "$0")"
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
cp -R "$here/files/base/." .
mkdir -p .agents/skills
cp -R "$here/files/old/." .
git init -q -b main
git add -A
git commit -q -m "Add release rules"
git branch old-work
cp -R "$here/files/new/." .
git commit -q -a -m "Write release notes per version"
git clone -q --bare . ../origin.git
git remote add origin ../origin.git
git fetch -q origin
git checkout -q old-work
