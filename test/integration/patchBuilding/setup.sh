#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo firstline > myfile2
echo secondline >> myfile2
echo thirdline >> myfile2
git add .
git commit -am "myfile2"
echo firstline2 > myfile2
echo secondline >> myfile2
echo thirdline2 >> myfile2
git commit -am "myfile2 update"
echo test3 > myfile3
git add .
git commit -am "myfile3"
