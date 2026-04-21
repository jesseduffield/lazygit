package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/direnv"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type onNewRepoFn func(startArgs appTypes.StartArgs, contextKey types.ContextKey) error

// helps switch back and forth between repos
type ReposHelper struct {
	c                     *HelperCommon
	recordDirectoryHelper *RecordDirectoryHelper
	onNewRepo             onNewRepoFn
}

func NewRecentReposHelper(
	c *HelperCommon,
	recordDirectoryHelper *RecordDirectoryHelper,
	onNewRepo onNewRepoFn,
) *ReposHelper {
	return &ReposHelper{
		c:                     c,
		recordDirectoryHelper: recordDirectoryHelper,
		onNewRepo:             onNewRepo,
	}
}

func (self *ReposHelper) EnterSubmodule(submodule *models.SubmoduleConfig) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	self.c.State().GetRepoPathStack().Push(wd)

	return self.DispatchSwitchToRepo(submodule.FullPath(), context.NO_CONTEXT)
}

func (self *ReposHelper) getCurrentBranch(path string) string {
	readHeadFile := func(path string) (string, error) {
		headFile, err := os.ReadFile(filepath.Join(path, "HEAD"))
		if err == nil {
			content := strings.TrimSpace(string(headFile))
			refsPrefix := "ref: refs/heads/"
			var branchDisplay string
			if bareName, ok := strings.CutPrefix(content, refsPrefix); ok {
				// is a branch
				branchDisplay = bareName
			} else {
				// detached HEAD state, displaying short hash
				branchDisplay = utils.ShortHash(content)
			}
			return branchDisplay, nil
		}
		return "", err
	}

	gitDirPath := filepath.Join(path, ".git")

	if gitDir, err := os.Stat(gitDirPath); err == nil {
		if gitDir.IsDir() {
			// ordinary repo
			if branch, err := readHeadFile(gitDirPath); err == nil {
				return branch
			}
		} else {
			// worktree
			if worktreeGitDir, err := os.ReadFile(gitDirPath); err == nil {
				content := strings.TrimSpace(string(worktreeGitDir))
				worktreePath := strings.TrimPrefix(content, "gitdir: ")
				if branch, err := readHeadFile(worktreePath); err == nil {
					return branch
				}
			}
		}
	}

	return self.c.Tr.BranchUnknown
}

func (self *ReposHelper) CreateRecentReposMenu() error {
	// we'll show an empty panel if there are no recent repos
	recentRepoPaths := []string{}
	if len(self.c.GetAppState().RecentRepos) > 0 {
		// we skip the first one because we're currently in it
		recentRepoPaths = self.c.GetAppState().RecentRepos[1:]
	}

	currentBranches := sync.Map{}

	wg := sync.WaitGroup{}
	wg.Add(len(recentRepoPaths))

	for _, path := range recentRepoPaths {
		go func(path string) {
			defer wg.Done()
			currentBranches.Store(path, self.getCurrentBranch(path))
		}(path)
	}

	wg.Wait()

	menuItems := lo.Map(recentRepoPaths, func(path string, _ int) *types.MenuItem {
		branchName, _ := currentBranches.Load(path)
		if icons.IsIconEnabled() {
			branchName = icons.BRANCH_ICON + " " + fmt.Sprintf("%v", branchName)
		}

		return &types.MenuItem{
			LabelColumns: []string{
				filepath.Base(path),
				style.FgCyan.Sprint(branchName),
				style.FgMagenta.Sprint(path),
			},
			OnPress: func() error {
				// if we were in a submodule, we want to forget about that stack of repos
				// so that hitting escape in the new repo does nothing
				self.c.State().GetRepoPathStack().Clear()
				return self.DispatchSwitchToRepo(path, context.NO_CONTEXT)
			},
		}
	})

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.RecentRepos, Items: menuItems})
}

func (self *ReposHelper) DispatchSwitchToRepo(path string, contextKey types.ContextKey) error {
	return self.DispatchSwitchTo(path, self.c.Tr.ErrRepositoryMovedOrDeleted, contextKey)
}

func (self *ReposHelper) DispatchSwitchTo(path string, errMsg string, contextKey types.ContextKey) error {
	return self.c.WithWaitingStatus(self.c.Tr.Switching, func(gocui.Task) error {
		env.UnsetGitLocationEnvVars()
		originalPath, err := os.Getwd()
		if err != nil {
			return nil
		}

		msg := utils.ResolvePlaceholderString(self.c.Tr.ChangingDirectoryTo, map[string]string{"path": path})
		self.c.LogCommand(msg, false)

		if err := os.Chdir(path); err != nil {
			if os.IsNotExist(err) {
				return errors.New(errMsg)
			}
			return err
		}

		if err := commands.VerifyInGitRepo(self.c.OS()); err != nil {
			if err := os.Chdir(originalPath); err != nil {
				return err
			}

			return err
		}

		direnvResult := self.logDirenvResult(direnv.Load(self.c.OS().Cmd))

		if err := self.recordDirectoryHelper.RecordCurrentDirectory(); err != nil {
			self.c.Log.Errorf("error recording current directory: %v", err)
		}

		self.c.Mutexes().RefreshingFilesMutex.Lock()
		defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

		if err := self.onNewRepo(appTypes.StartArgs{}, contextKey); err != nil {
			return err
		}

		if direnvResult.Blocked {
			self.c.OnUIThread(func() error {
				self.promptDirenvApproval(direnvResult.EnvrcPath)
				return nil
			})
			return nil
		}

		return direnvResult.Err
	})
}

// logDirenvResult writes whatever direnv emitted to the command log and the
// debug log; both happen for every load attempt regardless of outcome.
func (self *ReposHelper) logDirenvResult(result direnv.LoadResult) direnv.LoadResult {
	if result.Message != "" {
		self.c.LogCommand(result.Message, false)
	}
	if result.Err != nil {
		self.c.Log.WithError(result.Err).Warn("direnv load failed")
	}
	return result
}

// promptDirenvApproval shows the user the contents of an unapproved .envrc
// and offers to run `direnv allow` for them. On confirm, we approve the
// file and re-run Load so the new env reaches subprocesses; on cancel we
// leave the env as-is (the previous repo's vars are already unloaded by
// the initial Load call, which is the correct state).
func (self *ReposHelper) promptDirenvApproval(envrcPath string) {
	content, err := os.ReadFile(envrcPath)
	if err != nil {
		self.c.Log.WithError(err).Warn("could not read .envrc for approval prompt")
		return
	}

	indented := "  " + strings.ReplaceAll(strings.TrimRight(string(content), "\n"), "\n", "\n  ")
	prompt := utils.ResolvePlaceholderString(self.c.Tr.DirenvApprovalPrompt, map[string]string{
		"confirmKey": self.c.UserConfig().Keybinding.Universal.Confirm.String(),
		"cancelKey":  self.c.UserConfig().Keybinding.Universal.Return.String(),
		"content":    indented,
	})

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DirenvApprovalTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			if err := direnv.Allow(self.c.OS().Cmd, envrcPath); err != nil {
				return err
			}
			return self.logDirenvResult(direnv.Load(self.c.OS().Cmd)).Err
		},
	})
}
