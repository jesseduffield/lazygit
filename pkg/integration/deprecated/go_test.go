//go:build !windows
// +build !windows

package deprecated

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
)

// Deprecated.

// This file is quite similar to integration/main.go. The main difference is that this file is
// run via `go test` whereas the other is run via `test/lazyintegration/main.go` which provides
//  a convenient gui wrapper around our integration tests. The `go test` approach is better
// for CI and for running locally in the background to ensure you haven't broken
// anything while making changes. If you want to visually see what's happening when a test is run,
// you'll need to take the other approach
//
// As for this file, to run an integration test, e.g. for test 'commit', go:
// go test pkg/gui/old_gui_test.go -run /commit
//
// To update a snapshot for an integration test, pass UPDATE_SNAPSHOTS=true
// UPDATE_SNAPSHOTS=true go test pkg/gui/old_gui_test.go -run /commit
//
// integration tests are run in test/integration/<test_name>/actual and the final test does
// not clean up that directory so you can cd into it to see for yourself what
// happened when a test fails.
//
// To override speed, pass e.g. `SPEED=1` as an env var. Otherwise we start each test
// at a high speed and then drop down to lower speeds upon each failure until finally
// trying at the original playback speed (speed 1). A speed of 2 represents twice the
// original playback speed. Speed may be a decimal.

func Test(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	mode := GetModeFromEnv()
	speedEnv := os.Getenv("SPEED")
	includeSkipped := os.Getenv("INCLUDE_SKIPPED") != ""

	parallelTotal := tryConvert(os.Getenv("PARALLEL_TOTAL"), 1)
	parallelIndex := tryConvert(os.Getenv("PARALLEL_INDEX"), 0)
	testNumber := 0

	err := RunTests(
		t.Logf,
		runCmdHeadless,
		func(test *IntegrationTest, f func(*testing.T) error) {
			defer func() { testNumber += 1 }()
			if testNumber%parallelTotal != parallelIndex {
				return
			}

			t.Run(test.Name, func(t *testing.T) {
				err := f(t)
				assert.NoError(t, err)
			})
		},
		mode,
		speedEnv,
		func(t *testing.T, expected string, actual string, prefix string) {
			t.Helper()
			assert.Equal(t, expected, actual, fmt.Sprintf("Unexpected %s. Expected:\n%s\nActual:\n%s\n", prefix, expected, actual))
		},
		includeSkipped,
	)

	assert.NoError(t, err)
}

func tryConvert(numStr string, defaultVal int) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return defaultVal
	}

	return num
}

func runCmdHeadless(cmd *exec.Cmd) error {
	cmd.Env = append(
		cmd.Env,
		"HEADLESS=true",
		"TERM=xterm",
	)

	// not writing stderr to the pty because we want to capture a panic if
	// there is one. But some commands will not be in tty mode if stderr is
	// not a terminal. We'll need to keep an eye out for that.
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	// these rows and columns are ignored because internally we use tcell's
	// simulation screen. However we still need the pty for the sake of
	// running other commands in a pty.
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 100, Cols: 100})
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, f)

	if cmd.Wait() != nil {
		// return an error with the stderr output
		return errors.New(stderr.String())
	}

	return f.Close()
}
