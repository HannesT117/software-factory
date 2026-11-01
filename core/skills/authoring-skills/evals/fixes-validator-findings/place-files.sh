#!/usr/bin/env bash
set -eu
here="$(dirname "$0")"
export GIT_AUTHOR_NAME=agent GIT_AUTHOR_EMAIL=agent@example.com
export GIT_COMMITTER_NAME=agent GIT_COMMITTER_EMAIL=agent@example.com
cp -R "$here/files/base/." .
git init -q -b main
git add -A
git commit -q -m "Add the release helper"
