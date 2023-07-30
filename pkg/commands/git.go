package commands

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/sasha-s/go-deadlock"
	"github.com/spf13/afero"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// GitCommand is our main git interface
type GitCommand struct {
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
	syncMutex *deadlock.Mutex,
) (*GitCommand, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, utils.WrapError(err)
	}

	// converting to forward slashes for the sake of windows (which uses backwards slashes). We want everything
	// to have forward slashes internally
	currentPath = filepath.ToSlash(currentPath)

	gitDir := env.GetGitDirEnv()
	if gitDir != "" {
		// we've been given the git directory explicitly so no need to navigate to it
		_, err := cmn.Fs.Stat(gitDir)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	} else {
		// we haven't been given the git dir explicitly so we assume it's in the current working directory as `.git/` (or an ancestor directory)

		rootDirectory, err := findWorktreeRoot(cmn.Fs, currentPath)
		if err != nil {
			return nil, utils.WrapError(err)
		}
		currentPath = rootDirectory
		err = os.Chdir(rootDirectory)
		if err != nil {
			return nil, utils.WrapError(err)
		}
	}

	repoPaths, err := git_commands.GetRepoPaths(cmn.Fs, currentPath)
	if err != nil {
		return nil, errors.Errorf("Error getting repo paths: %v", err)
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
		syncMutex,
	), nil
}

func NewGitCommandAux(
	cmn *common.Common,
	version *git_commands.GitVersion,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
	repoPaths *git_commands.RepoPaths,
	repo *gogit.Repository,
	syncMutex *deadlock.Mutex,
) *GitCommand {
	cmd := NewGitCmdObjBuilder(cmn.Log, osCommand.Cmd)

	// here we're doing a bunch of dependency injection for each of our commands structs.
	// This is admittedly messy, but allows us to test each command struct in isolation,
	// and allows for better namespacing when compared to having every method living
	// on the one struct.
	// common ones are: cmn, osCommand, dotGitDir, configCommands
	configCommands := git_commands.NewConfigCommands(cmn, gitConfig, repo)

	gitCommon := git_commands.NewGitCommon(cmn, version, cmd, osCommand, repoPaths, repo, configCommands, syncMutex)

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
			// TODO: make patch builder take Gui.IgnoreWhitespaceInDiffView into
			// account. For now we just pass false.
			return workingTreeCommands.ShowFileDiff(from, to, reverse, filename, plain, false)
		})
	patchCommands := git_commands.NewPatchCommands(gitCommon, rebaseCommands, commitCommands, statusCommands, stashCommands, patchBuilder)
	bisectCommands := git_commands.NewBisectCommands(gitCommon)
	worktreeCommands := git_commands.NewWorktreeCommands(gitCommon)

	branchLoader := git_commands.NewBranchLoader(cmn, cmd, branchCommands.CurrentBranchInfo, configCommands)
	commitFileLoader := git_commands.NewCommitFileLoader(cmn, cmd)
	commitLoader := git_commands.NewCommitLoader(cmn, cmd, statusCommands.RebaseMode, gitCommon)
	reflogCommitLoader := git_commands.NewReflogCommitLoader(cmn, cmd)
	remoteLoader := git_commands.NewRemoteLoader(cmn, cmd, repo.Remotes)
	worktreeLoader := git_commands.NewWorktreeLoader(gitCommon)
	stashLoader := git_commands.NewStashLoader(cmn, cmd)
	tagLoader := git_commands.NewTagLoader(cmn, cmd)

	return &GitCommand{
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

// this returns the root of the current worktree. So if you start lazygit from within
// a subdirectory of the worktree, it will start in the context of the root of that worktree
func findWorktreeRoot(fs afero.Fs, currentPath string) (string, error) {
	for {
		// we don't care if .git is a directory or a file: either is okay.
		_, err := fs.Stat(path.Join(currentPath, ".git"))

		if err == nil {
			return currentPath, nil
		}

		if !os.IsNotExist(err) {
			return "", utils.WrapError(err)
		}

		currentPath = path.Dir(currentPath)

		atRoot := currentPath == path.Dir(currentPath)
		if atRoot {
			// we should never really land here: the code that creates GitCommand should
			// verify we're in a git directory
			return "", errors.New("Must open lazygit in a git repository")
		}
	}
}

func VerifyInGitRepo(osCommand *oscommands.OSCommand) error {
	return osCommand.Cmd.New(git_commands.NewGitCmd("rev-parse").Arg("--git-dir").ToArgv()).DontLog().Run()
}
