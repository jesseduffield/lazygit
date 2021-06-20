package commands

import (
	"fmt"
	"sync"
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

func (c *GitCommand) GetPullMode(mode string) string {
	if mode != "auto" {
		return mode
	}

	var isRebase bool
	var isFf bool
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		isRebase = c.GetConfigValue("pull.rebase") == "true"
		wg.Done()
	}()
	go func() {
		isFf = c.GetConfigValue("pull.ff") == "only"
		wg.Done()
	}()
	wg.Wait()

	if isRebase {
		return "rebase"
	} else if isFf {
		return "ff-only"
	} else {
		return "merge"
	}
}
