//go:build !windows
// +build !windows

package oscommands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestOSCommandOpenFileDarwin is a function.
func TestOSCommandOpenFileDarwin(t *testing.T) {
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
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `open "test"`}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"filename with spaces",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `open "filename with spaces"`}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := NewDummyOSCommand()
		OSCmd.Platform.OS = "darwin"
		OSCmd.Command = s.command
		OSCmd.Config.GetUserConfig().OS.OpenCommand = "open {{filename}}"

		s.test(OSCmd.OpenFile(s.filename))
	}
}

// TestOSCommandOpenFileLinux tests the OpenFile command on Linux
func TestOSCommandOpenFileLinux(t *testing.T) {
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
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `xdg-open "test" > /dev/null`}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"filename with spaces",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `xdg-open "filename with spaces" > /dev/null`}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"let's_test_with_single_quote",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `xdg-open "let's_test_with_single_quote" > /dev/null`}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"$USER.txt",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "bash", name)
				assert.Equal(t, []string{"-c", `xdg-open "\$USER.txt" > /dev/null`}, arg)
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
		OSCmd.Platform.OS = "linux"
		OSCmd.Config.GetUserConfig().OS.OpenCommand = `xdg-open {{filename}} > /dev/null`

		s.test(OSCmd.OpenFile(s.filename))
	}
}
