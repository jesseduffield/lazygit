package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewordWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reword a commit while the whole graph is shown and another branch's commit sits right below it",
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
			NavigateToLine(Contains("dev-b-1")).
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Reword commit")).
					Clear().
					Type("renamed dev-b-1").
					Confirm()
			})

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("dev-b")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("dev-b-2"),
				Contains("renamed dev-b-1"),
				Contains("base"),
			)
	},
})
