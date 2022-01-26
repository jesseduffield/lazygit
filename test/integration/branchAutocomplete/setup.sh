#!/bin/sh

set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

touch myfile.txt
git add .
git commit -m "initial commit"

git checkout -b one
git checkout -b two
git checkout -b three
git checkout -b four

