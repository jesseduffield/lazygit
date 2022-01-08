package git_commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type RemoteCommands struct {
	*common.Common

	cmd oscommands.ICmdObjBuilder
}

func NewRemoteCommands(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *RemoteCommands {
	return &RemoteCommands{
		Common: common,
		cmd:    cmd,
	}
}

func (self *RemoteCommands) AddRemote(name string, url string) error {
	return self.cmd.
		New(fmt.Sprintf("git remote add %s %s", self.cmd.Quote(name), self.cmd.Quote(url))).
		Run()
}

func (self *RemoteCommands) RemoveRemote(name string) error {
	return self.cmd.
		New(fmt.Sprintf("git remote remove %s", self.cmd.Quote(name))).
		Run()
}

func (self *RemoteCommands) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return self.cmd.
		New(fmt.Sprintf("git remote rename %s %s", self.cmd.Quote(oldRemoteName), self.cmd.Quote(newRemoteName))).
		Run()
}

func (self *RemoteCommands) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return self.cmd.
		New(fmt.Sprintf("git remote set-url %s %s", self.cmd.Quote(remoteName), self.cmd.Quote(updatedUrl))).
		Run()
}

func (self *RemoteCommands) DeleteRemoteBranch(remoteName string, branchName string) error {
	command := fmt.Sprintf("git push %s --delete %s", self.cmd.Quote(remoteName), self.cmd.Quote(branchName))
	return self.cmd.New(command).PromptOnCredentialRequest().Run()
}

// CheckRemoteBranchExists Returns remote branch
func (self *RemoteCommands) CheckRemoteBranchExists(branchName string) bool {
	_, err := self.cmd.
		New(
			fmt.Sprintf("git show-ref --verify -- refs/remotes/origin/%s",
				self.cmd.Quote(branchName),
			),
		).
		DontLog().
		RunWithOutput()

	return err == nil
}
