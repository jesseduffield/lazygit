#!/bin/sh

set -e

cd $1
git config user.email "CI@example.com"
git config user.name "CI"

git init


echo test1 > myfile1

