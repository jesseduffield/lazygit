package oscommands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestOSCommandRunCommandWithOutput is a function.
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
				assert.Regexp(t, "rmdir.*unexisting-folder.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s.test(NewDummyOSCommand().RunCommandWithOutput(s.command))
	}
}

// TestOSCommandRunCommand is a function.
func TestOSCommandRunCommand(t *testing.T) {
	type scenario struct {
		command string
		test    func(error)
	}

	scenarios := []scenario{
		{
			"rmdir unexisting-folder",
			func(err error) {
				assert.Regexp(t, "rmdir.*unexisting-folder.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s.test(NewDummyOSCommand().RunCommand(s.command))
	}
}

// TestOSCommandOpenFile is a function.
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
	if runtime.GOOS == "windows" {
		return
	}

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
				assert.Equal(t, []string{"-c", "xdg-open \"test\" > /dev/null"}, arg)
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
				assert.Equal(t, []string{"-c", "xdg-open \"filename with spaces\" > /dev/null"}, arg)
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
				assert.Equal(t, []string{"-c", "xdg-open \"let's_test_with_single_quote\" > /dev/null"}, arg)
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
				assert.Equal(t, []string{"-c", "xdg-open \"\\$USER.txt\" > /dev/null"}, arg)
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

// TestOSCommandOpenFileWindows tests the OpenFile command on Linux
func TestOSCommandOpenFileWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		return
	}

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
		OSCmd.Config.GetUserConfig().OS.OpenCommand = `cmd /c start "" {{filename}}`

		s.test(OSCmd.OpenFile(s.filename))
	}
}

// TestOSCommandQuote is a function.
func TestOSCommandQuote(t *testing.T) {
	osCommand := NewDummyOSCommand()

	osCommand.Platform.OS = "linux"

	actual := osCommand.Quote("hello `test`")

	expected := "\"hello \\`test\\`\""

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteSingleQuote tests the quote function with ' quotes explicitly for Linux
func TestOSCommandQuoteSingleQuote(t *testing.T) {
	osCommand := NewDummyOSCommand()

	osCommand.Platform.OS = "linux"

	actual := osCommand.Quote("hello 'test'")

	expected := `"hello 'test'"`

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteDoubleQuote tests the quote function with " quotes explicitly for Linux
func TestOSCommandQuoteDoubleQuote(t *testing.T) {
	osCommand := NewDummyOSCommand()

	osCommand.Platform.OS = "linux"

	actual := osCommand.Quote(`hello "test"`)

	expected := `"hello \"test\""`

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteWindows tests the quote function for Windows
func TestOSCommandQuoteWindows(t *testing.T) {
	osCommand := NewDummyOSCommand()

	osCommand.Platform.OS = "windows"

	actual := osCommand.Quote(`hello "test" 'test2'`)

	expected := `\"hello "'"'"test"'"'" 'test2'\"`

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandFileType is a function.
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
		s.test(NewDummyOSCommand().FileType(s.path))
		_ = os.RemoveAll(s.path)
	}
}

func TestOSCommandCreateTempFile(t *testing.T) {
	type scenario struct {
		testName string
		filename string
		content  string
		test     func(string, error)
	}

	scenarios := []scenario{
		{
			"valid case",
			"filename",
			"content",
			func(path string, err error) {
				assert.NoError(t, err)

				content, err := ioutil.ReadFile(path)
				assert.NoError(t, err)

				assert.Equal(t, "content", string(content))
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(NewDummyOSCommand().CreateTempFile(s.filename, s.content))
		})
	}
}
