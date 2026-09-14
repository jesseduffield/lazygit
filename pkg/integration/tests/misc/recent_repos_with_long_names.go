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
		other, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{current, target, other}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		// this one is cloned while master is checked out, so only its path is
		// long enough to be truncated
		shell.CloneNonBare("other")
		shell.NewBranch("branch-with-a-name-that-is-too-long")
		shell.CloneNonBare("repo-with-a-name-that-is-far-too-long")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Lines(
				// The repos live in the test's own directory, so their paths are
				// long enough to lose their middle
				Contains("repo-with-a-name-that-is-far-… branch-with-a-name-that-is-to… ~/_results/mi…names/actual").IsSelected(),
				Contains("other                          master                         ~/_results/mi…names/actual"),
				Contains("Cancel"),
			).
			Tooltip(Contains("Repo:   repo-with-a-name-that-is-far-too-long\n" +
				"Branch: branch-with-a-name-that-is-too-long\n" +
				"Path:   ~/_results/misc/recent_repos_with_long_names/actual")).
			// The values are lined up behind the labels that are there, so a
			// tooltip that only spells out the path doesn't indent it as if a
			// "Branch:" label were in front of it
			Select(Contains("other")).
			Tooltip(Equals("Path: ~/_results/misc/recent_repos_with_long_names/actual")).
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
