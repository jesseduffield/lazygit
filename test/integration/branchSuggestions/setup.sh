#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test0 > file0
git add .
git commit -am file0

git checkout -b new-branch
git checkout -b new-branch-2
git checkout -b new-branch-3
git checkout -b old-branch
git checkout -b old-branch-2
git checkout -b old-branch-3
