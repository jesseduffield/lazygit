//go:build !windows
// +build !windows

package integration

// this is the new way of running tests. See pkg/integration/integration_tests/commit.go
// for an example

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/creack/pty"
	"github.com/jesseduffield/lazygit/pkg/integration/helpers"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	mode := GetModeFromEnv()
	includeSkipped := os.Getenv("INCLUDE_SKIPPED") != ""

	parallelTotal := tryConvert(os.Getenv("PARALLEL_TOTAL"), 1)
	parallelIndex := tryConvert(os.Getenv("PARALLEL_INDEX"), 0)
	testNumber := 0

	err := RunTestsNew(
		t.Logf,
		runCmdHeadless,
		func(test *helpers.Test, f func() error) {
			defer func() { testNumber += 1 }()
			if testNumber%parallelTotal != parallelIndex {
				return
			}

			t.Run(test.Name(), func(t *testing.T) {
				err := f()
				assert.NoError(t, err)
			})
		},
		mode,
		includeSkipped,
	)

	assert.NoError(t, err)
}

func runCmdHeadless(cmd *exec.Cmd) error {
	cmd.Env = append(
		cmd.Env,
		"HEADLESS=true",
		"TERM=xterm",
	)

	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 100, Cols: 100})
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, f)

	return f.Close()
}

func tryConvert(numStr string, defaultVal int) int {
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return defaultVal
	}

	return num
}
