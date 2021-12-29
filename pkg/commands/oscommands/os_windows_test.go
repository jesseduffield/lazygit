//go:build windows
// +build windows

package oscommands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestOSCommandOpenFileWindows tests the OpenFile command on Linux
func TestOSCommandOpenFileWindows(t *testing.T) {
	type scenario struct {
		filename string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				return secureexec.Command("exit", "1")
			},
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "cmd", name)
				assert.Equal(t, []string{"/c", "start", "", "test"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"filename with spaces",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "cmd", name)
				assert.Equal(t, []string{"/c", "start", "", "filename with spaces"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"let's_test_with_single_quote",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "cmd", name)
				assert.Equal(t, []string{"/c", "start", "", "let's_test_with_single_quote"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"$USER.txt",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "cmd", name)
				assert.Equal(t, []string{"/c", "start", "", "$USER.txt"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := NewDummyOSCommand()
		OSCmd.Command = s.command
		OSCmd.Platform.OS = "windows"
		OSCmd.UserConfig.OS.OpenCommand = `start "" {{filename}}`

		s.test(OSCmd.OpenFile(s.filename))
	}
}
