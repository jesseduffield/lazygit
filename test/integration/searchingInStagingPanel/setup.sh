#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo "line 1" > myfile1
git add .
git commit -am "myfile1"

echo "line 2" >> myfile1
echo "line 3" >> myfile1
echo "line 4" >> myfile1
