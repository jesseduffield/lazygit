package commands

import (
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
