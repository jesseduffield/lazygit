#!/bin/bash

# this script makes a repo with lots of commits

# call this command from the test directory:
# ./lots_of_commits.sh; cd testrepo; gg; cd ..

# -e means exit if something fails
# -x means print out simple commands before running them
set -ex

reponame="testrepo"

rm -rf ${reponame}
mkdir ${reponame}
cd ${reponame}

git init

i=2
end=100
while [ $i -le $end ]; do
    echo "file${i}" > file${i}
    git add file${i}
    git commit -m file${i}

    i=$(($i+1))
done

echo "unstaged change" > file100
