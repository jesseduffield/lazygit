package git_commands

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/samber/lo"
)

type SyncCommands struct {
	*GitCommon
}

func NewSyncCommands(gitCommon *GitCommon) *SyncCommands {
	return &SyncCommands{
		GitCommon: gitCommon,
	}
}

type PushOpts struct {
	Force          bool
	ForceWithLease bool
	SetUpstream    bool
	// The remote to push to. If empty, git picks it from its configuration,
	// and Refspecs must be empty too.
	Remote string
	// What to push, each in the form "refs/heads/<local branch>:<remote ref>".
	// If empty, git decides what to push based on push.default and
	// remote.<name>.push.
	Refspecs []string
}

func (self *SyncCommands) PushCmdObj(task gocui.Task, opts PushOpts) (*oscommands.CmdObj, error) {
	if len(opts.Refspecs) > 0 && opts.Remote == "" {
		return nil, errors.New(self.Tr.MustSpecifyOriginError)
	}

	cmdArgs := NewGitCmd("push").
		ArgIf(opts.Force, "--force").
		ArgIf(opts.ForceWithLease, "--force-with-lease").
		ArgIf(opts.SetUpstream, "--set-upstream").
		ArgIf(opts.Remote != "", opts.Remote).
		Arg(opts.Refspecs...).
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
		Arg("--no-write-fetch-head")
}

func (self *SyncCommands) FetchCmdObj(task gocui.Task) *oscommands.CmdObj {
	cmdArgs := self.fetchCommandBuilder(self.UserConfig().Git.FetchAll).ToArgv()

	cmdObj := self.cmd.New(cmdArgs)
	cmdObj.PromptOnCredentialRequest(task)
	return cmdObj
}

func (self *SyncCommands) Fetch(task gocui.Task) error {
	return self.FetchCmdObj(task).Run()
}

func (self *SyncCommands) FetchBackgroundCmdObj() *oscommands.CmdObj {
	cmdArgs := self.fetchCommandBuilder(self.UserConfig().Git.FetchAll).ToArgv()

	cmdObj := self.cmd.New(cmdArgs)
	cmdObj.DontLog().FailOnCredentialRequest()
	cmdObj.SuppressOutputUnlessError()
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

// Fetches the given branches of the given remote, updating their
// remote-tracking branches. Local branches are left alone, including the ones
// that track them.
func (self *SyncCommands) FetchRemoteBranches(
	task gocui.Task,
	remoteName string,
	remoteBranchNames []string,
) error {
	// The explicit destinations and the leading + make sure that the
	// remote-tracking branches are updated even when the remote branches were
	// rewritten, whatever the remote's fetch refspec says
	refspecs := lo.Map(remoteBranchNames, func(remoteBranchName string, _ int) string {
		return fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s",
			remoteBranchName, remoteName, remoteBranchName)
	})

	cmdArgs := self.fetchCommandBuilder(false).
		Arg(remoteName).
		Arg(refspecs...).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

func (self *SyncCommands) FetchRemote(task gocui.Task, remoteName string) error {
	cmdArgs := self.fetchCommandBuilder(false).
		Arg(remoteName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}
