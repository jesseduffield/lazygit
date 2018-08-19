#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init
git config user.email "test@example.com"
git config user.name "Lazygit Tester"


function add_spacing {
  for i in {1..60}
  do
    echo "..." >> $1
  done
}

echo "Here is a story that has been told throuhg the ages" >> file1

git add file1
git commit -m "first commit"

git checkout -b develop

echo "once upon a time there was a dog" >> file1
add_spacing file1
echo "once upon a time there was another dog" >> file1
git add file1
git commit -m "first commit on develop"

git checkout master

echo "once upon a time there was a cat" >> file1
add_spacing file1
echo "once upon a time there was another cat" >> file1
git add file1
git commit -m "first commit on develop"

git merge develop # should have a merge conflict here
