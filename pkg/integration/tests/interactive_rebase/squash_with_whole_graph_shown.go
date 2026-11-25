package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash a commit while the whole graph is shown and another branch's commit sits between the current branch's commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.Order = "date-order"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		setupInterleavedSiblingBranches(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("dev-a-2"),
				Contains("dev-b-2"),
				Contains("dev-b-1"),
				Contains("dev-a-1"),
				Contains("base"),
			).
			NavigateToLine(Contains("dev-b-2")).
			Press(keys.Commits.SquashDown).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Squash")).
					Content(Contains("squash")).
					Confirm()
			})

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("dev-b")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("dev-b-1"),
				Contains("base"),
			)
	},
})
