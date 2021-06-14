package commands

import "fmt"

func (c *Git) CreateLightweightTag(tagName string, commitSha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag %s %s", tagName, commitSha))
}

func (c *Git) DeleteTag(tagName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag -d %s", tagName))
}

func (c *Git) PushTag(remoteName string, tagName string) error {
	return c.RunCommandWithCredentialsHandling(
		BuildGitCmdObjFromStr(fmt.Sprintf("push %s %s", remoteName, tagName)),
	)
}
