#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init
git remote add origin https://example.com/test.git
