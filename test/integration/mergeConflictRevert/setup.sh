#!/bin/sh

set -e

cd $1

git init
git config user.email "CI@example.com"
git config user.name "CI"

git checkout -b master

echo "test" > file1
git add .
git commit -m "test 1"

git checkout -b other

echo "test" > file2
git add .
git commit -m "test 2"

git checkout -b another

echo "test" > file3
git add .
git commit -m "test 3"

git checkout other

echo "test" > file4
git add .
git commit -m "test 4"

git merge another

echo "test" > file5
git add .
git commit -m "test 5"

