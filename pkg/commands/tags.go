package commands

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	cmdStr := BuildGitCmdStr("tag", []string{tagName, commitSha}, nil)
	return c.RunCommand(cmdStr)
}

func (c *GitCommand) DeleteTag(tagName string) error {
	cmdStr := BuildGitCmdStr("tag", []string{"-d", tagName}, nil)
	return c.RunCommand(cmdStr)
}

func (c *GitCommand) PushTag(remoteName string, tagName string) error {
	cmdObj := BuildGitCmdObj("push", []string{remoteName, tagName}, nil)

	return c.RunCommandWithCredentialsHandling(cmdObj)
}
