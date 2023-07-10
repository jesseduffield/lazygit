#!/bin/sh

# This is ugly, but older versions of git don't support the GIT_CONFIG_GLOBAL
# env var; the only way to run tests for these old versions is to copy our test
# config file to the actual global location. Move an existing file out of the
# way so that we can restore it at the end.
if test -f ~/.gitconfig; then
  mv ~/.gitconfig ~/.gitconfig.lazygit.bak
fi

cp test/global_git_config ~/.gitconfig

go test pkg/integration/clients/*.go
EXITCODE=$?

if test -f ~/.gitconfig.lazygit.bak; then
  mv ~/.gitconfig.lazygit.bak ~/.gitconfig
fi

exit $EXITCODE
