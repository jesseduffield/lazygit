#!/bin/sh

echo "Running integration tests with $(git --version)"

# if the LAZYGIT_GOCOVERDIR env var is set, we'll capture code coverage data
if [ -n "$LAZYGIT_GOCOVERDIR" ]; then
  # Go expects us to either be running the test binary directly or running `go test`, but because
  # we're doing both and because we want to combine coverage data for both, we need to be a little
  # hacky. To capture the coverage data for the test runner we pass the test.gocoverdir positional
  # arg, but if we do that then the GOCOVERDIR env var (which you typically pass to the test binary) will be overwritten by the test runner. So we're passing LAZYGIT_COCOVERDIR instead
  # and then internally passing that to the test binary as GOCOVERDIR.
  go test -timeout 30m -cover -coverpkg=github.com/jesseduffield/lazygit/pkg/... pkg/integration/clients/*.go -args -test.gocoverdir="/tmp/code_coverage"
  EXITCODE=$?
else
  go test -timeout 30m pkg/integration/clients/*.go
  EXITCODE=$?
fi

# If per-test timings were collected (LAZYGIT_TEST_TIMING points at the file the
# harness appends to), print them sorted by slowest first so they show up in the
# CI log.
if [ -n "$LAZYGIT_TEST_TIMING" ] && [ -f "$LAZYGIT_TEST_TIMING" ]; then
  echo "Test timings (seconds):"
  sort -rn "$LAZYGIT_TEST_TIMING"
fi

exit $EXITCODE
