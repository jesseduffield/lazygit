package commands

import (
	"testing"

	"github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandWithOutput(t *testing.T) {
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
		s.test(NewOSCommand(logrus.New()).RunCommandWithOutput(s.command))
	}
}

func TestRunCommand(t *testing.T) {
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
		s.test(NewOSCommand(logrus.New()).RunCommand(s.command))
	}
}

func TestQuote(t *testing.T) {
	osCommand := NewOSCommand(logrus.New())

	actual := osCommand.Quote("hello `test`")

	expected := osCommand.Platform.escapedQuote + "hello \\`test\\`" + osCommand.Platform.escapedQuote

	assert.EqualValues(t, expected, actual)
}

func TestUnquote(t *testing.T) {
	osCommand := NewOSCommand(logrus.New())

	actual := osCommand.Unquote(`hello "test"`)

	expected := "hello test"

	assert.EqualValues(t, expected, actual)
}
