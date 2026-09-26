#!/bin/sh

set -e

# Re-records the demos that a page embeds and points the page at the new
# videos. Use it when a change to demo/settings.tape, or to lazygit's
# appearance, leaves the existing recordings looking out of date.
#
# A page names the demo behind each video in a comment above it:
#
#   <!-- demo: commit_and_push -->
#   <video src="https://github.com/user-attachments/assets/..." controls></video>
#
# GitHub drops that comment when it renders the page, so it costs the reader
# nothing. This script re-records the demo the comment names and rewrites the
# URL on the line below it.

usage() {
    echo "Usage: $0 [--no-upload] [page ...]"
    echo "e.g. $0 README.md"
    echo
    echo "Re-records every demo the given pages embed. With no page given,"
    echo "that means README.md."
    echo
    echo "--no-upload leaves the videos in demo/output and the pages untouched,"
    echo "which is what you want for reviewing a change to demo/settings.tape."
    exit 1
}

NO_UPLOAD=

if [ "$1" = "--no-upload" ]
then
    NO_UPLOAD=--no-upload
    shift
fi

if [ "$1" = "-h" ] || [ "$1" = "--help" ]
then
    usage
fi

if [ ! -x demo/record_demo.sh ]
then
    echo "Run this from the root of the repository."
    exit 1
fi

if [ "$#" -eq 0 ]
then
    set -- README.md
fi

for PAGE in "$@"
do
    if [ ! -f "$PAGE" ]
    then
        echo "$PAGE does not exist"
        exit 1
    fi

    NAMES=$(sed -n 's/^<!-- demo: \([A-Za-z0-9_]*\) -->$/\1/p' "$PAGE")

    if [ -z "$NAMES" ]
    then
        echo "$PAGE embeds no demos. Each video needs a <!-- demo: <name> -->"
        echo "comment on the line above it to say which demo it came from."
        exit 1
    fi

    for NAME in $NAMES
    do
        TEST="pkg/integration/tests/demo/$NAME.go"

        if [ ! -f "$TEST" ]
        then
            echo "$PAGE asks for a demo called $NAME, but $TEST does not exist"
            exit 1
        fi

        echo
        echo "=== $NAME ==="

        if [ -n "$NO_UPLOAD" ]
        then
            demo/record_demo.sh --no-upload "$TEST"
            continue
        fi

        # Keep the recording chatter on screen, since a full run takes a while,
        # and read the new URL back out of it afterwards.
        LOG=$(mktemp)
        demo/record_demo.sh "$TEST" | tee "$LOG"
        URL=$(sed -n 's/.*<video src="\([^"]*\)".*/\1/p' "$LOG")
        rm -f "$LOG"

        if [ -z "$URL" ]
        then
            echo "Recording $NAME produced no URL"
            exit 1
        fi

        awk -v marker="<!-- demo: $NAME -->" -v url="$URL" '
            hit { sub(/src="[^"]*"/, "src=\"" url "\""); hit = 0 }
            $0 == marker { hit = 1 }
            { print }
        ' "$PAGE" > "$PAGE.new"

        mv "$PAGE.new" "$PAGE"

        echo "$PAGE now points at $URL"
    done
done
