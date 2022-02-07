//go:build !windows
// +build !windows

package oscommands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestOSCommandOpenFileDarwin(t *testing.T) {
	type scenario struct {
		filename string
		linenumber int
		runner   *FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			filename: "test",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "test"`}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "test"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			linenumber: 1,
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

		s.test(oSCmd.OpenFile(s.filename, s.linenumber))
	}
}

// TestOSCommandOpenFileLinux tests the OpenFile command on Linux
func TestOSCommandOpenFileLinux(t *testing.T) {
	type scenario struct {
		filename string
		linenumber int
		runner   *FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			filename: "test",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "test" > /dev/null`}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "test" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "filename with spaces" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "let's_test_with_single_quote",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `xdg-open "let's_test_with_single_quote" > /dev/null`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "$USER.txt",
			linenumber: 1,
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

		s.test(oSCmd.OpenFile(s.filename, s.linenumber))
	}
}

func TestOSCommandOpenFileWithFilenameAndLine(t *testing.T) {
	type scenario struct {
		filename   string
		linenumber int
		runner     *FakeCmdObjRunner
		test       func(error)
	}

	scenarios := []scenario{
		{
			filename:   "test",
			linenumber: 1,
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"bash", "-c", `open "test":1`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		oSCmd := NewDummyOSCommandWithRunner(s.runner)
		oSCmd.Platform.OS = "darwin"
		oSCmd.UserConfig.OS.OpenCommand = `open {{filename}}:{{line}}`

		s.test(oSCmd.OpenFile(s.filename, s.linenumber))
	}
}
