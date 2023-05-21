package oscommands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSCommandRun(t *testing.T) {
	type scenario struct {
		args []string
		test func(error)
	}

	scenarios := []scenario{
		{
			[]string{"rmdir", "unexisting-folder"},
			func(err error) {
				assert.Regexp(t, "rmdir.*unexisting-folder.*", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		c := NewDummyOSCommand()
		s.test(c.Cmd.New(s.args).Run())
	}
}

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
				if err := os.Mkdir("testDirectory", 0o644); err != nil {
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
		s.test(FileType(s.path))
		_ = os.RemoveAll(s.path)
	}
}

func TestOSCommandAppendLineToFile(t *testing.T) {
	type scenario struct {
		path  string
		setup func(string)
		test  func(string)
	}

	scenarios := []scenario{
		{
			filepath.Join(os.TempDir(), "testFile"),
			func(path string) {
				if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "hello\nworld\n", output)
			},
		},
		{
			filepath.Join(os.TempDir(), "emptyTestFile"),
			func(path string) {
				if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "world\n", output)
			},
		},
		{
			filepath.Join(os.TempDir(), "testFileWithNewline"),
			func(path string) {
				if err := os.WriteFile(path, []byte("hello\n"), 0o600); err != nil {
					panic(err)
				}
			},
			func(output string) {
				assert.EqualValues(t, "hello\nworld\n", output)
			},
		},
	}

	for _, s := range scenarios {
		s.setup(s.path)
		osCommand := NewDummyOSCommand()
		if err := osCommand.AppendLineToFile(s.path, "world"); err != nil {
			panic(err)
		}
		f, err := os.ReadFile(s.path)
		if err != nil {
			panic(err)
		}
		s.test(string(f))
		_ = os.RemoveAll(s.path)
	}
}
