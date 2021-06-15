package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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

// Git is our main git interface
type Git struct {
	*Commander
	*GitConfigMgr
	commitsMgr           *CommitsMgr
	branchesMgr          *BranchesMgr
	log                  *logrus.Entry
	os                   *oscommands.OS
	repo                 *gogit.Repository
	tr                   *i18n.TranslationSet
	config               config.AppConfigurer
	dotGitDir            string
	onSuccessfulContinue func() error

	promptUserForCredential func(CredentialKind) string
	handleCredentialError   func(error)
}

func (c *Git) GetLog() *logrus.Entry {
	return c.log
}

// NewGit it runs git commands
func NewGit(log *logrus.Entry, oS *oscommands.OS, tr *i18n.TranslationSet, config config.AppConfigurer) (*Git, error) {
	if err := navigateToRepoRootDirectory(os.Stat, os.Chdir); err != nil {
		return nil, err
	}

	repo, err := setupRepository(gogit.PlainOpen, tr.GitconfigParseErr)
	if err != nil {
		return nil, err
	}

	dotGitDir, err := findDotGitDir(os.Stat, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	commander := NewCommander(oS.RunWithOutput, log, oS.GetLazygitPath(), oS.Quote)
	gitConfig := NewGitConfigMgr(commander, config.GetUserConfig(), getGitConfigValue, log)
	commitsMgr := NewCommitsMgr(commander, gitConfig)
	branchesMgr := NewBranchesMgr(commander, gitConfig)

	gitCommand := &Git{
		Commander:    commander,
		GitConfigMgr: gitConfig,
		commitsMgr:   commitsMgr,
		branchesMgr:  branchesMgr,
		log:          log,
		os:           oS,
		tr:           tr,
		repo:         repo,
		config:       config,
		dotGitDir:    dotGitDir,
	}

	return gitCommand, nil
}

func (c *Git) Commits() ICommitsMgr {
	return c.commitsMgr
}

func (c *Git) Branches() IBranchesMgr {
	return c.branchesMgr
}

func (c *Git) Quote(str string) string {
	return c.os.Quote(str)
}

func (c *Git) NewPatchManager() *patch.PatchManager {
	return patch.NewPatchManager(c.log, c.ShowFileDiff)
}

func (c *Git) WithSpan(span string) IGit {
	// sometimes .WithSpan(span) will be called where span actually is empty, in
	// which case we don't need to log anything so we can just return early here
	// with the original struct
	if span == "" {
		return c
	}

	newGit := &Git{}
	*newGit = *c
	newGit.os = c.GetOS().WithSpan(span)

	return newGit
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
			// we should never really land here: the code that creates Git should
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

func VerifyInGitRepo(osCommand *oscommands.OS) error {
	return osCommand.Run(
		BuildGitCmdObjFromStr("rev-parse --git-dir"),
	)
}

func (c *Git) GetOS() *oscommands.OS {
	return c.os
}

func (c *Git) AllBranchesCmdObj() ICmdObj {
	cmdStr := c.cleanCustomGitCmdStr(
		c.config.GetUserConfig().Git.AllBranchesLogCmd,
	)

	return BuildGitCmdObjFromStr(cmdStr)
}

func (c *Git) cleanCustomGitCmdStr(cmdStr string) string {
	if strings.HasPrefix(cmdStr, "git ") {
		return GitCmdStr() + strings.TrimPrefix(cmdStr, "git")
	} else {
		return cmdStr
	}
}

func (c *Git) GenericAbortCmdObj() ICmdObj {
	return c.GenericMergeOrRebaseCmdObj("abort")
}

func (c *Git) GenericContinueCmdObj() ICmdObj {
	return c.GenericMergeOrRebaseCmdObj("continue")
}

func (c *Git) GenericMergeOrRebaseCmdObj(action string) ICmdObj {
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

func (c *Git) GetStatusFiles(opts loaders.LoadStatusFilesOpts) []*models.File {
	return loaders.NewStatusFileLoader(c).Load(opts)
}
