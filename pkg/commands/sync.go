package commands

import (
	"fmt"
	"sync"

	"github.com/go-errors/errors"
)

// Push pushes to a branch
type PushOpts struct {
	Force                   bool
	UpstreamRemote          string
	UpstreamBranch          string
	SetUpstream             bool
	PromptUserForCredential func(string) string
}

func (c *GitCommand) Push(opts PushOpts) error {
	cmd := "git push"

	if c.GetConfigValue("push.followTags") != "false" {
		cmd += " --follow-tags"
	}

	if opts.Force {
		cmd += " --force-with-lease"
	}

	if opts.SetUpstream {
		cmd += " --set-upstream"
	}

	if opts.UpstreamRemote != "" {
		cmd += " " + c.OSCommand.Quote(opts.UpstreamRemote)
	}

	if opts.UpstreamBranch != "" {
		if opts.UpstreamRemote == "" {
			return errors.New(c.Tr.MustSpecifyOriginError)
		}
		cmd += " " + c.OSCommand.Quote(opts.UpstreamBranch)
	}

	return c.OSCommand.DetectUnamePass(cmd, opts.PromptUserForCredential)
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
		command = fmt.Sprintf("%s %s", command, c.OSCommand.Quote(opts.RemoteName))
	}
	if opts.BranchName != "" {
		command = fmt.Sprintf("%s %s", command, c.OSCommand.Quote(opts.BranchName))
	}

	return c.OSCommand.DetectUnamePass(command, func(question string) string {
		if opts.PromptUserForCredential != nil {
			return opts.PromptUserForCredential(question)
		}
		return "\n"
	})
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git fetch %s %s:%s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(remoteBranchName), c.OSCommand.Quote(branchName))
	return c.OSCommand.DetectUnamePass(command, promptUserForCredential)
}

func (c *GitCommand) FetchRemote(remoteName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git fetch %s", c.OSCommand.Quote(remoteName))
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
