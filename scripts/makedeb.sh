#!/bin/bash
mkdir -p lazygit/{DEBIAN,usr/local/bin}
cat << EOF > lazygit/DEBIAN/control
Package: lazygit
Version: 0.44.1
Maintainer: jesseduffield
Architecture: amd64
Description: simple terminal UI for git commands
EOF
cp ../lazygit lazygit/usr/local/bin/
dpkg-deb --build lazygit
