#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init
git config user.email "test@example.com"
git config user.name "Lazygit Tester"


touch foo
git add foo
git commit -m "init"
git branch -a
git branch test
git branch TEST
git checkout TEST
git checkout TeST
git checkout TesT
git checkout TEsT
git branch -a