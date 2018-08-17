#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init

git checkout -b asdf
touch asdf
