package gui

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateRecentReposMenu() error {
	recentRepoPaths := gui.Config.GetAppState().RecentRepos
	reposCount := utils.Min(len(recentRepoPaths), 20)

	// we won't show the current repo hence the -1
	menuItems := make([]*menuItem, reposCount-1)
	for i, path := range recentRepoPaths[1:reposCount] {
		path := path // cos we're closing over the loop variable
		menuItems[i] = &menuItem{
			displayStrings: []string{
				filepath.Base(path),
				style.FgMagenta.Sprint(path),
			},
			onPress: func() error {
				// if we were in a submodule, we want to forget about that stack of repos
				// so that hitting escape in the new repo does nothing
				gui.RepoPathStack = []string{}
				return gui.dispatchSwitchToRepo(path, false)
			},
		}
	}

	return gui.createMenu(gui.Tr.RecentRepos, menuItems, createMenuOptions{showCancel: true})
}

func (gui *Gui) handleShowAllBranchLogs() error {
	cmdObj := gui.Git.Branch.AllBranchesLogCmdObj()
	task := NewRunPtyTask(cmdObj.GetCmd())

	return gui.refreshMainViews(refreshMainOpts{
		main: &viewUpdateOpts{
			title: "Log",
			task:  task,
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
			return gui.createErrorPanel(gui.Tr.ErrRepositoryMovedOrDeleted)
		}
		return err
	}

	if err := commands.VerifyInGitRepo(gui.OSCommand); err != nil {
		if err := os.Chdir(originalPath); err != nil {
			return err
		}

		return err
	}

	newGitCommand, err := commands.NewGitCommand(gui.Common, gui.OSCommand, git_config.NewStdCachedGitConfig(gui.Log))
	if err != nil {
		return err
	}
	gui.Git = newGitCommand

	// these two mutexes are used by our background goroutines (triggered via `gui.goEvery`. We don't want to
	// switch to a repo while one of these goroutines is in the process of updating something
	gui.Mutexes.FetchMutex.Lock()
	defer gui.Mutexes.FetchMutex.Unlock()

	gui.Mutexes.RefreshingFilesMutex.Lock()
	defer gui.Mutexes.RefreshingFilesMutex.Unlock()

	gui.resetState("", reuse)

	return nil
}

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu
func (gui *Gui) updateRecentRepoList() error {
	if gui.Git.Status.IsBareRepo() {
		// we could totally do this but it would require storing both the git-dir and the
		// worktree in our recent repos list, which is a change that would need to be
		// backwards compatible
		gui.Log.Info("Not appending bare repo to recent repo list")
		return nil
	}

	recentRepos := gui.Config.GetAppState().RecentRepos
	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}
	known, recentRepos := newRecentReposList(recentRepos, currentRepo)
	gui.IsNewRepo = known
	gui.Config.GetAppState().RecentRepos = recentRepos
	return gui.Config.SaveAppState()
}

// newRecentReposList returns a new repo list with a new entry but only when it doesn't exist yet
func newRecentReposList(recentRepos []string, currentRepo string) (bool, []string) {
	isNew := true
	newRepos := []string{currentRepo}
	for _, repo := range recentRepos {
		if repo != currentRepo {
			newRepos = append(newRepos, repo)
		} else {
			isNew = false
		}
	}
	return isNew, newRepos
}
