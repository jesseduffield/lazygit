#!/bin/bash

# Run some tests on the current commit, similar to what CI does; useful for
# checking every commit in a branch with `git rebase -x scripts/check_commit.sh master`.

set -e

git diff --quiet || {
  echo "Error: there are unstaged changes. Please stage or stash them before running this script."
  exit 1
}

just test
just lint
just generate
git diff --quiet || {
  echo "Error: auto-generated files not up to date."
  exit 1
}
