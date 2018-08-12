#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

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
