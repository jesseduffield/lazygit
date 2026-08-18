package gui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
)

// #5942: a repo whose git dir does not live at <work tree>/.git used to fall
// out of the recent-repos list — newRecentReposList kept an entry only when
// <path>/.git existed, so the next run in any other repo dropped it.
func TestNewRecentRepoLocationsListKeepsDotfileRepo(t *testing.T) {
	dir := t.TempDir()
	plainDir := t.TempDir()
	plainGit := filepath.Join(plainDir, ".git")
	if err := os.MkdirAll(plainGit, 0o755); err != nil {
		t.Fatal(err)
	}
	dotfileDir := t.TempDir() // no .git inside — the env carries the repo

	locations := newRecentRepoLocationsList(
		[]config.RecentRepoLocation{
			{Path: plainDir},
			{Path: dotfileDir, GitLocationEnvVars: []string{"GIT_DIR=/x/y", "GIT_WORK_TREE=" + dotfileDir}},
			{Path: filepath.Join(dir, "gone")}, // directory no longer exists: ages out
		},
		dir,
		nil,
	)

	paths := []string{}
	for _, location := range locations {
		paths = append(paths, location.Path)
	}
	if len(paths) != 3 || paths[0] != dir {
		t.Fatalf("unexpected head/got %v", paths)
	}
	foundDotfile := false
	for _, location := range locations {
		if location.Path == dotfileDir {
			foundDotfile = true
			if len(location.GitLocationEnvVars) != 2 {
				t.Fatalf("dotfile entry lost its env: %+v", location)
			}
		}
	}
	if !foundDotfile {
		t.Fatal("dotfile repo entry was dropped from the list")
	}
}

func TestNewRecentRepoLocationsListDropsLegacyPlainEntryWithoutDotGit(t *testing.T) {
	dir := t.TempDir()
	noGitDir := t.TempDir() // legacy entry with no env and no .git: was dropped before, still is

	locations := newRecentRepoLocationsList(
		[]config.RecentRepoLocation{{Path: noGitDir}},
		dir,
		nil,
	)

	if len(locations) != 1 || locations[0].Path != dir {
		t.Fatalf("legacy envless entry without .git must age out, got %+v", locations)
	}
}

func TestMigrateRecentReposPrependsCurrent(t *testing.T) {
	locations := migrateRecentRepos([]string{"/a", "/b"}, "/current")
	if len(locations) != 3 || locations[0].Path != "/current" {
		t.Fatalf("migration must prepend the current repo, got %+v", locations)
	}
	for _, location := range locations[1:] {
		if location.GitLocationEnvVars != nil {
			t.Fatalf("legacy entries carry no env, got %+v", location)
		}
	}
}
