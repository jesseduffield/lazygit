#!/bin/bash

# this script will make a repo with a master and develop branch, where we end up
# on the master branch and if we try and merge master we get a merge conflict

# -e means exit if something fails
# -x means print out simple commands before running them
set -ex

reponame="testrepo"

rm -rf ${reponame}
mkdir ${reponame}
cd ${reponame}

git init

echo "Here is a story that has been told throuhg the ages" >> file1
git add file1
git commit -m "first commit"

git checkout -b develop

echo "once upon a time there was a dog" >> file1
git add file1
git commit -m "first commit on develop"

git checkout master

echo "once upon a time there was a cat" >> file1
git add file1
git commit -m "first commit on develop"

git merge develop # should have a merge conflict here
