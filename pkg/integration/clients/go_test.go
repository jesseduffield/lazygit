//go:build !windows
// +build !windows

package clients

// this is the new way of running tests. See pkg/integration/integration_tests/commit.go
// for an example

import (
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
		tests.Tests,
		t.Logf,
		runCmdHeadless,
		func(test *components.IntegrationTest, f func() error) {
			defer func() { testNumber += 1 }()
			if testNumber%parallelTotal != parallelIndex {
				return
			}

			t.Run(test.Name(), func(t *testing.T) {
				err := f()
				assert.NoError(t, err)
			})
		},
		components.CHECK_SNAPSHOT,
		0,
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
