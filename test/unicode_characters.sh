#!/bin/bash

# this script will make a repo with a master and develop branch, where we end up
# on the master branch and if we try and merge master we get a merge conflict

# call this command from the test directory:
# ./generate_basic_repo.sh; cd testrepo; gg; cd ..

# -e means exit if something fails
# -x means print out simple commands before running them
set -ex

reponame="testrepo"

rm -rf ${reponame}
mkdir ${reponame}
cd ${reponame}

git init

# Add some ansi, unicode, zero width joiner caracters
cat <<EOT >> charstest.txt
ANSI      Œ	(U+0152	&OElig;	Latin capital ligature OE	Latin Extended-A)
          ¥	(0xA5	U+00A5	&yen;	yes sign)
          ƒ	(0x83	U+0192	&fnof;	Latin small letter f with hook)
ZWJ       https://en.wikipedia.org/wiki/Zero-width_joiner / https://unicode.org/Public/emoji/4.0/emoji-zwj-sequences.txt 👶(👨‍👦)
UNICODE   ☆ 🤓 え 术
EOT
git add charstest.txt
git commit -m "Test chars Œ¥ƒ👶👨‍👦☆ 🤓 え 术 commit"
echo "我喜歡編碼" >> charstest.txt
echo "நான் குறியீடு விரும்புகிறேன்" >> charstest.txt
git add charstest.txt
git commit -m "Test chars 我喜歡編碼 நான் குறியீடு விரும்புகிறேன் commit"

