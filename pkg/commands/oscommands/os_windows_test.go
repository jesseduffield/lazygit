//go:build windows
// +build windows

package oscommands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

// handling this in a separate file because str.ToArgv has different behaviour if we're on windows

func TestOSCommandOpenFileWindows(t *testing.T) {
	type scenario struct {
		filename string
		runner   *FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/c", "start", "", "test"}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/c", "start", "", "test"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/c", "start", "", "filename with spaces"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "let's_test_with_single_quote",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/c", "start", "", "let's_test_with_single_quote"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "$USER.txt",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/c", "start", "", "$USER.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		oSCmd := NewDummyOSCommandWithRunner(s.runner)
		platform := &Platform{
			OS:       "windows",
			Shell:    "cmd",
			ShellArg: "/c",
		}
		oSCmd.Platform = platform
		oSCmd.Cmd.platform = platform
		oSCmd.UserConfig().OS.OpenCommand = `start "" {{filename}}`

		s.test(oSCmd.OpenFile(s.filename))
	}
}
