package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteBranchesAtCommitRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete all branches pointing at a range of selected commits from the Commits panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.NewBranch("bookmark1")
		shell.Checkout("master")
		shell.EmptyCommit("three")
		shell.EmptyCommit("four")
		shell.NewBranchFrom("bookmark2", "master~1")
		shell.EmptyCommit("bookmark2 commit")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("bookmark2 commit"),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			SelectNextItem().
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Menu().
			Title(Equals("Drop commits or delete branches")).
			Lines(
				Contains("Drop commits").IsSelected(),
				Contains("Delete branch 'bookmark1'"),
				Contains("Delete branch 'bookmark2'"),
				Contains("Delete branches"),
				Contains("Cancel"),
			).
			Select(Contains("Delete branches")).
			Tooltip(Equals("bookmark1, bookmark2")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"a"})
			})

		t.ExpectPopup().Menu().
			Title(Equals("Delete selected branches?")).
			Tap(func() {
				t.Views().Menu().Press(config.Keybinding{"c"})
			})

		t.ExpectPopup().Confirmation().
			Title(Equals("Force delete branch")).
			Content(Equals("Some of the selected branches are not fully merged. Are you sure you want to delete them?")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master").IsSelected(),
			)

		t.Views().Commits().
			IsFocused().
			Lines(
				Contains("four").DoesNotContain("bookmark"),
				Contains("three").DoesNotContain("bookmark"),
				Contains("two").DoesNotContain("bookmark"),
				Contains("one"),
			)
	},
})
