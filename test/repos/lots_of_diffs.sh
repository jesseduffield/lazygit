#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init
git config user.email "test@example.com"
git config user.name "Lazygit Tester"


echo "deleted" > deleted_staged
echo "deleted_unstaged" > deleted_unstaged
echo "modified_staged" > modified_staged
echo "modified_unstaged" > modified_unstaged
echo "renamed" > renamed_before

git add .
git commit -m "files to delete"
rm deleted_staged
rm deleted_unstaged

rm renamed_before
echo "renamed" > renamed_after
echo "more" >> modified_staged
echo "more" >> modified_unstaged
echo "untracked_staged" > untracked_staged
echo "untracked_unstaged" > untracked_unstaged
echo "blah" > "file with space staged"
echo "blah" > "file with space unstaged"
echo "same name as branch" > master

git add deleted_staged
git add modified_staged
git add untracked_staged
git add "file with space staged"
git add renamed_before
git add renamed_after