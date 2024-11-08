package git_commands

import (
	"github.com/go-errors/errors"
	"github.com/jesseduffield/gocui"
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
	ForceWithLease bool
	UpstreamRemote string
	UpstreamBranch string
	SetUpstream    bool
}

func (self *SyncCommands) PushCmdObj(task gocui.Task, opts PushOpts) (oscommands.ICmdObj, error) {
	if opts.UpstreamBranch != "" && opts.UpstreamRemote == "" {
		return nil, errors.New(self.Tr.MustSpecifyOriginError)
	}

	cmdArgs := NewGitCmd("push").
		ArgIf(opts.Force, "--force").
		ArgIf(opts.ForceWithLease, "--force-with-lease").
		ArgIf(opts.SetUpstream, "--set-upstream").
		ArgIf(opts.UpstreamRemote != "", opts.UpstreamRemote).
		ArgIf(opts.UpstreamBranch != "", "HEAD:"+opts.UpstreamBranch).
		ToArgv()

	cmdObj := self.cmd.New(cmdArgs).PromptOnCredentialRequest(task)
	return cmdObj, nil
}

func (self *SyncCommands) Push(task gocui.Task, opts PushOpts) error {
	cmdObj, err := self.PushCmdObj(task, opts)
	if err != nil {
		return err
	}

	return cmdObj.Run()
}

func (self *SyncCommands) fetchCommandBuilder(fetchAll bool) *GitCommandBuilder {
	return NewGitCmd("fetch").
		ArgIf(fetchAll, "--all").
		// avoid writing to .git/FETCH_HEAD; this allows running a pull
		// concurrently without getting errors
		ArgIf(self.version.IsAtLeast(2, 29, 0), "--no-write-fetch-head")
}

func (self *SyncCommands) FetchCmdObj(task gocui.Task) oscommands.ICmdObj {
	cmdArgs := self.fetchCommandBuilder(self.UserConfig().Git.FetchAll).ToArgv()

	cmdObj := self.cmd.New(cmdArgs)
	cmdObj.PromptOnCredentialRequest(task)
	return cmdObj
}

func (self *SyncCommands) Fetch(task gocui.Task) error {
	return self.FetchCmdObj(task).Run()
}

func (self *SyncCommands) FetchBackgroundCmdObj() oscommands.ICmdObj {
	cmdArgs := self.fetchCommandBuilder(self.UserConfig().Git.FetchAll).ToArgv()

	cmdObj := self.cmd.New(cmdArgs)
	cmdObj.DontLog().FailOnCredentialRequest()
	return cmdObj
}

func (self *SyncCommands) FetchBackground() error {
	return self.FetchBackgroundCmdObj().Run()
}

type PullOptions struct {
	RemoteName      string
	BranchName      string
	FastForwardOnly bool
	WorktreeGitDir  string
	WorktreePath    string
}

func (self *SyncCommands) Pull(task gocui.Task, opts PullOptions) error {
	cmdArgs := NewGitCmd("pull").
		Arg("--no-edit").
		ArgIf(opts.FastForwardOnly, "--ff-only").
		ArgIf(opts.RemoteName != "", opts.RemoteName).
		ArgIf(opts.BranchName != "", "refs/heads/"+opts.BranchName).
		GitDirIf(opts.WorktreeGitDir != "", opts.WorktreeGitDir).
		WorktreePathIf(opts.WorktreePath != "", opts.WorktreePath).
		ToArgv()

	// setting GIT_SEQUENCE_EDITOR to ':' as a way of skipping it, in case the user
	// has 'pull.rebase = interactive' configured.
	return self.cmd.New(cmdArgs).AddEnvVars("GIT_SEQUENCE_EDITOR=:").PromptOnCredentialRequest(task).Run()
}

func (self *SyncCommands) FastForward(
	task gocui.Task,
	branchName string,
	remoteName string,
	remoteBranchName string,
) error {
	cmdArgs := self.fetchCommandBuilder(false).
		Arg(remoteName).
		Arg("refs/heads/" + remoteBranchName + ":" + branchName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

func (self *SyncCommands) FetchRemote(task gocui.Task, remoteName string) error {
	cmdArgs := self.fetchCommandBuilder(false).
		Arg(remoteName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}
