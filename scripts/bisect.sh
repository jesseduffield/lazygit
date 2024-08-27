#!/bin/sh

# How to use:
# 1) find a commit that is working fine.
# 2) Create an integration test capturing the fact that it works (Don't commit it). See https://github.com/jesseduffield/lazygit/blob/master/pkg/integration/README.md
# 3) checkout the commit that's known to be failing
# 4) run this script supplying the commit hash / tag name that works and the name of the newly created test

# usage: scripts/bisect.sh <ref that's broken> <ref that's working> <integration test name>
# e.g.   scripts/bisect.sh v0.32.1 mergeConflictsResolvedExternally
# It's assumed that the current commit (i.e. HEAD) is broken.

if [[ $# -ne 3 ]] ; then
    echo 'Usage: scripts/bisect.sh <ref that's broken> <ref that's working> <integration test name>'
    exit 1
fi

git bisect start $1 $2
git bisect run sh -c "(go build -o /dev/null || exit 125) && go test ./pkg/gui -run /$3"
git bisect reset
