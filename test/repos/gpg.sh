#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init

git config gpg.program $(which gpg)
git config user.signingkey E304229F # test key
git config commit.gpgsign true
git config credential.helper store
git config credential.helper cache 1

touch foo
git add foo

touch bar
git add bar