#!/bin/sh

set -e

cd $1

git init
git config user.email "CI@example.com"
git config user.name "CI"


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
