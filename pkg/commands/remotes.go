package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (c *GitCommand) AddRemote(name string, url string) error {
	return c.OSCommand.RunCommand("git remote add %s %s", name, url)
}

func (c *GitCommand) RemoveRemote(name string) error {
	return c.OSCommand.RunCommand("git remote remove %s", name)
}

func (c *GitCommand) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return c.OSCommand.RunCommand("git remote rename %s %s", oldRemoteName, newRemoteName)
}

func (c *GitCommand) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return c.OSCommand.RunCommand("git remote set-url %s %s", remoteName, updatedUrl)
}

func (c *GitCommand) DeleteRemoteBranch(remoteName string, branchName string) error {
	return c.OSCommand.RunCommand("git push %s --delete %s", remoteName, branchName)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branch *models.Branch) bool {
	_, err := c.OSCommand.RunCommandWithOutput(
		"git show-ref --verify -- refs/remotes/origin/%s",
		branch.Name,
	)

	return err == nil
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	return c.GetConfigValue("remote.origin.url")
}
