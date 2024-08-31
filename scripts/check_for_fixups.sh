#!/bin/sh

# We will have only done a shallow clone, so the git log will consist only of
# commits on the current PR
commits=$(git log --grep='^fixup!' --grep='^squash!' --grep='^amend!' --grep='^[^\n]*WIP' --grep='^[^\n]*DROPME' --format="%h %s")

if [ -z "$commits" ]; then
    echo "No fixup commits found."
    exit 0
else
    echo "Fixup or WIP commits found:"
    echo "$commits"
    exit 1
fi
