#!/bin/sh

# expecting an error so we're not setting this
# set -e

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

# common stuff
echo test > both-deleted.txt
git checkout -b conflict && git add both-deleted.txt
echo bothmodded > both-modded.txt && git add both-modded.txt
echo haha > deleted-them.txt && git add deleted-them.txt
echo haha2 > deleted-us.txt && git add deleted-us.txt
echo mod > modded.txt & git add modded.txt
echo mod > modded-staged.txt & git add modded-staged.txt
echo del > deleted.txt && git add deleted.txt
echo del > deleted-staged.txt && git add deleted-staged.txt
echo change-delete > change-delete.txt && git add change-delete.txt
echo delete-change > delete-change.txt && git add delete-change.txt
echo double-modded > double-modded.txt && git add double-modded.txt
echo "renamed\nhaha" > renamed.txt && git add renamed.txt
git commit -m one

# stuff on other branch
git branch conflict_second && git mv both-deleted.txt added-them-changed-us.txt
git commit -m "both-deleted.txt renamed in added-them-changed-us.txt"
echo blah > both-added.txt && git add both-added.txt
echo mod1 > both-modded.txt && git add both-modded.txt
rm deleted-them.txt && git add deleted-them.txt
echo modded > deleted-us.txt && git add deleted-us.txt
git commit -m "two"

# stuff on our branch
git checkout conflict_second
git mv both-deleted.txt changed-them-added-us.txt
git commit -m "both-deleted.txt renamed in changed-them-added-us.txt"
echo mod2 > both-modded.txt && git add both-modded.txt
echo blah2 > both-added.txt && git add both-added.txt
echo modded > deleted-them.txt && git add deleted-them.txt
rm deleted-us.txt && git add deleted-us.txt
git commit -m "three"
git reset --hard conflict_second
git merge conflict

echo "new" > new.txt
echo "new staged" > new-staged.txt && git add new-staged.txt
echo mod2 > modded.txt
echo mod2 > modded-staged.txt && git add modded-staged.txt
rm deleted.txt
rm deleted-staged.txt && git add deleted-staged.txt
echo change-delete2 > change-delete.txt && git add change-delete.txt
rm change-delete.txt
rm delete-change.txt && git add delete-change.txt
echo "changed" > delete-change.txt
echo "change1" > double-modded.txt && git add double-modded.txt
echo "change2" > double-modded.txt
echo before > added-changed.txt && git add added-changed.txt
echo after > added-changed.txt
rm renamed.txt && git add renamed.txt
echo "renamed\nhaha" > renamed2.txt && git add renamed2.txt
