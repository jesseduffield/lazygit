#!/bin/bash
mkdir -p lazygit_deb/{DEBIAN,usr/local/bin}
cat << EOF > lazygit_deb/DEBIAN/control
Package: lazygit
Version: 0.44.1
Maintainer: jesseduffield
Architecture: amd64
Description: simple terminal UI for git commands
EOF
cp ../lazygit lazygit_deb/usr/local/bin/
dpkg-deb --build lazygit_deb
mv lazygit_deb.deb lazygit.deb
