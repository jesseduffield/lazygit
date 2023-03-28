//go:build !windows
// +build !windows

package oscommands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestOSCommandRunWithOutput(t *testing.T) {
	type scenario struct {
		command string
		test    func(string, error)
	}

	scenarios := []scenario{
		{
			"echo -n '123'",
			func(output string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "123", output)
			},
		},
		{
			"rmdir unexisting-folder",
			func(output string, err error) {
				assert.Regexp(t, "rmdir.*unexisting-folder.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		c := NewDummyOSCommand()
		s.test(c.Cmd.New(s.command).RunWithOutput())
	}
}

func TestOSCommandOpenFileDarwin(t *testing.T) {
	type scenario struct {
		filename string
		runner   *FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "test"`}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "test"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "filename with spaces"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		oSCmd := NewDummyOSCommandWithRunner(s.runner)
		oSCmd.Platform.OS = "darwin"
		oSCmd.UserConfig.OS.OpenCommand = "open {{filename}}"

		s.test(oSCmd.OpenFile(s.filename))
	}
}

// TestOSCommandOpenFileLinux tests the OpenFile command on Linux
func TestOSCommandOpenFileLinux(t *testing.T) {
	type scenario struct {
		filename string
		runner   *FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "test" > /dev/null`}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "test" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "filename with spaces" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "let's_test_with_single_quote",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "let's_test_with_single_quote" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "$USER.txt",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "\$USER.txt" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		oSCmd := NewDummyOSCommandWithRunner(s.runner)
		oSCmd.Platform.OS = "linux"
		oSCmd.UserConfig.OS.OpenCommand = `xdg-open {{filename}} > /dev/null`

		s.test(oSCmd.OpenFile(s.filename))
	}
}
