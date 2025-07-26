package git_commands

import (
	"errors"
	"fmt"
	"path"
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

// CheckRemoteBranchExists returns a boolean indicating whether or not
// the given branch has an upstream.
func (self *RemoteCommands) CheckRemoteBranchExists(branchName string) bool {
	_, err := self.getRemoteRef(branchName)
	return err == nil
}

// Resolve what might be a aliased URL into a full URL
// SEE: `man -P 'less +/--get-url +n' git-ls-remote`
func (self *RemoteCommands) GetRemoteURL() (string, error) {
	remoteName := self.getRemoteName()
	if remoteName == "" {
		return "", errors.New("could not find upstream remote")
	}

	cmdArgs := NewGitCmd("ls-remote").
		Arg("--get-url", remoteName).
		ToArgv()

	url, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return strings.TrimSpace(url), err
}

func (self *RemoteCommands) getRemoteRef(branchName string) (string, error) {
	cmdArgs := NewGitCmd("rev-parse").
		Arg("--symbolic-full-name", fmt.Sprintf("%s@{upstream}", branchName)).
		ToArgv()

	remote, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil && branchName == "" {
		// if we couldn't find an upstream and the caller isn't asking about a specific
		// branch we'll return the first valid remote we find (if any)
		cmdArgs := NewGitCmd("rev-parse").
			Arg("--symbolic-full-name", "--remotes").
			ToArgv()
		remote, err = self.cmd.New(cmdArgs).DontLog().RunWithOutput()
		remote, _, _ = strings.Cut(remote, "\n")
	}

	return remote, err
}

func (self *RemoteCommands) getRemoteName() string {
	ref, err := self.getRemoteRef("")
	if err != nil {
		return ""
	}
	// refs/remotes/remote-name/branch
	//              ^^^^^^^^^^^
	return path.Base(path.Dir(ref))
}
