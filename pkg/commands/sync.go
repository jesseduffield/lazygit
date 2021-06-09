package commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type PushOpts struct {
	Force             bool
	SetUpstream       bool
	DestinationRemote string
	DestinationBranch string
}

func (c *GitCommand) Push(opts PushOpts) (bool, error) {
	cmdObj := BuildGitCmdObj("push", []string{opts.DestinationRemote, opts.DestinationBranch},
		map[string]bool{
			"--follow-tags":      c.GetConfigValue("push.followTags") != "false",
			"--force-with-lease": opts.Force,
			"--set-upstream":     opts.SetUpstream,
		})

	err := c.RunCommandWithCredentialsPrompt(cmdObj)

	if isRejectionErr(err) {
		return true, nil
	}

	c.handleCredentialError(err)

	return false, nil
}

func isRejectionErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Updates were rejected")
}

type FetchOptions struct {
	RemoteName string
	BranchName string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	cmdObj := GetFetchCommandObj(opts)

	return c.RunCommandWithCredentialsHandling(cmdObj)
}

// FetchInBackground fails if credentials are requested
func (c *GitCommand) FetchInBackground(opts FetchOptions) error {
	cmdObj := GetFetchCommandObj(opts)

	cmdObj = c.FailOnCredentialsRequest(cmdObj)
	return c.oSCommand.RunExecutable(cmdObj)
}

func GetFetchCommandObj(opts FetchOptions) *oscommands.CmdObj {
	return BuildGitCmdObj("fetch", []string{opts.RemoteName, opts.BranchName}, nil)
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName, remoteBranchName + ":" + branchName}, nil)
	return c.RunCommandWithCredentialsHandling(cmdObj)
}

func (c *GitCommand) FetchRemote(remoteName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName}, nil)
	return c.RunCommandWithCredentialsHandling(cmdObj)
}
