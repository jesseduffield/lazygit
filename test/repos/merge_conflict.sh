#!/bin/sh
set -ex; rm -rf repo; mkdir repo; cd repo

git init
git config user.email "test@example.com"
git config user.name "Lazygit Tester"


function add_spacing {
  for i in {1..60}
  do
    echo "..." >> $1
  done
}

mkdir directory
echo "test1" > directory/file
echo "test1" > directory/file2


echo "Here is a story that has been told throuhg the ages" >> file1

git add file1
git add directory
git commit -m "first commit"

git checkout -b feature/cherry-picking

echo "this is file number 1 that I'm going to cherry-pick" > cherrypicking1
echo "this is file number 2 that I'm going to cherry-pick" > cherrypicking2

git add .

git commit -am "first commit freshman year"

echo "this is file number 3 that I'm going to cherry-pick" > cherrypicking3

git add .

git commit -am "second commit subway eat fresh"

echo "this is file number 4 that I'm going to cherry-pick" > cherrypicking4

git add .

git commit -am "third commit fresh"

echo "this is file number 5 that I'm going to cherry-pick" > cherrypicking5

git add .

git commit -am "fourth commit cool"

echo "this is file number 6 that I'm going to cherry-pick" > cherrypicking6

git add .

git commit -am "fifth commit nice"

echo "this is file number 7 that I'm going to cherry-pick" > cherrypicking7

git add .

git commit -am "sixth commit haha"

echo "this is file number 8 that I'm going to cherry-pick" > cherrypicking8

git add .

git commit -am "seventh commit yeah"

echo "this is file number 9 that I'm going to cherry-pick" > cherrypicking9

git add .

git commit -am "eighth commit woo"


git checkout -b develop
echo "once upon a time there was a dog" >> file1
add_spacing file1
echo "once upon a time there was another dog" >> file1
git add file1
echo "test2" > directory/file
echo "test2" > directory/file2
git add directory
git commit -m "first commit on develop"


git checkout master
echo "once upon a time there was a cat" >> file1
add_spacing file1
echo "once upon a time there was another cat" >> file1
git add file1
echo "test3" > directory/file
echo "test3" > directory/file2
git add directory
git commit -m "first commit on master"


git checkout develop
echo "once upon a time there was a mouse" >> file3
git add file3
git commit -m "second commit on develop"


git checkout master
echo "once upon a time there was a horse" >> file3
git add file3
git commit -m "second commit on master"


git checkout develop
echo "once upon a time there was a mouse" >> file4
git add file4
git commit -m "third commit on develop"


git checkout master
echo "once upon a time there was a horse" >> file4
git add file4
git commit -m "third commit on master"


git checkout develop
echo "once upon a time there was a mouse" >> file5
git add file5
git commit -m "fourth commit on develop"


git checkout master
echo "once upon a time there was a horse" >> file5
git add file5
git commit -m "fourth commit on master"


# this is for the autostash feature

git checkout -b base_branch

echo "original1\noriginal2\noriginal3" > file
git add file
git commit -m "file"

git checkout -b other_branch

git checkout base_branch

echo "new1\noriginal2\noriginal3" > file
git add file
git commit -m "file changed"

git checkout other_branch

echo "new2\noriginal2\noriginal3" > file
