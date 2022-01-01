package commands

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// this takes something like:
// * (HEAD detached at 264fc6f5)
//	remotes
// and returns '264fc6f5' as the second match
const CurrentBranchNameRegex = `(?m)^\*.*?([^ ]*?)\)?$`

type Loaders struct {
	Commits       *loaders.CommitLoader
	Branches      *loaders.BranchLoader
	Files         *loaders.FileLoader
	CommitFiles   *loaders.CommitFileLoader
	Remotes       *loaders.RemoteLoader
	ReflogCommits *loaders.ReflogCommitLoader
	Stash         *loaders.StashLoader
	Tags          *loaders.TagLoader
}

// GitCommand is our main git interface
type GitCommand struct {
	*common.Common
	OSCommand            *oscommands.OSCommand
	Repo                 *gogit.Repository
	DotGitDir            string
	onSuccessfulContinue func() error
	PatchManager         *patch.PatchManager
	GitConfig            git_config.IGitConfig
	Loaders              Loaders

	// Push to current determines whether the user has configured to push to the remote branch of the same name as the current or not
	PushToCurrent bool

	// this is just a view that we write to when running certain commands.
	// Coincidentally at the moment it's the same view that OnRunCommand logs to
	// but that need not always be the case.
	GetCmdWriter func() io.Writer

	Cmd oscommands.ICmdObjBuilder

	Submodules SubmoduleCommands
}

// NewGitCommand it runs git commands
func NewGitCommand(
	cmn *common.Common,
	osCommand *oscommands.OSCommand,
	gitConfig git_config.IGitConfig,
) (*GitCommand, error) {
	var repo *gogit.Repository

	pushToCurrent := gitConfig.Get("push.default") == "current"

	if err := navigateToRepoRootDirectory(os.Stat, os.Chdir); err != nil {
		return nil, err
	}

	var err error
	if repo, err = setupRepository(gogit.PlainOpen, cmn.Tr.GitconfigParseErr); err != nil {
		return nil, err
	}

	dotGitDir, err := findDotGitDir(os.Stat, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	cmd := NewGitCmdObjBuilder(cmn.Log, osCommand.Cmd)

	gitCommand := &GitCommand{
		Common:        cmn,
		OSCommand:     osCommand,
		Repo:          repo,
		DotGitDir:     dotGitDir,
		PushToCurrent: pushToCurrent,
		GitConfig:     gitConfig,
		GetCmdWriter:  func() io.Writer { return ioutil.Discard },
		Cmd:           cmd,
	}

	gitCommand.Loaders = Loaders{
		Commits:       loaders.NewCommitLoader(cmn, gitCommand),
		Branches:      loaders.NewBranchLoader(cmn, gitCommand),
		Files:         loaders.NewFileLoader(cmn, cmd, gitConfig),
		CommitFiles:   loaders.NewCommitFileLoader(cmn, cmd),
		Remotes:       loaders.NewRemoteLoader(cmn, cmd, gitCommand.Repo.Remotes),
		ReflogCommits: loaders.NewReflogCommitLoader(cmn, cmd),
		Stash:         loaders.NewStashLoader(cmn, cmd),
		Tags:          loaders.NewTagLoader(cmn, cmd),
	}

	gitCommand.Submodules = NewSubmoduleCommands(cmn, cmd, dotGitDir)

	gitCommand.PatchManager = patch.NewPatchManager(gitCommand.Log, gitCommand.ApplyPatch, gitCommand.ShowFileDiff)

	return gitCommand, nil
}

func (c *GitCommand) WithSpan(span string) *GitCommand {
	// sometimes .WithSpan(span) will be called where span actually is empty, in
	// which case we don't need to log anything so we can just return early here
	// with the original struct
	if span == "" {
		return c
	}

	newGitCommand := &GitCommand{}
	*newGitCommand = *c
	newGitCommand.OSCommand = c.OSCommand.WithSpan(span)

	newGitCommand.Cmd = NewGitCmdObjBuilder(c.Log, newGitCommand.OSCommand.Cmd)

	// NOTE: unlike the other things here which create shallow clones, this will
	// actually update the PatchManager on the original struct to have the new span.
	// This means each time we call ApplyPatch in PatchManager, we need to ensure
	// we've called .WithSpan() ahead of time with the new span value
	newGitCommand.PatchManager.ApplyPatch = newGitCommand.ApplyPatch

	return newGitCommand
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
	return osCommand.Cmd.New("git rev-parse --git-dir").Run()
}

func (c *GitCommand) GetDotGitDir() string {
	return c.DotGitDir
}

func (c *GitCommand) GetCmd() oscommands.ICmdObjBuilder {
	return c.Cmd
}
