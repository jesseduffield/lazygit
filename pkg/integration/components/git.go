package components

import (
	"fmt"
	"log"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
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

func (self *Git) RemoteTagDeleted(ref string, tagName string) *Git {
	return self.expect([]string{"git", "ls-remote", ref, fmt.Sprintf("refs/tags/%s", tagName)}, func(s string) (bool, string) {
		return len(s) == 0, fmt.Sprintf("Expected tag %s to have been removed from %s", tagName, ref)
	})
}

func (self *Git) assert(cmdArgs []string, expected string) *Git {
	self.expect(cmdArgs, func(output string) (bool, string) {
		return output == expected, fmt.Sprintf("Expected current branch name to be '%s', but got '%s'", expected, output)
	})

	return self
}

func (self *Git) expect(cmdArgs []string, condition func(string) (bool, string)) *Git {
	self.assertWithRetries(func() (bool, string) {
		output, err := self.shell.runCommandWithOutput(cmdArgs)
		if err != nil {
			return false, fmt.Sprintf("Unexpected error running command: `%v`. Error: %s", cmdArgs, err.Error())
		}
		actual := strings.TrimSpace(output)
		return condition(actual)
	})

	return self
}

func (self *Git) Version() *git_commands.GitVersion {
	version, err := getGitVersion()
	if err != nil {
		log.Fatalf("Could not get git version: %v", err)
	}
	return version
}
