#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

for i in {1..20}
do
  echo "$i" > file
  git add .
  git commit -m "commit $i"
done
