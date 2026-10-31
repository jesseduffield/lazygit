package bisect

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectFirstBadCommitOnStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select the first bad commit when starting with a finished bisect",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(10).
			StartBisect("HEAD~2", "HEAD~5").
			// Mark commit-07, the commit that git checked out to test. This
			// finishes the bisect, because commit-08 is bad.
			RunCommand([]string{"git", "bisect", "good"})
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsFocused().
			/* EXPECTED:
			SelectedLine(Contains("CI commit-08").Contains("<-- bad"))
			ACTUAL: */
			SelectedLine(Contains("CI commit-07"))
	},
})
