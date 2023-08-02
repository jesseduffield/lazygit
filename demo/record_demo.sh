#!/bin/sh

set -e

TYPE=$1
TEST=$2

usage() {
    echo "Usage: $0 [gif|mp4] <test path>"
    echo "e.g. using full path: $0 gif pkg/integration/tests/demo/nuke_working_tree.go"
    exit 1
}

if [ "$#" -ne 2 ]
then
    usage
fi

if [ "$TYPE" != "gif" ] && [ "$TYPE" != "mp4" ]
then
    usage
    exit 1
fi

if [ -z "$TEST" ]
then
    usage
fi

WORKTREE_PATH=$(git worktree list | grep assets | awk '{print $1}')

if [ -z "$WORKTREE_PATH" ]
then
    echo "Could not find assets worktree. You'll need to create a worktree for the assets branch using the following command:"
    echo "git worktree add .worktrees/assets assets"
    echo "The assets branch has no shared history with the main branch: it exists to store assets which are too large to store in the main branch."
    exit 1
fi

OUTPUT_DIR="$WORKTREE_PATH/demo"

if ! command -v terminalizer &> /dev/null
then
    echo "terminalizer could not be found"
    echo "Install it with: npm install -g terminalizer"
    exit 1
fi

if ! command -v "gifsicle" &> /dev/null
then
    echo "gifsicle could not be found"
    echo "Install it with: npm install -g gifsicle"
    exit 1
fi

# Get last part of the test path and set that as the output name
# example test path: pkg/integration/tests/01_basic_test.go
# For that we want: NAME=01_basic_test
NAME=$(echo "$TEST" | sed -e 's/.*\///' | sed -e 's/\..*//')

# Add the demo to the tests list (if missing) so that it can be run
go generate pkg/integration/tests/tests.go

mkdir -p "$OUTPUT_DIR"

# First we record the demo into a yaml representation
terminalizer -c demo/config.yml record --skip-sharing -d "go run cmd/integration_test/main.go cli --slow $TEST" "$OUTPUT_DIR/$NAME"
# Then we render it into a gif
terminalizer render "$OUTPUT_DIR/$NAME" -o "$OUTPUT_DIR/$NAME.gif"

# Then we convert it to either an mp4 or gif based on the command line argument
if [ "$TYPE" = "mp4" ]
then
    COMPRESSED_PATH="$OUTPUT_DIR/$NAME.mp4"
    ffmpeg -y -i "$OUTPUT_DIR/$NAME.gif" -movflags faststart -pix_fmt yuv420p -vf "scale=trunc(iw/2)*2:trunc(ih/2)*2" "$COMPRESSED_PATH"
else
    COMPRESSED_PATH="$OUTPUT_DIR/$NAME-compressed.gif"
    gifsicle --colors 256 --use-col=web -O3 < "$OUTPUT_DIR/$NAME.gif" > "$COMPRESSED_PATH"
fi

echo "Demo recorded to $COMPRESSED_PATH"
