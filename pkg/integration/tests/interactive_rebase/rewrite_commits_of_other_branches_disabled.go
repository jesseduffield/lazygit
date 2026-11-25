package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RewriteCommitsOfOtherBranchesDisabled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With the whole graph shown, commits of other branches can't be rewritten, and the HEAD commit can't be moved up past them",
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
		disabledReason := "Disabled: Only commits of the current branch can be changed; the others are only listed because the whole git graph is shown"

		t.Views().Commits().
			Focus().
			Lines(
				Contains("dev-a-2").IsSelected(),
				Contains("dev-b-2"),
				Contains("dev-b-1"),
				Contains("dev-a-1"),
				Contains("base"),
			).
			Press(keys.Commits.RenameCommit).
			Tap(func() { t.ExpectToast(Equals(disabledReason)) }).
			Press(keys.Commits.SquashDown).
			Tap(func() { t.ExpectToast(Equals(disabledReason)) }).
			Press(keys.Universal.Remove).
			Tap(func() { t.ExpectToast(Equals(disabledReason)) }).
			Press(keys.Universal.Edit).
			Tap(func() { t.ExpectToast(Equals(disabledReason)) }).
			NavigateToLine(Contains("dev-b-1")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Commits.MoveDownCommit).
			Tap(func() { t.ExpectToast(Equals(disabledReason)) }).
			Press(keys.Universal.Return).
			NavigateToLine(Contains("dev-b-2")).
			Press(keys.Commits.MoveUpCommit).
			Tap(func() { t.ExpectToast(Equals("Disabled: Cannot move any further")) }).
			Lines(
				Contains("dev-a-2"),
				Contains("dev-b-2").IsSelected(),
				Contains("dev-b-1"),
				Contains("dev-a-1"),
				Contains("base"),
			)
	},
})
