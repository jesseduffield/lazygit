package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		DotGitDir:         dotGitDir,
		PushToCurrent:     pushToCurrent,
	}

	gitCommand.PatchManager = patch.NewPatchManager(log, gitCommand.ApplyPatch, gitCommand.ShowFileDiff)

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
	return osCommand.RunCommand("git rev-parse --git-dir")
}

func (c *GitCommand) RunCommand(formatString string, formatArgs ...interface{}) error {
	_, err := c.RunCommandWithOutput(formatString, formatArgs...)
	return err
}

func (c *GitCommand) RunCommandWithOutput(formatString string, formatArgs ...interface{}) (string, error) {
	// TODO: have this retry logic in other places we run the command
	waitTime := 50 * time.Millisecond
	retryCount := 5
	attempt := 0

	for {
		output, err := c.OSCommand.RunCommandWithOutput(formatString, formatArgs...)
		if err != nil {
			// if we have an error based on the index lock, we should wait a bit and then retry
			if strings.Contains(output, ".git/index.lock") {
				c.Log.Error(output)
				c.Log.Info("index.lock prevented command from running. Retrying command after a small wait")
				attempt++
				time.Sleep(waitTime)
				if attempt < retryCount {
					continue
				}
			}
		}
		return output, err
	}
}
