#!/bin/sh

set -e

cd $1

export GIT_COMMITTER_DATE="Mon 20 Aug 2018 20:19:19 BST"
export GIT_AUTHOR_DATE="Mon 20 Aug 2018 20:19:19 BST"

git init

git config user.email "CI@example.com"
git config user.name "CI"
# see https://vielmetti.typepad.com/logbook/2022/10/git-security-fixes-lead-to-fatal-transport-file-not-allowed-error-in-ci-systems-cve-2022-39253.html
# NOTE: I don't think this actually works if it's only applied to the repo.
# On CI we set the global setting, but given it's a security concern I don't want
# people to do that for their locals.
git config protocol.file.allow always

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo test2 > myfile2
git add .
git commit -am "myfile2"

cd ..
git clone --bare ./repo other_repo
cd repo
