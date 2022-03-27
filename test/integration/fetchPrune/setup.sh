#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

# we're setting this to ensure that it's honoured by the fetch command
git config fetch.prune true

echo test1 > myfile1
git add .
git commit -am "myfile1"

git checkout -b other_branch
git checkout master

cd ..
git clone --bare ./repo origin

cd repo

git remote add origin ../origin
git fetch origin
git branch --set-upstream-to=origin/master master
git branch --set-upstream-to=origin/other_branch other_branch

# unbenownst to our test repo we're removing the branch on the remote, so upon
# fetching with prune: true we expect git to realise the remote branch is gone
git -C ../origin branch -d other_branch
