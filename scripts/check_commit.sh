#!/bin/bash

# Run some tests on the current commit, similar to what CI does; useful for
# checking every commit in a branch with `git rebase -x scripts/check_commit.sh master`.

set -e

just test
just lint

status_before_generate=$(git status --porcelain=v1)
just generate
status_after_generate=$(git status --porcelain=v1)
if [[ "$status_after_generate" != "$status_before_generate" ]]; then
  echo "Error: auto-generated files not up to date."
  exit 1
fi
