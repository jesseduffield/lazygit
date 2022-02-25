#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

# need a history of commits each containing various files

echo test0 > file0
git add .
git commit -am file0

echo test1 > file1
git add .
git commit -am file1

echo test2 > file2
git add .
git commit -am "file2"

echo test3 > file1
echo test4 > file2
git add .
git commit -am "file1 and file2"

echo test4 > file
git add .
git commit -am "file"
