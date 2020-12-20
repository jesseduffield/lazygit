package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

// this takes something like:
// * (HEAD detached at 264fc6f5)
//	remotes
// and returns '264fc6f5' as the second match
const CurrentBranchNameRegex = `(?m)^\*.*?([^ ]*?)\)?$`

// GitCommand is our main git interface
type GitCommand struct {
	Log                  *logrus.Entry
	OSCommand            *oscommands.OSCommand
	Repo                 *gogit.Repository
	Tr                   *i18n.TranslationSet
	Config               config.AppConfigurer
	getGitConfigValue    func(string) (string, error)
	removeFile           func(string) error
	DotGitDir            string
	onSuccessfulContinue func() error
	PatchManager         *patch.PatchManager

	// Push to current determines whether the user has configured to push to the remote branch of the same name as the current or not
	PushToCurrent bool
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Entry, osCommand *oscommands.OSCommand, tr *i18n.TranslationSet, config config.AppConfigurer) (*GitCommand, error) {
	var repo *gogit.Repository

	// see what our default push behaviour is
	output, err := osCommand.RunCommandWithOutput("git config --get push.default")
	pushToCurrent := false
	if err != nil {
		log.Errorf("error reading git config: %v", err)
	} else {
		pushToCurrent = strings.TrimSpace(output) == "current"
	}

	if err := verifyInGitRepo(osCommand.RunCommand); err != nil {
		return nil, err
	}

	if err := navigateToRepoRootDirectory(os.Stat, os.Chdir); err != nil {
		return nil, err
	}

	if repo, err = setupRepository(gogit.PlainOpen, tr.GitconfigParseErr); err != nil {
		return nil, err
	}

	dotGitDir, err := findDotGitDir(os.Stat, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	gitCommand := &GitCommand{
		Log:               log,
		OSCommand:         osCommand,
		Tr:                tr,
		Repo:              repo,
		Config:            config,
		getGitConfigValue: getGitConfigValue,
		removeFile:        os.RemoveAll,
		DotGitDir:         dotGitDir,
		PushToCurrent:     pushToCurrent,
	}

	gitCommand.PatchManager = patch.NewPatchManager(log, gitCommand.ApplyPatch, gitCommand.ShowFileDiff)

	return gitCommand, nil
}

func verifyInGitRepo(runCmd func(string, ...interface{}) error) error {
	return runCmd("git status")
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
