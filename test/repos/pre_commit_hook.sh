#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init
cp ../extras/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "file" > file
git add file
