package commands

import (
	"fmt"
)

// Push pushes to a branch
func (c *GitCommand) Push(branchName string, force bool, upstream string, args string, promptUserForCredential func(string) string) error {
	followTagsFlag := "--follow-tags"
	if c.GetConfigValue("push.followTags") == "false" {
		followTagsFlag = ""
	}

	forceFlag := ""
	if force {
		forceFlag = "--force-with-lease"
	}

	setUpstreamArg := ""
	if upstream != "" {
		setUpstreamArg = "--set-upstream " + upstream
	}

	cmd := fmt.Sprintf("git push %s %s %s %s", followTagsFlag, forceFlag, setUpstreamArg, args)
	return c.OSCommand.DetectUnamePass(cmd, promptUserForCredential)
}

type FetchOptions struct {
	PromptUserForCredential func(string) string
	RemoteName              string
	BranchName              string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	command := "git fetch"

	if opts.RemoteName != "" {
		command = fmt.Sprintf("%s %s", command, opts.RemoteName)
	}
	if opts.BranchName != "" {
		command = fmt.Sprintf("%s %s", command, opts.BranchName)
	}

	return c.OSCommand.DetectUnamePass(command, func(question string) string {
		if opts.PromptUserForCredential != nil {
			return opts.PromptUserForCredential(question)
		}
		return "\n"
	})
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git fetch %s %s:%s", remoteName, remoteBranchName, branchName)
	return c.OSCommand.DetectUnamePass(command, promptUserForCredential)
}

func (c *GitCommand) FetchRemote(remoteName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git fetch %s", remoteName)
	return c.OSCommand.DetectUnamePass(command, promptUserForCredential)
}
