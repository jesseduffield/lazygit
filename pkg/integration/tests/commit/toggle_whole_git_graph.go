package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleWholeGitGraph = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle showing the whole git graph from the commit log options menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithDate("base", "2026-01-01T00:00:00Z").
			NewBranch("other").
			EmptyCommitWithDate("other commit", "2026-01-01T00:01:00Z").
			Checkout("master").
			EmptyCommitWithDate("master commit", "2026-01-01T00:02:00Z")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		toggle := func() {
			t.Views().Commits().
				Press(keys.Commits.OpenLogMenu).
				Tap(func() {
					t.ExpectPopup().Menu().
						Title(Equals("Commit Log Options")).
						Select(Contains("Toggle show whole git graph")).
						Confirm()
				})
		}

		t.Views().Commits().
			Focus().
			Lines(
				Contains("master commit"),
				Contains("base"),
			)

		toggle()

		t.Views().Commits().
			Lines(
				Contains("master commit"),
				Contains("other commit"),
				Contains("base"),
			)

		toggle()

		t.Views().Commits().
			Lines(
				Contains("master commit"),
				Contains("base"),
			)
	},
})
