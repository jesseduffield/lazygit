package git_commands

import (
	"fmt"

	"github.com/jesseduffield/gocui"
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
	cmdArgs := NewGitCmd("remote").
		Arg("add", name, url).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *RemoteCommands) RemoveRemote(name string) error {
	cmdArgs := NewGitCmd("remote").
		Arg("remove", name).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *RemoteCommands) RenameRemote(oldRemoteName string, newRemoteName string) error {
	cmdArgs := NewGitCmd("remote").
		Arg("rename", oldRemoteName, newRemoteName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *RemoteCommands) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	cmdArgs := NewGitCmd("remote").
		Arg("set-url", remoteName, updatedUrl).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *RemoteCommands) DeleteRemoteBranch(task gocui.Task, remoteName string, branchName string) error {
	cmdArgs := NewGitCmd("push").
		Arg(remoteName, "--delete", branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).WithMutex(self.syncMutex).Run()
}

// CheckRemoteBranchExists Returns remote branch
func (self *RemoteCommands) CheckRemoteBranchExists(branchName string) bool {
	cmdArgs := NewGitCmd("show-ref").
		Arg("--verify", "--", fmt.Sprintf("refs/remotes/origin/%s", branchName)).
		ToArgv()

	_, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()

	return err == nil
}
