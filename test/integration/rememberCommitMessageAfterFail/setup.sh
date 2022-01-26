#!/bin/sh

set -e

cd $1

git init
git config user.email "CI@example.com"
git config user.name "CI"

git checkout -b master

echo "file1" > file1
echo "file2" > file2
echo "disruptive" > disruptive
cat > .git/hooks/pre-commit <<EOL
#!/bin/bash
if [ -f disruptive ]; then
  exit 1
fi
exit 0
EOL
chmod +x .git/hooks/pre-commit
