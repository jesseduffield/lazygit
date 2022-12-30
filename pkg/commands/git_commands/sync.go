package git_commands

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type SyncCommands struct {
	*GitCommon
}

func NewSyncCommands(gitCommon *GitCommon) *SyncCommands {
	return &SyncCommands{
		GitCommon: gitCommon,
	}
}

// Push pushes to a branch
type PushOpts struct {
	Force          bool
	UpstreamRemote string
	UpstreamBranch string
	SetUpstream    bool
}

func (self *SyncCommands) PushCmdObj(opts PushOpts) (oscommands.ICmdObj, error) {
	cmdStr := "git push"

	if opts.Force {
		if self.version.IsOlderThan(2, 30, 0) {
			cmdStr += " --force-with-lease"
		} else {
			cmdStr += " --force-with-lease --force-if-includes"
		}
	}

	if opts.SetUpstream {
		cmdStr += " --set-upstream"
	}

	if opts.UpstreamRemote != "" {
		cmdStr += " " + self.cmd.Quote(opts.UpstreamRemote)
	}

	if opts.UpstreamBranch != "" {
		if opts.UpstreamRemote == "" {
			return nil, errors.New(self.Tr.MustSpecifyOriginError)
		}
		cmdStr += " " + self.cmd.Quote(opts.UpstreamBranch)
	}

	cmdObj := self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex)
	return cmdObj, nil
}

func (self *SyncCommands) Push(opts PushOpts) error {
	cmdObj, err := self.PushCmdObj(opts)
	if err != nil {
		return err
	}

	return cmdObj.Run()
}

type FetchOptions struct {
	Background bool
	RemoteName string
	BranchName string
}

// Fetch fetch git repo
func (self *SyncCommands) Fetch(opts FetchOptions) error {
	cmdStr := "git fetch"

	if opts.RemoteName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, self.cmd.Quote(opts.RemoteName))
	}
	if opts.BranchName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, self.cmd.Quote(opts.BranchName))
	}

	cmdObj := self.cmd.New(cmdStr)
	if opts.Background {
		cmdObj.DontLog().FailOnCredentialRequest()
	} else {
		cmdObj.PromptOnCredentialRequest()
	}
	return cmdObj.WithMutex(self.syncMutex).Run()
}

type PullOptions struct {
	RemoteName      string
	BranchName      string
	FastForwardOnly bool
}

func (self *SyncCommands) Pull(opts PullOptions) error {
	cmdStr := "git pull --no-edit"

	if opts.FastForwardOnly {
		cmdStr += " --ff-only"
	}

	if opts.RemoteName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, self.cmd.Quote(opts.RemoteName))
	}
	if opts.BranchName != "" {
		cmdStr = fmt.Sprintf("%s %s", cmdStr, self.cmd.Quote(opts.BranchName))
	}

	// setting GIT_SEQUENCE_EDITOR to ':' as a way of skipping it, in case the user
	// has 'pull.rebase = interactive' configured.
	return self.cmd.New(cmdStr).AddEnvVars("GIT_SEQUENCE_EDITOR=:").PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdStr := fmt.Sprintf("git fetch %s %s:%s", self.cmd.Quote(remoteName), self.cmd.Quote(remoteBranchName), self.cmd.Quote(branchName))
	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FetchRemote(remoteName string) error {
	cmdStr := fmt.Sprintf("git fetch %s", self.cmd.Quote(remoteName))
	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
