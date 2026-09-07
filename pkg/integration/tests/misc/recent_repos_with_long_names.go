package misc

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RecentReposWithLongNames = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Long repo and branch names are truncated in the recent repositories menu, and shown in full in the tooltip",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip: false,
	SetupConfig: func(cfg *config.AppConfig) {
		// the first entry is the repo we're in, so it isn't offered
		current, _ := filepath.Abs(".")
		target, _ := filepath.Abs("../repo-with-a-name-that-is-far-too-long")
		cfg.GetAppState().RecentRepos = []string{current, target}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("branch-with-a-name-that-is-too-long")
		shell.CloneNonBare("repo-with-a-name-that-is-far-too-long")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Lines(
				Contains("repo-with-a-name-that-is-far-… branch-with-a-name-that-is-to…").IsSelected(),
				Contains("Cancel"),
			).
			Tooltip(Equals("Repo:   repo-with-a-name-that-is-far-too-long\nBranch: branch-with-a-name-that-is-too-long")).
			// Filtering matches the full names, including the part that the
			// menu truncates
			Filter("too-long").
			Lines(
				Contains("repo-with-a-name-that-is-far-…").IsSelected(),
			).
			Confirm()

		t.Views().Status().Content(Contains("repo-with-a-name"))
	},
})
