#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

git commit --allow-empty -m "Initial commit"

echo test1 > .gitignore

