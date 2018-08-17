#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init

echo "deleted1" > deleted1
echo "deleted2" > deleted2
echo "modified1" > modified1
echo "modified2" > modified2
echo "renamed" > renamed1

git add .
git commit -m "files to delete"
rm deleted1
rm deleted2

rm renamed1
echo "renamed" > renamed2
echo "more" >> modified1
echo "more" >> modified2
echo "untracked1" > untracked1
echo "untracked2" > untracked2
echo "blah" > "file with space1"
echo "blah" > "file with space2"
echo "same name as branch" > master

git add deleted1
git add modified1
git add untracked1
git add "file with space2"
git add renamed1
git add renamed2