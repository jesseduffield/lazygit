//go:build !windows
// +build !windows

package clients

// this is the new way of running tests. See pkg/integration/integration_tests/commit.go
// for an example

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/creack/pty"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	parallelTotal := tryConvert(os.Getenv("PARALLEL_TOTAL"), 1)
	parallelIndex := tryConvert(os.Getenv("PARALLEL_INDEX"), 0)
	testNumber := 0

	err := components.RunTests(
		tests.GetTests(),
		t.Logf,
		runCmdHeadless,
		func(test *components.IntegrationTest, f func() error) {
			defer func() { testNumber += 1 }()
			if testNumber%parallelTotal != parallelIndex {
				return
			}

			t.Run(test.Name(), func(t *testing.T) {
				t.Parallel()
				err := f()
				assert.NoError(t, err)
			})
		},
		false,
		0,
		// allowing two attempts at the test. If a test fails intermittently,
		// there may be a concurrency issue that we need to resolve.
		2,
	)

	assert.NoError(t, err)
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
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 300, Cols: 300})
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
