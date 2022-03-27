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

cd ..
git clone --bare ./repo origin

cd repo

echo test3 > myfile3
git add .
git commit -am "myfile3"
echo test4 > myfile4
git add .
git commit -am "myfile4"

git remote add origin ../origin
git fetch origin
git branch --set-upstream-to=origin/master master

# actually getting a password prompt is tricky: it requires SSH'ing into localhost under a newly created, restricted, user. This is not easy to do in a cross-platform way, nor is it easy to do in a docker container. If you can think of a way to do it, please let me know!
cp ../../../../hooks/pre-push .git/hooks/pre-push
chmod +x .git/hooks/pre-push
