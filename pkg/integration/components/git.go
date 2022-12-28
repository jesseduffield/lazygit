package components

import (
	"fmt"
	"strings"
)

type Git struct {
	*assertionHelper
	shell *Shell
}

func (self *Git) CurrentBranchName(expectedName string) *Git {
	return self.assert("git rev-parse --abbrev-ref HEAD", expectedName)
}

func (self *Git) assert(cmdStr string, expected string) *Git {
	self.assertWithRetries(func() (bool, string) {
		output, err := self.shell.runCommandWithOutput(cmdStr)
		if err != nil {
			return false, fmt.Sprintf("Unexpected error running command: `%s`. Error: %s", cmdStr, err.Error())
		}
		actual := strings.TrimSpace(output)
		return actual == expected, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expected, actual)
	})

	return self
}
