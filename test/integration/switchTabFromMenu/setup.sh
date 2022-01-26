#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

touch file0
git add file0
git commit -am file0
git tag 0.0.1

touch file1
git add file1
git commit -am file0
git tag 0.0.2
