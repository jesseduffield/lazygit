#!/bin/sh

set -e

cd $1

git init
git config user.email "CI@example.com"
git config user.name "CI"

git checkout -b master

echo "master1" > file
git add .
git commit -m "master1"

git checkout -b other

echo "other1" > file
git add .
git commit -m "other1"

git checkout master

echo "master2" > file
git add .
git commit -m "master2"

git checkout other
