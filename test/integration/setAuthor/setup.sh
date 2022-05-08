#!/bin/sh

set -e

cd $1

git init

git config user.email "Author1@example.com"
git config user.name "Author1"

echo test1 > myfile1
git add .
git commit -am "myfile1"
echo test2 > myfile2
git add .
git commit -am "myfile2"

git config user.email "Author2@example.com"
git config user.name "Author2"

echo test3 > myfile3
git add .
git commit -am "myfile3"
