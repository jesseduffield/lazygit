package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartInCommitsPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focus the commits panel at startup when a bisect is in progress",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3).
			// The bisect is in progress, although git knows only a bad commit
			// so far
			RunCommand([]string{"git", "bisect", "start"}).
			RunCommand([]string{"git", "bisect", "bad"})
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("CI commit-03").Contains("<-- bad").IsSelected(),
				Contains("CI commit-02"),
				Contains("CI commit-01"),
			)
	},
})
