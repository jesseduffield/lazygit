#!/bin/sh

set -e

# Since I couldn't get crowdin-cli to work yet, I'm doing things a bit more
# manually for now. The process is as follows:
#
# 1. Download the translations from Crowdin as a zip file
# 2. Unzip the file
# 3. Run this script with the path to the unzipped directory as an argument
#
# Requires jq (1.7 or later): https://github.com/jqlang/jq

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <download_dir>"
    exit 2
fi

download_dir="$1"

for d in "$download_dir"/*
do
    # We need to remove empty strings from the JSON files; those are the ones
    # that haven't been translated yet. Crowdin has an option to skip these when
    # exporting, but unfortunately it doesn't work for json files.
    jq 'del(..|select(. == ""))' < "$d/en.json" > pkg/i18n/translations/$(basename "$d").json
done
