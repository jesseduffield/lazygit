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

git remote add origin ../origin
git fetch origin
git branch --set-upstream-to=origin/master master
