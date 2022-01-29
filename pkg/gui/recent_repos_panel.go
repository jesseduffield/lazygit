package gui

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateRecentReposMenu() error {
	recentRepoPaths := gui.c.GetAppState().RecentRepos
	reposCount := utils.Min(len(recentRepoPaths), 20)

	// we won't show the current repo hence the -1
	menuItems := make([]*types.MenuItem, reposCount-1)
	for i, path := range recentRepoPaths[1:reposCount] {
		path := path // cos we're closing over the loop variable
		menuItems[i] = &types.MenuItem{
			DisplayStrings: []string{
				filepath.Base(path),
				style.FgMagenta.Sprint(path),
			},
			OnPress: func() error {
				// if we were in a submodule, we want to forget about that stack of repos
				// so that hitting escape in the new repo does nothing
				gui.RepoPathStack.Clear()
				return gui.dispatchSwitchToRepo(path, false)
			},
		}
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: gui.c.Tr.RecentRepos, Items: menuItems})
}

func (gui *Gui) handleShowAllBranchLogs() error {
	cmdObj := gui.git.Branch.AllBranchesLogCmdObj()
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
			return gui.c.ErrorMsg(gui.c.Tr.ErrRepositoryMovedOrDeleted)
		}
		return err
	}

	if err := commands.VerifyInGitRepo(gui.OSCommand); err != nil {
		if err := os.Chdir(originalPath); err != nil {
			return err
		}

		return err
	}

	newGitCommand, err := commands.NewGitCommand(
		gui.Common,
		gui.OSCommand,
		git_config.NewStdCachedGitConfig(gui.Log),
		gui.Mutexes.SyncMutex,
	)
	if err != nil {
		return err
	}
	gui.git = newGitCommand

	// these two mutexes are used by our background goroutines (triggered via `gui.goEvery`. We don't want to
	// switch to a repo while one of these goroutines is in the process of updating something
	gui.Mutexes.SyncMutex.Lock()
	defer gui.Mutexes.SyncMutex.Unlock()

	gui.Mutexes.RefreshingFilesMutex.Lock()
	defer gui.Mutexes.RefreshingFilesMutex.Unlock()

	if err := gui.recordCurrentDirectory(); err != nil {
		return err
	}

	gui.resetState("", reuse)

	return nil
}

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu
func (gui *Gui) updateRecentRepoList() error {
	if gui.git.Status.IsBareRepo() {
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
			newRepos = append(newRepos, repo)
		} else {
			isNew = false
		}
	}
	return isNew, newRepos
}
