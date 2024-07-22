#!/bin/sh

base_ref=$1

# Determine the base commit
base_commit=$(git merge-base HEAD origin/"$base_ref")

# Check if base_commit is set correctly
if [ -z "$base_commit" ]; then
    echo "Failed to determine base commit."
    exit 1
fi
echo "Base commit: $base_commit"

# Get commits with "fixup!" in the message from base_commit to HEAD
commits=$(git log -i -E --grep '^fixup!' --grep '^squash!' --grep '^amend!' --grep '^[^\n]*WIP' --grep '^[^\n]*DROPME' --format="%h %s" "$base_commit..HEAD")

if [ -z "$commits" ]; then
    echo "No fixup commits found."
    exit 0
else
    echo "Fixup or WIP commits found:"
    echo "$commits"
    exit 1
fi
