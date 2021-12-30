package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (c *GitCommand) AddRemote(name string, url string) error {
	return c.Cmd.
		New(fmt.Sprintf("git remote add %s %s", c.Cmd.Quote(name), c.Cmd.Quote(url))).
		Run()
}

func (c *GitCommand) RemoveRemote(name string) error {
	return c.Cmd.
		New(fmt.Sprintf("git remote remove %s", c.Cmd.Quote(name))).
		Run()
}

func (c *GitCommand) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return c.Cmd.
		New(fmt.Sprintf("git remote rename %s %s", c.Cmd.Quote(oldRemoteName), c.Cmd.Quote(newRemoteName))).
		Run()
}

func (c *GitCommand) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return c.Cmd.
		New(fmt.Sprintf("git remote set-url %s %s", c.Cmd.Quote(remoteName), c.Cmd.Quote(updatedUrl))).
		Run()
}

func (c *GitCommand) DeleteRemoteBranch(remoteName string, branchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git push %s --delete %s", c.Cmd.Quote(remoteName), c.Cmd.Quote(branchName))
	cmdObj := c.Cmd.
		New(command)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}

func (c *GitCommand) DetectUnamePass(cmdObj oscommands.ICmdObj, promptUserForCredential func(string) string) error {
	return c.OSCommand.DetectUnamePass(cmdObj, c.GetCmdWriter(), promptUserForCredential)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branchName string) bool {
	_, err := c.Cmd.
		New(
			fmt.Sprintf("git show-ref --verify -- refs/remotes/origin/%s",
				c.Cmd.Quote(branchName),
			)).
		RunWithOutput()

	return err == nil
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	return c.GitConfig.Get("remote.origin.url")
}
