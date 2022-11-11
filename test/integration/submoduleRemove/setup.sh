#!/bin/sh

set -e

cd $1

export GIT_COMMITTER_DATE="Mon 20 Aug 2018 20:19:19 BST"
export GIT_AUTHOR_DATE="Mon 20 Aug 2018 20:19:19 BST"

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo test2 > myfile2
git add .
git commit -am "myfile2"

cd ..
git clone --bare ./repo other_repo
cd repo

git -c protocol.file.allow=always submodule add ../other_repo
git commit -am "add submodule"
