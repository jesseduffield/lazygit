#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

git commit --allow-empty -m "first commit"

mkdir -p one/two/three
echo test1 > one/two/three/file1
echo test2 > one/two/three/file2
echo test3 > one/two/three/file3
echo test4 > one/two/three/file4
echo test5 > one/two/file1
echo test6 > one/two/file2

git add .
git commit -m "blah"
