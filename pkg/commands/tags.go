package commands

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	return c.OSCommand.RunCommand("git tag %s %s", tagName, commitSha)
}

func (c *GitCommand) DeleteTag(tagName string) error {
	return c.OSCommand.RunCommand("git tag -d %s", tagName)
}

func (c *GitCommand) PushTag(remoteName string, tagName string) error {
	return c.OSCommand.RunCommand("git push %s %s", remoteName, tagName)
}
