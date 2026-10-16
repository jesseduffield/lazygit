package misc

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RecentReposColumnWidths = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The directory column of the recent repositories menu gets the room that the short columns before it don't need",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip: false,
	SetupConfig: func(cfg *config.AppConfig) {
		// the first entry is the repo we're in, so it isn't offered
		current, _ := filepath.Abs(".")
		target, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{current, target}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneNonBare("other")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The name and the branch are short, so the directory is shown in full
		// even though it is longer than the third of the row it would get if
		// they each took their maximum width
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Lines(
				Contains("other  master ~/_results/misc/recent_repos_column_widths/actual").IsSelected(),
				Contains("Cancel"),
			).
			Tap(func() {
				// Nothing is truncated, so there is nothing to spell out
				t.Views().Tooltip().IsInvisible()
			}).
			Confirm()

		t.Views().Status().Content(Contains("other → master"))
	},
})
