package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteBranchAtCommitRangeWithMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a branch in a range of selected commits that contains a merge commit, which can't be dropped",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("feature")
		shell.EmptyCommit("feature commit")
		shell.Checkout("master")
		shell.EmptyCommit("two")
		shell.Merge("feature")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Merge branch 'feature'").IsSelected(),
				Contains("feature commit"),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Drop commits or delete branches")).
			Lines(
				Contains("Drop commits").IsSelected(),
				Contains("Delete branch 'feature'"),
				Contains("Cancel"),
			).
			Tooltip(Contains("Disabled: Dropping a merge commit requires a single selected item")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"1"})
			})

		t.ExpectPopup().Menu().
			Title(Equals("Delete branch 'feature'?")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"c"})
			})

		t.Views().Branches().
			Lines(
				Contains("master").IsSelected(),
			)
	},
})
