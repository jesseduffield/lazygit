package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashFixupsInCurrentBranchWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash the fixups of the current branch while the whole graph is shown, so that another branch's commits come first and in between",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.Order = "date-order"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		setupInterleavedSiblingBranches(shell)
		shell.EmptyCommitWithDate("fixup! dev-b-1", "2026-01-01T00:05:00Z")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("fixup! dev-b-1").IsSelected(),
				Contains("dev-a-2"),
				Contains("dev-b-2"),
				Contains("dev-b-1"),
				Contains("dev-a-1"),
				Contains("base"),
			).
			Press(keys.Commits.SquashAboveCommits).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Apply fixup commits")).
					Select(Contains("In current branch")).
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
				Contains("dev-b-1"),
				Contains("base"),
			)
	},
})
