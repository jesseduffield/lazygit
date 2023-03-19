package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/slices"
	appTypes "github.com/jesseduffield/lazygit/pkg/app/types"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/icons"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) getCurrentBranch(path string) string {
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

	return gui.c.Tr.LcBranchUnknown
}

func (gui *Gui) handleCreateRecentReposMenu() error {
	// we'll show an empty panel if there are no recent repos
	recentRepoPaths := []string{}
	if len(gui.c.GetAppState().RecentRepos) > 0 {
		// we skip the first one because we're currently in it
		recentRepoPaths = gui.c.GetAppState().RecentRepos[1:]
	}

	currentBranches := sync.Map{}

	wg := sync.WaitGroup{}
	wg.Add(len(recentRepoPaths))

	for _, path := range recentRepoPaths {
		go func(path string) {
			defer wg.Done()
			currentBranches.Store(path, gui.getCurrentBranch(path))
		}(path)
	}

	wg.Wait()

	menuItems := slices.Map(recentRepoPaths, func(path string) *types.MenuItem {
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
				gui.RepoPathStack.Clear()
				return gui.dispatchSwitchToRepo(path, false)
			},
		}
	})

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.RecentRepos, Items: menuItems})
}

func (gui *Gui) handleShowAllBranchLogs() error {
	cmdObj := gui.git.Branch.AllBranchesLogCmdObj()
	task := types.NewRunPtyTask(cmdObj.GetCmd())

	return gui.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: gui.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title: gui.c.Tr.LogTitle,
			Task:  task,
		},
	})
}

func (gui *Gui) dispatchSwitchToRepo(path string, reuse bool) error {
	env.UnsetGitDirEnvs()
	originalPath, err := os.Getwd()
	if err != nil {
		return nil
	}

	if err := os.Chdir(path); err != nil {
		if os.IsNotExist(err) {
			return gui.c.ErrorMsg(gui.c.Tr.ErrRepositoryMovedOrDeleted)
		}
		return err
	}

	if err := commands.VerifyInGitRepo(gui.os); err != nil {
		if err := os.Chdir(originalPath); err != nil {
			return err
		}

		return err
	}

	if err := gui.recordCurrentDirectory(); err != nil {
		return err
	}

	// these two mutexes are used by our background goroutines (triggered via `gui.goEvery`. We don't want to
	// switch to a repo while one of these goroutines is in the process of updating something
	gui.Mutexes.SyncMutex.Lock()
	defer gui.Mutexes.SyncMutex.Unlock()

	gui.Mutexes.RefreshingFilesMutex.Lock()
	defer gui.Mutexes.RefreshingFilesMutex.Unlock()

	return gui.onNewRepo(appTypes.StartArgs{}, reuse)
}

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu
func (gui *Gui) updateRecentRepoList() error {
	isBareRepo, err := gui.git.Status.IsBareRepo()
	if err != nil {
		return err
	}

	if isBareRepo {
		// we could totally do this but it would require storing both the git-dir and the
		// worktree in our recent repos list, which is a change that would need to be
		// backwards compatible
		gui.c.Log.Info("Not appending bare repo to recent repo list")
		return nil
	}

	recentRepos := gui.c.GetAppState().RecentRepos
	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}
	known, recentRepos := newRecentReposList(recentRepos, currentRepo)
	gui.IsNewRepo = known
	gui.c.GetAppState().RecentRepos = recentRepos
	return gui.c.SaveAppState()
}

// newRecentReposList returns a new repo list with a new entry but only when it doesn't exist yet
func newRecentReposList(recentRepos []string, currentRepo string) (bool, []string) {
	isNew := true
	newRepos := []string{currentRepo}
	for _, repo := range recentRepos {
		if repo != currentRepo {
			if _, err := os.Stat(filepath.Join(repo, ".git")); err != nil {
				continue
			}
			newRepos = append(newRepos, repo)
		} else {
			isNew = false
		}
	}
	return isNew, newRepos
}
