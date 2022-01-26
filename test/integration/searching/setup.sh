#!/bin/sh

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
echo test5 > myfile5
