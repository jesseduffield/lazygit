#!/bin/bash

# this script will make a repo with a master and develop branch, where we end up
# on the master branch and if we try and merge master we get a merge conflict

# call this command from the test directory:
# ./merge_conflict.sh; cd testrepo; gg; cd ..

# -e means exit if something fails
# -x means print out simple commands before running them
set -ex

reponame="merge_conflict"

rm -rf ${reponame}
mkdir ${reponame}
cd ${reponame}

git init

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
