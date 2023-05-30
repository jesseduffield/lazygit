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

	cmdArgs := NewGitCmd("push").
		ArgIf(opts.Force, "--force-with-lease").
		ArgIf(opts.SetUpstream, "--set-upstream").
		ArgIf(opts.UpstreamRemote != "", opts.UpstreamRemote).
		ArgIf(opts.UpstreamBranch != "", opts.UpstreamBranch).
		ToArgv()

	cmdObj := self.cmd.New(cmdArgs).PromptOnCredentialRequest().WithMutex(self.syncMutex)
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
}

// Fetch fetch git repo
func (self *SyncCommands) FetchCmdObj(opts FetchOptions) oscommands.ICmdObj {
	cmdArgs := NewGitCmd("fetch").ToArgv()

	cmdObj := self.cmd.New(cmdArgs)
	if opts.Background {
		cmdObj.DontLog().FailOnCredentialRequest()
	} else {
		cmdObj.PromptOnCredentialRequest()
	}
	return cmdObj.WithMutex(self.syncMutex)
}

func (self *SyncCommands) Fetch(opts FetchOptions) error {
	cmdObj := self.FetchCmdObj(opts)
	return cmdObj.Run()
}

type PullOptions struct {
	RemoteName      string
	BranchName      string
	FastForwardOnly bool
}

func (self *SyncCommands) Pull(opts PullOptions) error {
	cmdArgs := NewGitCmd("pull").
		Arg("--no-edit").
		ArgIf(opts.FastForwardOnly, "--ff-only").
		ArgIf(opts.RemoteName != "", opts.RemoteName).
		ArgIf(opts.BranchName != "", opts.BranchName).
		ToArgv()

	// setting GIT_SEQUENCE_EDITOR to ':' as a way of skipping it, in case the user
	// has 'pull.rebase = interactive' configured.
	return self.cmd.New(cmdArgs).AddEnvVars("GIT_SEQUENCE_EDITOR=:").PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdArgs := NewGitCmd("fetch").
		Arg(remoteName).
		Arg(remoteBranchName + ":" + branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}

func (self *SyncCommands) FetchRemote(remoteName string) error {
	cmdArgs := NewGitCmd("fetch").
		Arg(remoteName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
