package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
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

	Loaders Loaders
}

type Loaders struct {
	Branches      *loaders.BranchLoader
	CommitFiles   *loaders.CommitFileLoader
	Commits       *loaders.CommitLoader
	Files         *loaders.FileLoader
	ReflogCommits *loaders.ReflogCommitLoader
	Remotes       *loaders.RemoteLoader
	Stash         *loaders.StashLoader
	Tags          *loaders.TagLoader
}

func NewGitCommand(
	cmn *common.Common,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
) (*GitCommand, error) {
	if err := navigateToRepoRootDirectory(os.Stat, os.Chdir); err != nil {
		return nil, err
	}

	repo, err := setupRepository(gogit.PlainOpen, cmn.Tr.GitconfigParseErr)
	if err != nil {
		return nil, err
	}

	dotGitDir, err := findDotGitDir(os.Stat, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	return NewGitCommandAux(
		cmn,
		osCommand,
		gitConfig,
		dotGitDir,
		repo,
	), nil
}

func NewGitCommandAux(
	cmn *common.Common,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
	dotGitDir string,
	repo *gogit.Repository,
) *GitCommand {
	cmd := NewGitCmdObjBuilder(cmn.Log, osCommand.Cmd)

	// here we're doing a bunch of dependency injection for each of our commands structs.
	// This is admittedly messy, but allows us to test each command struct in isolation,
	// and allows for better namespacing when compared to having every method living
	// on the one struct.
	// common ones are: cmn, osCommand, dotGitDir, configCommands
	configCommands := git_commands.NewConfigCommands(cmn, gitConfig, repo)
	gitCommon := git_commands.NewGitCommon(cmn, cmd, osCommand, dotGitDir, repo, configCommands)

	statusCommands := git_commands.NewStatusCommands(gitCommon)
	fileLoader := loaders.NewFileLoader(cmn, cmd, configCommands)
	flowCommands := git_commands.NewFlowCommands(gitCommon)
	remoteCommands := git_commands.NewRemoteCommands(gitCommon)
	branchCommands := git_commands.NewBranchCommands(gitCommon)
	syncCommands := git_commands.NewSyncCommands(gitCommon)
	tagCommands := git_commands.NewTagCommands(gitCommon)
	commitCommands := git_commands.NewCommitCommands(gitCommon)
	customCommands := git_commands.NewCustomCommands(gitCommon)
	fileCommands := git_commands.NewFileCommands(gitCommon)
	submoduleCommands := git_commands.NewSubmoduleCommands(gitCommon)
	workingTreeCommands := git_commands.NewWorkingTreeCommands(gitCommon, submoduleCommands, fileLoader)
	rebaseCommands := git_commands.NewRebaseCommands(gitCommon, commitCommands, workingTreeCommands)
	stashCommands := git_commands.NewStashCommands(gitCommon, fileLoader, workingTreeCommands)
	// TODO: have patch manager take workingTreeCommands in its entirety
	patchManager := patch.NewPatchManager(cmn.Log, workingTreeCommands.ApplyPatch, workingTreeCommands.ShowFileDiff)
	patchCommands := git_commands.NewPatchCommands(gitCommon, rebaseCommands, commitCommands, statusCommands, stashCommands, patchManager)
	bisectCommands := git_commands.NewBisectCommands(gitCommon)

	return &GitCommand{
		Branch:      branchCommands,
		Commit:      commitCommands,
		Config:      configCommands,
		Custom:      customCommands,
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
		Loaders: Loaders{
			Branches:      loaders.NewBranchLoader(cmn, branchCommands.GetRawBranches, branchCommands.CurrentBranchName, configCommands),
			CommitFiles:   loaders.NewCommitFileLoader(cmn, cmd),
			Commits:       loaders.NewCommitLoader(cmn, cmd, dotGitDir, branchCommands.CurrentBranchName, statusCommands.RebaseMode),
			Files:         fileLoader,
			ReflogCommits: loaders.NewReflogCommitLoader(cmn, cmd),
			Remotes:       loaders.NewRemoteLoader(cmn, cmd, repo.Remotes),
			Stash:         loaders.NewStashLoader(cmn, cmd),
			Tags:          loaders.NewTagLoader(cmn, cmd),
		},
	}
}

func navigateToRepoRootDirectory(stat func(string) (os.FileInfo, error), chdir func(string) error) error {
	gitDir := env.GetGitDirEnv()
	if gitDir != "" {
		// we've been given the git directory explicitly so no need to navigate to it
		_, err := stat(gitDir)
		if err != nil {
			return utils.WrapError(err)
		}

		return nil
	}

	// we haven't been given the git dir explicitly so we assume it's in the current working directory as `.git/` (or an ancestor directory)

	for {
		_, err := stat(".git")

		if err == nil {
			return nil
		}

		if !os.IsNotExist(err) {
			return utils.WrapError(err)
		}

		if err = chdir(".."); err != nil {
			return utils.WrapError(err)
		}

		currentPath, err := os.Getwd()
		if err != nil {
			return err
		}

		atRoot := currentPath == filepath.Dir(currentPath)
		if atRoot {
			// we should never really land here: the code that creates GitCommand should
			// verify we're in a git directory
			return errors.New("Must open lazygit in a git repository")
		}
	}
}

// resolvePath takes a path containing a symlink and returns the true path
func resolvePath(path string) (string, error) {
	l, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if l.Mode()&os.ModeSymlink == 0 {
		return path, nil
	}

	return filepath.EvalSymlinks(path)
}

func setupRepository(openGitRepository func(string) (*gogit.Repository, error), gitConfigParseErrorStr string) (*gogit.Repository, error) {
	unresolvedPath := env.GetGitDirEnv()
	if unresolvedPath == "" {
		var err error
		unresolvedPath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	path, err := resolvePath(unresolvedPath)
	if err != nil {
		return nil, err
	}

	repository, err := openGitRepository(path)

	if err != nil {
		if strings.Contains(err.Error(), `unquoted '\' must be followed by new line`) {
			return nil, errors.New(gitConfigParseErrorStr)
		}

		return nil, err
	}

	return repository, err
}

func findDotGitDir(stat func(string) (os.FileInfo, error), readFile func(filename string) ([]byte, error)) (string, error) {
	if env.GetGitDirEnv() != "" {
		return env.GetGitDirEnv(), nil
	}

	f, err := stat(".git")
	if err != nil {
		return "", err
	}

	if f.IsDir() {
		return ".git", nil
	}

	fileBytes, err := readFile(".git")
	if err != nil {
		return "", err
	}
	fileContent := string(fileBytes)
	if !strings.HasPrefix(fileContent, "gitdir: ") {
		return "", errors.New(".git is a file which suggests we are in a submodule but the file's contents do not contain a gitdir pointing to the actual .git directory")
	}
	return strings.TrimSpace(strings.TrimPrefix(fileContent, "gitdir: ")), nil
}

func VerifyInGitRepo(osCommand *oscommands.OSCommand) error {
	return osCommand.Cmd.New("git rev-parse --git-dir").DontLog().Run()
}
