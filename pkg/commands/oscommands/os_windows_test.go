package oscommands

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"unsafe"

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
				ExpectArgs([]string{"cmd", "/s", "/c", `start "" "test"`}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			filename: "test",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/s", "/c", `start "" "test"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "filename with spaces",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/s", "/c", `start "" "filename with spaces"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "let's_test_with_single_quote",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/s", "/c", `start "" "let's_test_with_single_quote"`}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			filename: "$USER.txt",
			runner: NewFakeRunner(t).
				ExpectArgs([]string{"cmd", "/s", "/c", `start "" "$USER.txt"`}, "", nil),
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
		oSCmd.UserConfig().OS.Open = `start "" {{filename}}`

		s.test(oSCmd.OpenFile(s.filename))
	}
}

var procGetConsoleTitle = kernel32.NewProc("GetConsoleTitleW")

// getConsoleTitle reads back the current console window title via the
// Win32 GetConsoleTitleW API, mirroring the SetConsoleTitleW call that
// UpdateWindowTitle makes.
func getConsoleTitle(t *testing.T) string {
	t.Helper()
	buf := make([]uint16, 1024)
	r1, _, err := procGetConsoleTitle.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if r1 == 0 {
		t.Skipf("GetConsoleTitleW unavailable in this environment (no attached console?): %v", err)
	}
	return syscall.UTF16ToString(buf[:r1])
}

// UpdateWindowTitle previously shelled out to `cmd /c title <name> - Lazygit`,
// which crashed lazygit whenever the current directory's basename contained
// a cmd.exe metacharacter such as & (see #5766: "test&aaa" was split into
// two commands, and cmd.exe tried to run "aaa" as a program). Calling
// SetConsoleTitleW directly bypasses cmd.exe's parsing, so the title is set
// verbatim regardless of what characters the directory name contains.
func TestUpdateWindowTitle_NameWithAmpersand(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	dir := t.TempDir()
	// t.TempDir() names are safe; nest a directory whose name reproduces
	// the exact character that crashed cmd.exe in #5766.
	problematic := filepath.Join(dir, "test&aaa")
	if err := os.Mkdir(problematic, 0o755); err != nil {
		t.Fatalf("os.Mkdir: %v", err)
	}
	if err := os.Chdir(problematic); err != nil {
		t.Fatalf("os.Chdir: %v", err)
	}

	osCommand := NewDummyOSCommand()
	err = osCommand.UpdateWindowTitle()
	assert.NoError(t, err, "UpdateWindowTitle must not error out for directory names containing '&'")

	title := getConsoleTitle(t)
	assert.True(t, strings.HasPrefix(title, "test&aaa"),
		"expected console title to start with the literal directory name %q, got %q", "test&aaa", title)
	assert.Contains(t, title, "Lazygit")
}

func TestUpdateWindowTitle_PlainName(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	dir := t.TempDir()
	plain := filepath.Join(dir, "plain-repo")
	if err := os.Mkdir(plain, 0o755); err != nil {
		t.Fatalf("os.Mkdir: %v", err)
	}
	if err := os.Chdir(plain); err != nil {
		t.Fatalf("os.Chdir: %v", err)
	}

	osCommand := NewDummyOSCommand()
	if err := osCommand.UpdateWindowTitle(); err != nil {
		t.Fatalf("UpdateWindowTitle: %v", err)
	}

	title := getConsoleTitle(t)
	assert.Equal(t, "plain-repo - Lazygit", title)
}
