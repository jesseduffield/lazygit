package git_commands

import (
	"fmt"
	"strings"

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

func (self *RemoteCommands) DeleteRemoteBranch(task gocui.Task, remoteName string, branchNames []string) error {
	cmdArgs := NewGitCmd("push").
		Arg(remoteName, "--delete").
		Arg(branchNames...).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

func (self *RemoteCommands) DeleteRemoteTag(task gocui.Task, remoteName string, tagName string) error {
	cmdArgs := NewGitCmd("push").
		Arg(remoteName, "--delete", tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

// CheckRemoteBranchExists Returns remote branch
func (self *RemoteCommands) CheckRemoteBranchExists(branchName string) bool {
	cmdArgs := NewGitCmd("show-ref").
		Arg("--verify", "--", fmt.Sprintf("refs/remotes/origin/%s", branchName)).
		ToArgv()

	_, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()

	return err == nil
}

// Resolve what might be a aliased URL into a full URL
// SEE: `man -P 'less +/--get-url +n' git-ls-remote`
func (self *RemoteCommands) GetRemoteURL(remoteName string) (string, error) {
	cmdArgs := NewGitCmd("ls-remote").
		Arg("--get-url", remoteName).
		ToArgv()

	url, err := self.cmd.New(cmdArgs).RunWithOutput()
	return strings.TrimSpace(url), err
}
