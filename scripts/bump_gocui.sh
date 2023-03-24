#!/bin/sh

# Go's proxy servers are not very up-to-date so that's why we use `GOPROXY=direct`
# We specify the `master` branch to avoid the default behaviour of looking for a semver tag.
GOPROXY=direct go get -u github.com/jesseduffield/gocui@master && go mod vendor && go mod tidy

# Note to self if you ever want to fork a repo be sure to use this same approach: it's important to use the branch name (e.g. master)
