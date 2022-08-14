#!/bin/sh

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

cp ../../files/one.txt one.txt
cp ../../files/two.txt two.txt
cp ../../files/three.txt three.txt
git add .
git commit -am file1

cp ../../files/one_new.txt one.txt
cp ../../files/two_new.txt two.txt
cp ../../files/three_new.txt three.txt
