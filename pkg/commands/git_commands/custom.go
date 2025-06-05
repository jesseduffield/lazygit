package git_commands

import (
	"fmt"
	"strings"

	"github.com/mgutz/str"
)

type CustomCommands struct {
	*GitCommon
}

func NewCustomCommands(gitCommon *GitCommon) *CustomCommands {
	return &CustomCommands{
		GitCommon: gitCommon,
	}
}

// Only to be used for the sake of running custom commands specified by the user.
// If you want to run a new command, try finding a place for it in one of the neighbouring
// files, or creating a new BlahCommands struct to hold it.
func (self *CustomCommands) RunWithOutput(cmdStr string) (string, error) {
	return self.cmd.New(str.ToArgv(cmdStr)).RunWithOutput()
}

// A function that can be used as a "runCommand" entry in the template.FuncMap of templates.
func (self *CustomCommands) TemplateFunctionRunCommand(cmdStr string) (string, error) {
	output, err := self.RunWithOutput(cmdStr)
	if err != nil {
		return "", err
	}
	output = strings.TrimRight(output, "\r\n")

	if strings.Contains(output, "\r\n") {
		return "", fmt.Errorf("command output contains newlines: %s", output)
	}

	return output, nil
}
