package commands

import (
	"os"
	"strings"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// GitCommand is our main git interface
type GitCommand struct {
	Blame       *git_commands.BlameCommands
	Branch      *git_commands.BranchCommands
	Commit      *git_commands.CommitCommands
	Config      *git_commands.ConfigCommands
	Custom      *git_commands.CustomCommands
	Diff        *git_commands.DiffCommands
	File        *git_commands.FileCommands
	Flow        *git_commands.FlowCommands
	Patch       *git_commands.PatchCommands
	Rebase      *git_commands.RebaseCommands
	Remote      *git_commands.RemoteCommands
	Stash       *git_commands.StashCommands
	Status      *git_commands.StatusCommands
	Submodule   *git_commands.SubmoduleCommands
	Sync        *git_commands.SyncCommands
	Tag         *git_commands.TagCommands
	WorkingTree *git_commands.WorkingTreeCommands
	Bisect      *git_commands.BisectCommands
	Worktree    *git_commands.WorktreeCommands
	Version     *git_commands.GitVersion
	RepoPaths   *git_commands.RepoPaths

	Loaders Loaders
}

type Loaders struct {
	BranchLoader       *git_commands.BranchLoader
	CommitFileLoader   *git_commands.CommitFileLoader
	CommitLoader       *git_commands.CommitLoader
	FileLoader         *git_commands.FileLoader
	ReflogCommitLoader *git_commands.ReflogCommitLoader
	RemoteLoader       *git_commands.RemoteLoader
	StashLoader        *git_commands.StashLoader
	TagLoader          *git_commands.TagLoader
	Worktrees          *git_commands.WorktreeLoader
}

func NewGitCommand(
	cmn *common.Common,
	version *git_commands.GitVersion,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
) (*GitCommand, error) {
	repoPaths, err := git_commands.GetRepoPaths(osCommand.Cmd, version)
	if err != nil {
		return nil, errors.Errorf("Error getting repo paths: %v", err)
	}

	err = os.Chdir(repoPaths.WorktreePath())
	if err != nil {
		return nil, utils.WrapError(err)
	}

	repository, err := gogit.PlainOpenWithOptions(
		repoPaths.WorktreeGitDirPath(),
		&gogit.PlainOpenOptions{DetectDotGit: false, EnableDotGitCommonDir: true},
	)
	if err != nil {
		if strings.Contains(err.Error(), `unquoted '\' must be followed by new line`) {
			return nil, errors.New(cmn.Tr.GitconfigParseErr)
		}
		return nil, err
	}

	return NewGitCommandAux(
		cmn,
		version,
		osCommand,
		gitConfig,
		repoPaths,
		repository,
	), nil
}

func NewGitCommandAux(
	cmn *common.Common,
	version *git_commands.GitVersion,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
	repoPaths *git_commands.RepoPaths,
	repo *gogit.Repository,
) *GitCommand {
	cmd := NewGitCmdObjBuilder(cmn.Log, osCommand.Cmd)

	// here we're doing a bunch of dependency injection for each of our commands structs.
	// This is admittedly messy, but allows us to test each command struct in isolation,
	// and allows for better namespacing when compared to having every method living
	// on the one struct.
	// common ones are: cmn, osCommand, dotGitDir, configCommands
	configCommands := git_commands.NewConfigCommands(cmn, gitConfig, repo)

	gitCommon := git_commands.NewGitCommon(cmn, version, cmd, osCommand, repoPaths, repo, configCommands)

	fileLoader := git_commands.NewFileLoader(gitCommon, cmd, configCommands)
	statusCommands := git_commands.NewStatusCommands(gitCommon)
	flowCommands := git_commands.NewFlowCommands(gitCommon)
	remoteCommands := git_commands.NewRemoteCommands(gitCommon)
	branchCommands := git_commands.NewBranchCommands(gitCommon)
	syncCommands := git_commands.NewSyncCommands(gitCommon)
	tagCommands := git_commands.NewTagCommands(gitCommon)
	commitCommands := git_commands.NewCommitCommands(gitCommon)
	customCommands := git_commands.NewCustomCommands(gitCommon)
	diffCommands := git_commands.NewDiffCommands(gitCommon)
	fileCommands := git_commands.NewFileCommands(gitCommon)
	submoduleCommands := git_commands.NewSubmoduleCommands(gitCommon)
	workingTreeCommands := git_commands.NewWorkingTreeCommands(gitCommon, submoduleCommands, fileLoader)
	rebaseCommands := git_commands.NewRebaseCommands(gitCommon, commitCommands, workingTreeCommands)
	stashCommands := git_commands.NewStashCommands(gitCommon, fileLoader, workingTreeCommands)
	patchBuilder := patch.NewPatchBuilder(cmn.Log,
		func(from string, to string, reverse bool, filename string, plain bool) (string, error) {
			return workingTreeCommands.ShowFileDiff(from, to, reverse, filename, plain)
		})
	patchCommands := git_commands.NewPatchCommands(gitCommon, rebaseCommands, commitCommands, statusCommands, stashCommands, patchBuilder)
	bisectCommands := git_commands.NewBisectCommands(gitCommon)
	worktreeCommands := git_commands.NewWorktreeCommands(gitCommon)
	blameCommands := git_commands.NewBlameCommands(gitCommon)

	branchLoader := git_commands.NewBranchLoader(cmn, cmd, branchCommands.CurrentBranchInfo, configCommands)
	commitFileLoader := git_commands.NewCommitFileLoader(cmn, cmd)
	commitLoader := git_commands.NewCommitLoader(cmn, cmd, statusCommands.RebaseMode, gitCommon)
	reflogCommitLoader := git_commands.NewReflogCommitLoader(cmn, cmd)
	remoteLoader := git_commands.NewRemoteLoader(cmn, cmd, repo.Remotes)
	worktreeLoader := git_commands.NewWorktreeLoader(gitCommon)
	stashLoader := git_commands.NewStashLoader(cmn, cmd)
	tagLoader := git_commands.NewTagLoader(cmn, cmd)

	return &GitCommand{
		Blame:       blameCommands,
		Branch:      branchCommands,
		Commit:      commitCommands,
		Config:      configCommands,
		Custom:      customCommands,
		Diff:        diffCommands,
		File:        fileCommands,
		Flow:        flowCommands,
		Patch:       patchCommands,
		Rebase:      rebaseCommands,
		Remote:      remoteCommands,
		Stash:       stashCommands,
		Status:      statusCommands,
		Submodule:   submoduleCommands,
		Sync:        syncCommands,
		Tag:         tagCommands,
		Bisect:      bisectCommands,
		WorkingTree: workingTreeCommands,
		Worktree:    worktreeCommands,
		Version:     version,
		Loaders: Loaders{
			BranchLoader:       branchLoader,
			CommitFileLoader:   commitFileLoader,
			CommitLoader:       commitLoader,
			FileLoader:         fileLoader,
			ReflogCommitLoader: reflogCommitLoader,
			RemoteLoader:       remoteLoader,
			Worktrees:          worktreeLoader,
			StashLoader:        stashLoader,
			TagLoader:          tagLoader,
		},
		RepoPaths: repoPaths,
	}
}

func VerifyInGitRepo(osCommand *oscommands.OSCommand) error {
	return osCommand.Cmd.New(git_commands.NewGitCmd("rev-parse").Arg("--git-dir").ToArgv()).DontLog().Run()
}
