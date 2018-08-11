#!/bin/bash

# this script makes a repo with lots of commits

# call this command from the test directory:
# ./lots_of_commits.sh; cd testrepo; gg; cd ..

# -e means exit if something fails
# -x means print out simple commands before running them
set -ex

reponame="case_insensitive_checkouts"

rm -rf ${reponame}
mkdir ${reponame}
cd ${reponame}

git init

touch foo
git add foo
git commit -m "init"
git branch -a
git branch test
git branch TEST
git checkout TEST
git checkout TeST
git checkout TesT
git checkout TEsT
git branch -a