#!/bin/sh

set -e

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo test2 > myfile2
git add .
git commit -am "myfile2"
echo test3 > myfile3
git add .
git commit -am "myfile3"
echo test4 > myfile4
git add .
git commit -am "myfile4"

cd ..
git clone --bare ./repo origin

cd repo

git reset --hard HEAD~2

echo conflict > myfile4
git add .
git commit -am "myfile4 conflict"

echo test > myfile5
git add .
git commit -am "5"

echo test > myfile6
git add .
git commit -am "6"

echo test > myfile7
git add .
git commit -am "7"

git remote add origin ../origin
git fetch origin
git branch --set-upstream-to=origin/master master

git config pull.rebase interactive
