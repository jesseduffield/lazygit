package oscommands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestOSCommandRunWithOutput is a function.
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
		s.test(NewDummyOS().RunWithOutput(s.command))
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
		s.test(NewDummyOS().RunCommand(s.command))
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
				assert.Equal(t, "open", name)
				assert.Equal(t, []string{"test"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"filename with spaces",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "open", name)
				assert.Equal(t, []string{"filename with spaces"}, arg)
				return secureexec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		OSCmd := NewDummyOS()
		OSCmd.command = s.command
		OSCmd.config.GetUserConfig().OS.OpenCommand = "open {{filename}}"

		s.test(OSCmd.OpenFile(s.filename))
	}
}

// TestOSCommandQuote is a function.
func TestOSCommandQuote(t *testing.T) {
	osCommand := NewDummyOS()

	osCommand.platform.OS = "linux"

	actual := osCommand.Quote("hello `test`")

	expected := osCommand.platform.EscapedQuote + "hello \\`test\\`" + osCommand.platform.EscapedQuote

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteSingleQuote tests the quote function with ' quotes explicitly for Linux
func TestOSCommandQuoteSingleQuote(t *testing.T) {
	osCommand := NewDummyOS()

	osCommand.platform.OS = "linux"

	actual := osCommand.Quote("hello 'test'")

	expected := osCommand.platform.EscapedQuote + "hello 'test'" + osCommand.platform.EscapedQuote

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteDoubleQuote tests the quote function with " quotes explicitly for Linux
func TestOSCommandQuoteDoubleQuote(t *testing.T) {
	osCommand := NewDummyOS()

	osCommand.platform.OS = "linux"

	actual := osCommand.Quote(`hello "test"`)

	expected := osCommand.platform.EscapedQuote + `hello \"test\"` + osCommand.platform.EscapedQuote

	assert.EqualValues(t, expected, actual)
}

// TestOSCommandQuoteWindows tests the quote function for Windows
func TestOSCommandQuoteWindows(t *testing.T) {
	osCommand := NewDummyOS()

	osCommand.platform.OS = "windows"

	actual := osCommand.Quote(`hello "test"`)

	expected := osCommand.platform.EscapedQuote + `hello "'"'"test"'"'"` + osCommand.platform.EscapedQuote

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
		s.test(NewDummyOS().FileType(s.path))
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
			s.test(NewDummyOS().CreateTempFile(s.filename, s.content))
		})
	}
}
