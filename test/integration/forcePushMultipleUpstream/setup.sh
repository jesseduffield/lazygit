#!/bin/sh

set -e

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"
git config push.default upstream

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo test2 > myfile2
git add .
git commit -am "myfile2"

git checkout -b other_branch
git checkout master

cd ..
git clone --bare ./repo origin

cd repo

git remote add origin ../origin
git fetch origin
git branch --set-upstream-to=origin/master master
git branch --set-upstream-to=origin/other_branch other_branch

echo test3 > myfile3
git add .
git commit -am "myfile3"

git push origin master
git reset --hard HEAD^

git checkout other_branch

echo test4 > myfile4
git add .
git commit -am "myfile4"

git push origin other_branch
git reset --hard HEAD^

git checkout master

# at this point, both branches have diverged from their remote counterparts, meaning if you
# attempt to push either, it'll ask if you want to force push.
