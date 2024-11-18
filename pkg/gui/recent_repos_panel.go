package gui

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

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

	originalRepos := gui.c.GetAppState().RecentRepos
	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}

	isNew, updatedRepos := newRecentReposList(originalRepos, currentRepo)

	setRecentRepos := func(repos []string) error {
		// TODO: migrate this file to use forward slashes on all OSes for consistency
		// (windows uses backslashes at the moment)
		gui.c.GetAppState().RecentRepos = repos

		return gui.c.SaveAppState()
	}

	denyRepo := func() error { return setRecentRepos(originalRepos) }
	acceptRepo := func() error { return setRecentRepos(updatedRepos) }

	if !isNew {
		return acceptRepo()
	}

	menuSection := types.MenuSection{
		Title:  "Policy",
		Column: 0,
	}

	getRadioState := func(policy config.RecentReposPolicy) bool {
		return gui.c.GetAppState().RecentReposPolicy == policy
	}

	// TODO: how can we make our radio selection without exiting the menu?
	// I want to select one then press confirm
	setRecentReposPolicy := func(policy config.RecentReposPolicy) error {
		gui.c.GetAppState().RecentReposPolicy = policy
		return gui.c.SaveAppState()
	}

	menuItems := []*types.MenuItem{
		{
			Label:   "Accept",
			OnPress: acceptRepo,
		},
		{
			Label:   "Deny",
			OnPress: denyRepo,
		},
		{
			Label:   "Per-repo confirmation",
			Widget:  types.MakeMenuRadioButton(getRadioState(config.RecentReposPolicyPerRepoConfirmation)),
			OnPress: func() error { return setRecentReposPolicy(config.RecentReposPolicyPerRepoConfirmation) },
			Section: &menuSection,
		},
		{
			Label:   "Accept all",
			Widget:  types.MakeMenuRadioButton(getRadioState(config.RecentReposPolicyAcceptAll)),
			OnPress: func() error { return setRecentReposPolicy(config.RecentReposPolicyAcceptAll) },
			Section: &menuSection,
		},
		{
			Label:   "Reject all",
			Widget:  types.MakeMenuRadioButton(getRadioState(config.RecentReposPolicyRejectAll)),
			OnPress: func() error { return setRecentReposPolicy(config.RecentReposPolicyRejectAll) },
			Section: &menuSection,
		},
	}

	return gui.c.Menu(types.CreateMenuOptions{
		Title:      "Add to recent repos",
		Prompt:     "Repo: " + currentRepo,
		Items:      menuItems,
		HideCancel: true,
	})
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
