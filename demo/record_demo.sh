#!/bin/sh

set -e

TEST=$1

# The repository the demo is uploaded to. GitHub only plays videos that live in
# its own attachment store, and an attachment is tied to one repository.
REPO=jesseduffield/lazygit

usage() {
    echo "Usage: $0 <test path>"
    echo "e.g. $0 pkg/integration/tests/demo/nuke_working_tree.go"
    exit 1
}

if [ "$#" -ne 1 ]
then
    usage
fi

for TOOL in vhs ttyd ffmpeg gh
do
    if ! command -v "$TOOL" > /dev/null 2>&1
    then
        echo "$TOOL could not be found"
        echo "Install it with: brew install $TOOL"
        exit 1
    fi
done

WORKTREE_PATH=$(git worktree list | grep assets | awk '{print $1}')

if [ -z "$WORKTREE_PATH" ]
then
    echo "Could not find assets worktree. You'll need to create a worktree for the assets branch using the following command:"
    echo "git worktree add .worktrees/assets assets"
    echo "The assets branch has no shared history with the main branch: it exists to store assets which are too large to store in the main branch."
    exit 1
fi

OUTPUT_DIR="$WORKTREE_PATH/demo"

# Get last part of the test path and set that as the output name
# example test path: pkg/integration/tests/01_basic_test.go
# For that we want: NAME=01_basic_test
NAME=$(echo "$TEST" | sed -e 's/.*\///' | sed -e 's/\..*//')

# Add the demo to the tests list (if missing) so that it can be run
go generate pkg/integration/tests/tests.go

mkdir -p "$OUTPUT_DIR"

SCRATCH=$(mktemp -d)
trap 'rm -rf "$SCRATCH"' EXIT

TAPE="$SCRATCH/$NAME.tape"
RECORDING="$SCRATCH/$NAME.mp4"
OUTPUT="$OUTPUT_DIR/$NAME.mp4"

# The two quotes in the marker keep the literal VHSDONE out of the command line
# that stays on screen while we wait for the marker to be printed.
cat > "$TAPE" <<EOF
Output "$RECORDING"

Source demo/settings.tape

# The command is typed while the recording is hidden, so there is nothing to
# gain from animating it.
Set TypingSpeed 0ms

Hide
Type "go run cmd/integration_test/main.go cli --slow $TEST; echo VHS''DONE"
Enter
Wait+Screen@180s /Local branches/
Show
Wait+Screen@600s /VHSDONE/
Hide
EOF

vhs "$TAPE"

if [ ! -f "$RECORDING" ]
then
    echo "vhs recorded the demo but wrote no video."
    echo "vhs 0.12.0 does this; see https://github.com/charmbracelet/vhs/issues/787."
    echo "Install a working version with: go install github.com/charmbracelet/vhs@v0.11.0"
    exit 1
fi

# Hold the last frame for a moment so that the end state stays readable, and
# move the moov atom to the front so that the video starts playing before it
# has fully downloaded.
ffmpeg -y -loglevel error -i "$RECORDING" \
    -vf "tpad=stop_mode=clone:stop_duration=1.2" \
    -c:v libx264 -crf 23 -preset slow -pix_fmt yuv420p \
    -movflags +faststart -an "$OUTPUT"

# GitHub's web editor posts to this endpoint when you drag a file into a
# comment box. It is undocumented, but it accepts an ordinary token, so we can
# upload from here. You need push access to $REPO for it to work. If the
# endpoint ever goes away, drag the video into a comment box on github.com
# instead and copy the URL that GitHub inserts.
REPOSITORY_ID=$(gh api "repos/$REPO" --jq .id)

RESPONSE=$(curl --silent --show-error --fail \
    --request POST \
    --header "Authorization: Bearer $(gh auth token)" \
    --header "Accept: application/json" \
    --header "Content-Type: video/mp4" \
    --data-binary "@$OUTPUT" \
    "https://uploads.github.com/user-attachments/assets?name=$NAME.mp4&content_type=video%2Fmp4&repository_id=$REPOSITORY_ID")

URL=$(echo "$RESPONSE" | sed -e 's/.*"url":"//' -e 's/".*//')

if [ -z "$URL" ]
then
    echo "Could not read an attachment URL out of GitHub's response:"
    echo "$RESPONSE"
    exit 1
fi

echo "Demo recorded to $OUTPUT"
echo
echo "Embed it with:"
echo "<video src=\"$URL\" controls></video>"
