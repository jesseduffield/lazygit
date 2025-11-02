#!/bin/sh

set -euo pipefail

st=$(git status --porcelain)
if [ -n "$st" ]; then
    echo "Working directory is not clean; aborting."
    exit 1
fi

if diff -r -q docs docs-master > /dev/null && diff -r -q schema schema-master > /dev/null; then
    echo "No changes to docs or schema; nothing to do."
    exit 0
fi

branch_name=update-docs-for-release

if git show-ref --verify --quiet refs/heads/"$branch_name"; then
    echo "Branch '$branch_name' already exists; aborting."
    exit 1
fi

git checkout -b "$branch_name" --no-track origin/master

git rm -r docs schema
cp -r docs-master docs
cp -r schema-master schema
git add docs schema
git commit -m "Update docs and schema for release"

git push -u origin "$branch_name"
