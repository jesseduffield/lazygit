#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

touch myfile
git add myfile
git commit -m "first commit"

git checkout -b other

for i in {1..20}
do
  echo "$i" > file
  git add .
  git commit -m "commit $i"
done

git checkout master

git bisect start other~2 other~13
