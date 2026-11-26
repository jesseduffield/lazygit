package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectCurrentCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select the current bisect commit whenever it changes, but not on every refresh",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			StartBisect("HEAD", "HEAD~9")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("CI commit-05").Contains("<-- current")).
			NavigateToLine(Contains("CI commit-08")).
			Tap(func() {
				t.RefreshInBackground()
			}).
			SelectedLine(Contains("CI commit-08")).
			Tap(func() {
				t.Shell().RunCommand([]string{"git", "bisect", "good"})
				t.RefreshInBackground()
			}).
			SelectedLine(Contains("CI commit-07").Contains("<-- current"))
	},
})
