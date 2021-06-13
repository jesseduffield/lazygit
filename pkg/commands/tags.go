package commands

import "fmt"

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag %s %s", tagName, commitSha))
}

func (c *GitCommand) DeleteTag(tagName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag -d %s", tagName))
}

func (c *GitCommand) PushTag(remoteName string, tagName string) error {
	return c.RunCommandWithCredentialsHandling(
		BuildGitCmdObjFromStr(fmt.Sprintf("push %s %s", remoteName, tagName)),
	)
}
