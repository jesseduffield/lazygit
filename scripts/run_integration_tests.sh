#!/bin/sh

# This is ugly, but older versions of git don't support the GIT_CONFIG_GLOBAL
# env var; the only way to run tests for these old versions is to copy our test
# config file to the actual global location. Move an existing file out of the
# way so that we can restore it at the end.
if test -f ~/.gitconfig; then
  mv ~/.gitconfig ~/.gitconfig.lazygit.bak
fi

cp test/global_git_config ~/.gitconfig

# if the LAZYGIT_GOCOVERDIR env var is set, we'll capture code coverage data
if [ -n "$LAZYGIT_GOCOVERDIR" ]; then
  # Go expects us to either be running the test binary directly or running `go test`, but because
  # we're doing both and because we want to combine coverage data for both, we need to be a little
  # hacky. To capture the coverage data for the test runner we pass the test.gocoverdir positional
  # arg, but if we do that then the GOCOVERDIR env var (which you typically pass to the test binary) will be overwritten by the test runner. So we're passing LAZYGIT_COCOVERDIR instead
  # and then internally passing that to the test binary as GOCOVERDIR.
  go test -cover -coverpkg=github.com/jesseduffield/lazygit/pkg/... pkg/integration/clients/*.go -args -test.gocoverdir="/tmp/code_coverage"
  EXITCODE=$?

  # We're merging the coverage data for the sake of having fewer artefacts to upload.
  # We can't merge inline so we're merging to a tmp dir then moving back to the original.
  mkdir -p /tmp/code_coverage_merged
  go tool covdata merge -i=/tmp/code_coverage -o=/tmp/code_coverage_merged
  rm -rf /tmp/code_coverage
  mv /tmp/code_coverage_merged /tmp/code_coverage
else
  go test pkg/integration/clients/*.go
  EXITCODE=$?
fi

if test -f ~/.gitconfig.lazygit.bak; then
  mv ~/.gitconfig.lazygit.bak ~/.gitconfig
fi

exit $EXITCODE
