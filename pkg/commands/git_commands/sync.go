package git_commands

import (
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
	if opts.UpstreamBranch != "" && opts.UpstreamRemote == "" {
		return nil, errors.New(self.Tr.MustSpecifyOriginError)
	}

	cmdStr := NewGitCmd("push").
		ArgIf(opts.Force, "--force-with-lease").
		ArgIf(opts.SetUpstream, "--set-upstream").
		ArgIf(opts.UpstreamRemote != "", self.cmd.Quote(opts.UpstreamRemote)).
		ArgIf(opts.UpstreamBranch != "", self.cmd.Quote(opts.UpstreamBranch)).
		ToString()

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
	cmdStr := NewGitCmd("fetch").
		ArgIf(opts.RemoteName != "", self.cmd.Quote(opts.RemoteName)).
		ArgIf(opts.BranchName != "", self.cmd.Quote(opts.BranchName)).
		ToString()

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
	cmdStr := NewGitCmd("pull").
		Arg("--no-edit").
		ArgIf(opts.FastForwardOnly, "--ff-only").
		ArgIf(opts.RemoteName != "", self.cmd.Quote(opts.RemoteName)).
		ArgIf(opts.BranchName != "", self.cmd.Quote(opts.BranchName)).
		ToString()

	// setting GIT_SEQUENCE_EDITOR to ':' as a way of skipping it, in case the user
	// has 'pull.rebase = interactive' configured.
	return self.cmd.New(cmdStr).AddEnvVars("GIT_SEQUENCE_EDITOR=:").PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdStr := NewGitCmd("fetch").
		Arg(self.cmd.Quote(remoteName)).
		Arg(self.cmd.Quote(remoteBranchName) + ":" + self.cmd.Quote(branchName)).
		ToString()

	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FetchRemote(remoteName string) error {
	cmdStr := NewGitCmd("fetch").
		Arg(self.cmd.Quote(remoteName)).
		ToString()

	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
