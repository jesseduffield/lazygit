package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/env"
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
			if strings.HasPrefix(content, refsPrefix) {
				// is a branch
				branchDisplay = strings.TrimPrefix(content, refsPrefix)
			} else {
				// detached HEAD state, displaying short SHA
				branchDisplay = utils.ShortSha(content)
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
				return self.c.ErrorMsg(errMsg)
			}
			return err
		}

		if err := commands.VerifyInGitRepo(self.c.OS()); err != nil {
			if err := os.Chdir(originalPath); err != nil {
				return err
			}

			return err
		}

		if err := self.recordDirectoryHelper.RecordCurrentDirectory(); err != nil {
			return err
		}

		self.c.Mutexes().RefreshingFilesMutex.Lock()
		defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

		return self.onNewRepo(appTypes.StartArgs{}, contextKey)
	})
}
