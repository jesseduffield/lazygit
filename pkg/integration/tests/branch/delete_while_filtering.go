package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Regression test for deleting the last branch in the unfiltered list while
// filtering is on. This used to cause a segfault.
var DeleteWhileFiltering = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a local branch while there's a filter in the branches panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetAppState().LocalBranchSortOrder = "alphabetic"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("branch1")
		shell.NewBranch("branch2")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("branch1"),
				Contains("branch2"),
			).
			FilterOrSearch("branch").
			Lines(
				Contains("branch1").IsSelected(),
				Contains("branch2"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'branch2'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Lines(
				Contains("branch1").IsSelected(),
			)
	},
})
