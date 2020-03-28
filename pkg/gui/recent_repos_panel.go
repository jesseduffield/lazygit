package gui

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) handleCreateRecentReposMenu(g *gocui.Gui, v *gocui.View) error {
	recentRepoPaths := gui.Config.GetAppState().RecentRepos
	reposCount := utils.Min(len(recentRepoPaths), 20)
	yellow := color.New(color.FgMagenta)
	// we won't show the current repo hence the -1
	menuItems := make([]*menuItem, reposCount-1)
	for i, path := range recentRepoPaths[1:reposCount] {
		innerPath := path
		menuItems[i] = &menuItem{
			displayStrings: []string{
				filepath.Base(innerPath),
				yellow.Sprint(innerPath),
			},
			onPress: func() error {
				if err := os.Chdir(innerPath); err != nil {
					return err
				}
				newGitCommand, err := commands.NewGitCommand(gui.Log, gui.OSCommand, gui.Tr, gui.Config)
				if err != nil {
					return err
				}
				gui.GitCommand = newGitCommand
				gui.State.FilterPath = ""
				return gui.Errors.ErrSwitchRepo
			},
		}
	}

	return gui.createMenu(gui.Tr.SLocalize("RecentRepos"), menuItems, createMenuOptions{showCancel: true})
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
