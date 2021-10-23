package commands

import (
	"fmt"
)

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	return c.RunCommand("git tag -- %s %s", c.OSCommand.Quote(tagName), commitSha)
}

func (c *GitCommand) DeleteTag(tagName string) error {
	return c.RunCommand("git tag -d %s", c.OSCommand.Quote(tagName))
}

func (c *GitCommand) PushTag(remoteName string, tagName string, promptUserForCredential func(string) string) error {
	cmdStr := fmt.Sprintf("git push %s %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(tagName))
	cmdObj := c.NewCmdObjFromStr(cmdStr)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}
