package misc

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterRecentRepos = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching to a recent repository by typing part of its name",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip: false,
	SetupConfig: func(cfg *config.AppConfig) {
		// the first entry is the repo we're in, so it isn't offered
		current, _ := filepath.Abs(".")
		other, _ := filepath.Abs("../other")
		target, _ := filepath.Abs("../target")
		cfg.GetAppState().RecentRepos = []string{current, other, target}
	},
	SetupRepo: func(shell *Shell) {
		shell.CloneNonBare("other")
		shell.CloneNonBare("target")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Filter("target").
			Lines(Contains("target").IsSelected()).
			Confirm()

		t.Views().Status().Content(Contains("target → master"))
	},
})
