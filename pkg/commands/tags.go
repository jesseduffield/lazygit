package commands

import (
	"fmt"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	return c.RunCommand("git tag %s %s", tagName, commitSha)
}

func (c *GitCommand) DeleteTag(tagName string) error {
	return c.RunCommand("git tag -d %s", tagName)
}

func (c *GitCommand) PushTag(remoteName string, tagName string, promptUserForCredential func(CredentialKind) string) error {
	command := fmt.Sprintf("git push %s %s", remoteName, tagName)
	return c.GetOSCommand().DetectUnamePass(command, promptUserForCredential)
}
