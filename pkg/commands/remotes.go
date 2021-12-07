package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (c *GitCommand) AddRemote(name string, url string) error {
	return c.Run(c.NewCmdObj(fmt.Sprintf("git remote add %s %s", c.OSCommand.Quote(name), c.OSCommand.Quote(url))))
}

func (c *GitCommand) RemoveRemote(name string) error {
	return c.Run(c.NewCmdObj(fmt.Sprintf("git remote remove %s", c.OSCommand.Quote(name))))
}

func (c *GitCommand) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return c.Run(c.NewCmdObj(fmt.Sprintf("git remote rename %s %s", c.OSCommand.Quote(oldRemoteName), c.OSCommand.Quote(newRemoteName))))
}

func (c *GitCommand) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return c.Run(c.NewCmdObj(fmt.Sprintf("git remote set-url %s %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(updatedUrl))))
}

func (c *GitCommand) DeleteRemoteBranch(remoteName string, branchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git push %s --delete %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(branchName))
	cmdObj := c.NewCmdObj(command)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}

func (c *GitCommand) DetectUnamePass(cmdObj oscommands.ICmdObj, promptUserForCredential func(string) string) error {
	return c.OSCommand.DetectUnamePass(cmdObj, c.GetCmdWriter(), promptUserForCredential)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branchName string) bool {
	_, err := c.RunWithOutput(c.NewCmdObj(
		fmt.Sprintf("git show-ref --verify -- refs/remotes/origin/%s",
			c.OSCommand.Quote(branchName),
		)))

	return err == nil
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	return c.GitConfig.Get("remote.origin.url")
}
