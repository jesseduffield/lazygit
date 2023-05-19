package git_commands

import (
	"fmt"
)

type RemoteCommands struct {
	*GitCommon
}

func NewRemoteCommands(gitCommon *GitCommon) *RemoteCommands {
	return &RemoteCommands{
		GitCommon: gitCommon,
	}
}

func (self *RemoteCommands) AddRemote(name string, url string) error {
	cmdStr := NewGitCmd("remote").
		Arg("add", self.cmd.Quote(name), self.cmd.Quote(url)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *RemoteCommands) RemoveRemote(name string) error {
	cmdStr := NewGitCmd("remote").
		Arg("remove", self.cmd.Quote(name)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *RemoteCommands) RenameRemote(oldRemoteName string, newRemoteName string) error {
	cmdStr := NewGitCmd("remote").
		Arg("rename", self.cmd.Quote(oldRemoteName), self.cmd.Quote(newRemoteName)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *RemoteCommands) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	cmdStr := NewGitCmd("remote").
		Arg("set-url", self.cmd.Quote(remoteName), self.cmd.Quote(updatedUrl)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *RemoteCommands) DeleteRemoteBranch(remoteName string, branchName string) error {
	cmdStr := NewGitCmd("push").
		Arg(self.cmd.Quote(remoteName), "--delete", self.cmd.Quote(branchName)).
		ToString()

	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

// CheckRemoteBranchExists Returns remote branch
func (self *RemoteCommands) CheckRemoteBranchExists(branchName string) bool {
	cmdStr := NewGitCmd("show-ref").
		Arg("--verify", "--", fmt.Sprintf("refs/remotes/origin/%s", self.cmd.Quote(branchName))).
		ToString()

	_, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()

	return err == nil
}
