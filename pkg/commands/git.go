package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
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
	log                  *logrus.Entry
	oSCommand            *oscommands.OSCommand
	repo                 *gogit.Repository
	tr                   *i18n.TranslationSet
	config               config.AppConfigurer
	getGitConfigValue    func(string) (string, error)
	dotGitDir            string
	onSuccessfulContinue func() error

	// Push to current determines whether the user has configured to push to the remote branch of the same name as the current or not
	pushToCurrent bool

	promptUserForCredential func(CredentialKind) string
	handleCredentialError   func(error)
}

func (c *GitCommand) GetPushToCurrent() bool {
	return c.pushToCurrent
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Entry, osCommand *oscommands.OSCommand, tr *i18n.TranslationSet, config config.AppConfigurer) (*GitCommand, error) {
	var repo *gogit.Repository

	// see what our default push behaviour is
	output, err := osCommand.RunCommandWithOutput(
		BuildGitCmdObjFromStr("config --get push.default"),
	)
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
		log:               log,
		oSCommand:         osCommand,
		tr:                tr,
		repo:              repo,
		config:            config,
		getGitConfigValue: getGitConfigValue,
		dotGitDir:         dotGitDir,
		pushToCurrent:     pushToCurrent,
	}

	return gitCommand, nil
}

func (c *GitCommand) NewPatchManager() *patch.PatchManager {
	return patch.NewPatchManager(c.log, c.ShowFileDiff)
}

func (c *GitCommand) WithSpan(span string) IGitCommand {
	// sometimes .WithSpan(span) will be called where span actually is empty, in
	// which case we don't need to log anything so we can just return early here
	// with the original struct
	if span == "" {
		return c
	}

	newGitCommand := &GitCommand{}
	*newGitCommand = *c
	newGitCommand.oSCommand = c.GetOSCommand().WithSpan(span)

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
	return osCommand.RunExecutable(
		BuildGitCmdObjFromStr("rev-parse --git-dir"),
	)
}

func (c *GitCommand) RunExecutable(cmdObj ICmdObj) error {
	_, err := c.RunCommandWithOutput(cmdObj)
	return err
}

func (c *GitCommand) RunCommandWithOutput(cmdObj ICmdObj) (string, error) {
	// TODO: have this retry logic in other places we run the command
	waitTime := 50 * time.Millisecond
	retryCount := 5
	attempt := 0

	for {
		output, err := c.GetOSCommand().RunCommandWithOutput(cmdObj)
		if err != nil {
			// if we have an error based on the index lock, we should wait a bit and then retry
			if strings.Contains(output, ".git/index.lock") {
				c.log.Error(output)
				c.log.Info("index.lock prevented command from running. Retrying command after a small wait")
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

func (c *GitCommand) GetOSCommand() *oscommands.OSCommand {
	return c.oSCommand
}

func BuildGitCmdStr(command string, positionalArgs []string, kwArgs map[string]bool) string {
	parts := []string{command}

	if len(kwArgs) > 0 {
		args := make([]string, 0, len(kwArgs))
		for arg, include := range kwArgs {
			if include {
				args = append(args, arg)
			}
		}
		utils.SortAlphabeticalInPlace(args)

		parts = append(parts, args...)
	}

	if len(positionalArgs) > 0 {
		parts = append(parts, positionalArgs...)
	}

	parts = utils.ExcludeEmpty(parts)

	return strings.Join(parts, " ")
}

func BuildGitCmdObj(command string, positionalArgs []string, kwArgs map[string]bool) ICmdObj {
	return BuildGitCmdObjFromStr(BuildGitCmdStr(command, positionalArgs, kwArgs))
}

// returns a command object from a command string. Prepends the `git ` part itself so
// if you want to do `git diff` just pass `diff` as the cmdStr
func BuildGitCmdObjFromStr(cmdStr string) ICmdObj {
	cmdObj := oscommands.NewCmdObjFromStr(GitCmdStr() + " " + cmdStr)
	SetDefaultEnvVars(cmdObj)

	return cmdObj
}

func BuildGitCmdObjFromArgs(args []string) ICmdObj {
	cmdObj := oscommands.NewCmdObjFromArgs(append([]string{GitCmdStr()}, args...))
	SetDefaultEnvVars(cmdObj)

	return cmdObj
}

func GitInitCmd() ICmdObj {
	return BuildGitCmdObjFromStr("init")
}

func GitVersionCmd() ICmdObj {
	return BuildGitCmdObjFromStr("--version")
}

func SetDefaultEnvVars(cmdObj ICmdObj) {
	cmdObj.GetCmd().Env = os.Environ()
	DisableOptionalLocks(cmdObj)
}

func DisableOptionalLocks(cmdObj ICmdObj) {
	cmdObj.AddEnvVars("GIT_OPTIONAL_LOCKS=0")
}

func (c *GitCommand) SkipEditor(cmdObj ICmdObj) {
	lazyGitPath := c.GetOSCommand().GetLazygitPath()

	cmdObj.AddEnvVars(
		"LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY",
		"GIT_EDITOR="+lazyGitPath,
		"EDITOR="+lazyGitPath,
		"VISUAL="+lazyGitPath,
	)
}

func (c *GitCommand) AllBranchesCmdObj() ICmdObj {
	cmdStr := c.cleanCustomGitCmdStr(
		c.config.GetUserConfig().Git.AllBranchesLogCmd,
	)

	return BuildGitCmdObjFromStr(cmdStr)
}

func (c *GitCommand) cleanCustomGitCmdStr(cmdStr string) string {
	if strings.HasPrefix(cmdStr, "git ") {
		return GitCmdStr() + strings.TrimPrefix(cmdStr, "git")
	} else {
		return cmdStr
	}
}

// We may have use for centralising this e.g. so that we call a specific executable
// or so that we can prepend some flags for bare repos.
// TODO: make this a method on the GitCommand struct
func GitCmdStr() string {
	return "git"
}

// BuildShellCmdObj returns the pointer to a custom command
func (c *GitCommand) BuildShellCmdObj(command string) ICmdObj {
	return oscommands.NewCmdObjFromArgs([]string{c.oSCommand.Platform.Shell, c.oSCommand.Platform.ShellArg, command})
}

func (c *GitCommand) GenericAbortCmdObj() ICmdObj {
	return c.GenericMergeOrRebaseCmdObj("abort")
}

func (c *GitCommand) GenericContinueCmdObj() ICmdObj {
	return c.GenericMergeOrRebaseCmdObj("continue")
}

func (c *GitCommand) GenericMergeOrRebaseCmdObj(action string) ICmdObj {
	status := c.WorkingTreeState()
	switch status {
	case REBASE_MODE_REBASING:
		return BuildGitCmdObjFromStr(fmt.Sprintf("rebase --%s", action))
	case REBASE_MODE_MERGING:
		return BuildGitCmdObjFromStr(fmt.Sprintf("merge --%s", action))
	default:
		panic("expected rebase mode")
	}
}

func (c *GitCommand) RunGitCmdFromStr(cmdStr string) error {
	return c.RunGitCmdFromStr(cmdStr)
}
