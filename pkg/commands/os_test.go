package commands

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newDummyOSCommand() *OSCommand {
	return NewOSCommand(newDummyLog())
}

func TestOSCommandRunCommandWithOutput(t *testing.T) {
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
				assert.Regexp(t, ".*No such file or directory.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s.test(newDummyOSCommand().RunCommandWithOutput(s.command))
	}
}

func TestOSCommandRunCommand(t *testing.T) {
	type scenario struct {
		command string
		test    func(error)
	}

	scenarios := []scenario{
		{
			"rmdir unexisting-folder",
			func(err error) {
				assert.Regexp(t, ".*No such file or directory.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s.test(newDummyOSCommand().RunCommand(s.command))
	}
}

func TestOSCommandGetOpenCommand(t *testing.T) {
	type scenario struct {
		command func(string, ...string) *exec.Cmd
		test    func(string, string, error)
	}

	scenarios := []scenario{
		{
			func(name string, arg ...string) *exec.Cmd {
				return exec.Command("exit", "1")
			},
			func(name string, trail string, err error) {
				assert.EqualError(t, err, "Unsure what command to use to open this file")
			},
		},
		{
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				assert.Len(t, arg, 1)
				assert.Regexp(t, "xdg-open|cygstart|open", arg[0])
				return exec.Command("echo")
			},
			func(name string, trail string, err error) {
				assert.NoError(t, err)
				assert.Regexp(t, "xdg-open|cygstart|open", name)
				assert.Regexp(t, " \\&\\>/dev/null \\&|", trail)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := newDummyOSCommand()
		OSCmd.command = s.command

		s.test(OSCmd.getOpenCommand())
	}
}

func TestOSCommandOpenFile(t *testing.T) {
	type scenario struct {
		filename string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				return exec.Command("exit", "1")
			},
			func(err error) {
				assert.EqualError(t, err, "Unsure what command to use to open this file")
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				if name == "which" {
					return exec.Command("echo")
				}

				switch len(arg) {
				case 1:
					assert.Regexp(t, "open|cygstart", name)
					assert.EqualValues(t, "test", arg[0])
				case 3:
					assert.Equal(t, "xdg-open", name)
					assert.EqualValues(t, "test", arg[0])
					assert.Regexp(t, " \\&\\>/dev/null \\&|", arg[1])
					assert.EqualValues(t, "&", arg[2])
				default:
					assert.Fail(t, "Unexisting command given")
				}

				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := newDummyOSCommand()
		OSCmd.command = s.command

		s.test(OSCmd.OpenFile(s.filename))
	}
}

func TestOSCommandEditFile(t *testing.T) {
	type scenario struct {
		filename           string
		command            func(string, ...string) *exec.Cmd
		getenv             func(string) string
		getGlobalGitConfig func(string) (string, error)
		test               func(*exec.Cmd, error)
	}

	scenarios := []scenario{
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				return exec.Command("exit", "1")
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmd *exec.Cmd, err error) {
				assert.EqualError(t, err, "No editor defined in $VISUAL, $EDITOR, or git config")
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				if name == "which" {
					return exec.Command("exit", "1")
				}

				assert.EqualValues(t, "nano", name)

				return nil
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "nano", nil
			},
			func(cmd *exec.Cmd, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				if name == "which" {
					return exec.Command("exit", "1")
				}

				assert.EqualValues(t, "nano", name)

				return nil
			},
			func(env string) string {
				if env == "VISUAL" {
					return "nano"
				}

				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmd *exec.Cmd, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				if name == "which" {
					return exec.Command("exit", "1")
				}

				assert.EqualValues(t, "emacs", name)

				return nil
			},
			func(env string) string {
				if env == "EDITOR" {
					return "emacs"
				}

				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmd *exec.Cmd, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				if name == "which" {
					return exec.Command("echo")
				}

				assert.EqualValues(t, "vi", name)

				return nil
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmd *exec.Cmd, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := newDummyOSCommand()
		OSCmd.command = s.command
		OSCmd.getGlobalGitConfig = s.getGlobalGitConfig
		OSCmd.getenv = s.getenv

		s.test(OSCmd.EditFile(s.filename))
	}
}

func TestOSCommandQuote(t *testing.T) {
	osCommand := newDummyOSCommand()

	actual := osCommand.Quote("hello `test`")

	expected := osCommand.Platform.escapedQuote + "hello \\`test\\`" + osCommand.Platform.escapedQuote

	assert.EqualValues(t, expected, actual)
}

func TestOSCommandUnquote(t *testing.T) {
	osCommand := newDummyOSCommand()

	actual := osCommand.Unquote(`hello "test"`)

	expected := "hello test"

	assert.EqualValues(t, expected, actual)
}

func TestOSCommandFileType(t *testing.T) {
	type scenario struct {
		path  string
		setup func()
		test  func(string)
	}

	scenarios := []scenario{
		{
			"testFile",
			func() {
				if _, err := os.Create("testFile"); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "file", output)
			},
		},
		{
			"file with spaces",
			func() {
				if _, err := os.Create("file with spaces"); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "file", output)
			},
		},
		{
			"testDirectory",
			func() {
				if err := os.Mkdir("testDirectory", 0644); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "directory", output)
			},
		},
		{
			"nonExistant",
			func() {},
			func(output string) {
				assert.EqualValues(t, "other", output)
			},
		},
	}

	for _, s := range scenarios {
		s.setup()
		s.test(newDummyOSCommand().FileType(s.path))
		_ = os.RemoveAll(s.path)
	}
}
