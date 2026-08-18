package gui

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
)

// updateRecentRepoList registers the fact that we opened lazygit in this repo,
// so that we can open the same repo via the 'recent repos' menu.
//
// Bare repos are still skipped (there is no work tree to change back to), but
// a repo whose git dir does not live at <work tree>/.git — opened with
// --git-dir/--work-tree (yadm) or found through core.worktree (vcsh) — is now
// recorded with the environment it takes to reopen it (#5942). The plain-path
// RecentRepos list is still maintained in parallel with the repos git can
// find on its own, so an older lazygit reading the same state file keeps
// working.
func (gui *Gui) updateRecentRepoList() error {
	if gui.git.Status.IsBareRepo() {
		gui.c.Log.Info("Not appending bare repo to recent repo list")
		return nil
	}

	currentRepo, err := os.Getwd()
	if err != nil {
		return err
	}

	appState := gui.c.GetAppState()

	// Migration: an existing install has the old plain-path list. Fold it
	// into the richer one so those entries keep working (their env is empty,
	// which is correct for repos git finds from the work tree).
	if len(appState.RecentRepoLocations) == 0 && len(appState.RecentRepos) > 0 {
		appState.RecentRepoLocations = migrateRecentRepos(appState.RecentRepos, currentRepo)
	}

	// The env we must restore to reopen THIS repo: empty for ordinary repos.
	gitLocationEnvVars := env.GetGitLocationEnvVars()

	locations := newRecentRepoLocationsList(
		appState.RecentRepoLocations, currentRepo, gitLocationEnvVars)
	// The plain-path companion keeps only the repos git can find from the
	// path alone, matching what the old list could represent.
	plainRepos := []string{}
	for _, location := range locations {
		if len(location.GitLocationEnvVars) == 0 {
			plainRepos = append(plainRepos, location.Path)
		}
	}

	appState.RecentRepoLocations = locations
	appState.RecentRepos = plainRepos
	// TODO: migrate this file to use forward slashes on all OSes for consistency
	// (windows uses backslashes at the moment)
	return gui.c.SaveAppState()
}

// migrateRecentRepos converts a legacy plain-path list into locations,
// prepending the current repo (the list's head is the repo we are in).
func migrateRecentRepos(recentRepos []string, currentRepo string) []config.RecentRepoLocation {
	locations := []config.RecentRepoLocation{{Path: currentRepo}}
	for _, repo := range recentRepos {
		if repo == currentRepo {
			continue
		}
		locations = append(locations, config.RecentRepoLocation{Path: repo})
	}
	return locations
}

// newRecentRepoLocationsList returns a new location list with the current
// repo's entry refreshed, keeping the others in order. Entries whose
// directory no longer exists are dropped, so a deleted repo ages out of the
// menu. A dotfile repo's entry is kept even though <path>/.git does not
// exist: its GitLocationEnvVars are what make it openable again.
func newRecentRepoLocationsList(
	recentLocations []config.RecentRepoLocation,
	currentRepo string,
	gitLocationEnvVars []string,
) []config.RecentRepoLocation {
	newLocations := []config.RecentRepoLocation{{Path: currentRepo, GitLocationEnvVars: gitLocationEnvVars}}
	for _, location := range recentLocations {
		if location.Path == currentRepo {
			continue
		}
		if _, err := os.Stat(location.Path); err != nil {
			continue
		}
		if len(location.GitLocationEnvVars) == 0 {
			// A legacy entry that had no env recorded: it was only kept in
			// the old list when <path>/.git existed, so preserve that
			// contract for the plain entries.
			if _, err := os.Stat(filepath.Join(location.Path, ".git")); err != nil {
				continue
			}
		}
		newLocations = append(newLocations, location)
	}
	return newLocations
}

// recentRepoLocationsOrLegacy is the accessor other components use: it
// prefers the richer list and falls back to the legacy plain-path list when
// the state file was last written by an older version after a rollback.
func recentRepoLocationsOrLegacy(appState *config.AppState) []config.RecentRepoLocation {
	if len(appState.RecentRepoLocations) > 0 {
		return appState.RecentRepoLocations
	}
	locations := make([]config.RecentRepoLocation, 0, len(appState.RecentRepos))
	for _, repo := range appState.RecentRepos {
		locations = append(locations, config.RecentRepoLocation{Path: repo})
	}
	return locations
}


// RecentRepoMenuEntries is the menu's slice of the recent-repo list,
// skipping the head entry (the repo we are currently in).
func RecentRepoMenuEntries(appState *config.AppState) []config.RecentRepoLocation {
	locations := recentRepoLocationsOrLegacy(appState)
	if len(locations) > 0 {
		return locations[1:]
	}
	return nil
}
