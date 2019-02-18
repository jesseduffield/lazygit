package gui

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type recentRepo struct {
	path string
}

// GetDisplayStrings returns the path from a recent repo.
func (r *recentRepo) GetDisplayStrings(isFocused bool) []string {
	yellow := color.New(color.FgMagenta)
	base := filepath.Base(r.path)
	path := yellow.Sprint(r.path)
	return []string{base, path}
}

func (gui *Gui) handleCreateRecentReposMenu(g *gocui.Gui, v *gocui.View) error {
	recentRepoPaths := gui.Config.GetAppState().RecentRepos
	reposCount := utils.Min(len(recentRepoPaths), 20)
	// we won't show the current repo hence the -1
	recentRepos := make([]*recentRepo, reposCount-1)
	for i, path := range recentRepoPaths[1:reposCount] {
		recentRepos[i] = &recentRepo{path: path}
	}

	handleMenuPress := func(index int) error {
		repo := recentRepos[index]
		if err := os.Chdir(repo.path); err != nil {
			return err
		}
		newGitCommand, err := commands.NewGitCommand(gui.Log, gui.OSCommand, gui.Tr, gui.Config)
		if err != nil {
			return err
		}
		gui.GitCommand = newGitCommand
		return gui.Errors.ErrSwitchRepo
	}

	return gui.createMenu(gui.Tr.SLocalize("RecentRepos"), recentRepos, handleMenuPress)
}

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu
func (gui *Gui) updateRecentRepoList() error {
	recentRepos := gui.Config.GetAppState().RecentRepos
	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}
	known, recentRepos := newRecentReposList(recentRepos, currentRepo)
	gui.Config.SetIsNewRepo(known)
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
