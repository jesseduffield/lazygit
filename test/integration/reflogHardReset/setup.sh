#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test0 > file0
git add .
git commit -am file0

echo test1 > file1
git add .
git commit -am file1

echo test2 > file2
git add .
git commit -am file2

echo "line one" > file4
git add .
git commit -am file4

git checkout -b branch2

echo "line two" >> file4
git add .
git commit -am file4

echo "line three" >> file4
git add .
git commit -am file4

echo "line two" >> file2
git add .
git commit -am file2
