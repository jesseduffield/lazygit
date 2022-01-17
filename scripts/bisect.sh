#!/bin/sh

# How to use:
# 1) find a commit that is working fine.
# 2) Create an integration test capturing the fact that it works (Don't commit it). See https://github.com/jesseduffield/lazygit/blob/master/docs/Integration_Tests.md
# 3) checkout the commit that's known to be failing
# 4) run this script supplying the commit sha / tag name that works and the name of the newly created test

# usage: scripts/bisect.sh <ref that works> <integration test name>
# e.g.   scripts/bisect.sh v0.32.1 mergeConflictsResolvedExternally
# It's assumed that the current commit (i.e. HEAD) is broken.

set -o pipefail

echo $1
echo $2

git bisect start HEAD $1
git bisect run go test ./pkg/gui -run /$2
git bisect reset
