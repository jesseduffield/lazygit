package commands

type PushOpts struct {
	Force             bool
	SetUpstream       bool
	DestinationRemote string
	DestinationBranch string
}

func (c *GitCommand) Push(opts PushOpts) error {
	cmdObj := BuildGitCmdObj("push", []string{opts.DestinationRemote, opts.DestinationBranch},
		map[string]bool{
			"--follow-tags":      c.GetConfigValue("push.followTags") != "false",
			"--force-with-lease": opts.Force,
			"--set-upstream":     opts.SetUpstream,
		})

	return c.DetectUnamePass(cmdObj)
}

type FetchOptions struct {
	// if Background is true, we will not prompt the user for a credential
	Background bool
	RemoteName string
	BranchName string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	cmdObj := BuildGitCmdObj("fetch", []string{opts.RemoteName, opts.BranchName}, nil)

	if opts.Background {
		cmdObj = c.FailOnCredentialsRequest(cmdObj)
		return c.oSCommand.RunExecutable(cmdObj)
	} else {
		return c.DetectUnamePass(cmdObj)
	}
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName, remoteBranchName + ":" + branchName}, nil)
	return c.DetectUnamePass(cmdObj)
}

func (c *GitCommand) FetchRemote(remoteName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName}, nil)
	return c.DetectUnamePass(cmdObj)
}
