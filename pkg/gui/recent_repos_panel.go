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

func (r *recentRepo) GetDisplayStrings() []string {
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
		newGitCommand, err := commands.NewGitCommand(gui.Log, gui.OSCommand, gui.Tr)
		if err != nil {
			return err
		}
		gui.GitCommand = newGitCommand
		return gui.Errors.ErrSwitchRepo
	}

	return gui.createMenu(recentRepos, handleMenuPress)
}

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu
func (gui *Gui) updateRecentRepoList() error {
	recentRepos := gui.Config.GetAppState().RecentRepos
	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}
	gui.Config.GetAppState().RecentRepos = newRecentReposList(recentRepos, currentRepo)
	return gui.Config.SaveAppState()
}

func newRecentReposList(recentRepos []string, currentRepo string) []string {
	newRepos := []string{currentRepo}
	for _, repo := range recentRepos {
		if repo != currentRepo {
			newRepos = append(newRepos, repo)
		}
	}
	return newRepos
}
