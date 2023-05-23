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
	return self.assert([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"}, expectedName)
}

func (self *Git) TagNamesAt(ref string, expectedNames []string) *Git {
	return self.assert([]string{"git", "tag", "--sort=v:refname", "--points-at", ref}, strings.Join(expectedNames, "\n"))
}

func (self *Git) assert(cmdArgs []string, expected string) *Git {
	self.assertWithRetries(func() (bool, string) {
		output, err := self.shell.runCommandWithOutput(cmdArgs)
		if err != nil {
			return false, fmt.Sprintf("Unexpected error running command: `%v`. Error: %s", cmdArgs, err.Error())
		}
		actual := strings.TrimSpace(output)
		return actual == expected, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expected, actual)
	})

	return self
}
