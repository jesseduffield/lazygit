package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashMergeWithStagedChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to squash merge a branch while there are staged changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			EmptyCommit("one").
			NewBranch("feature-branch").
			CreateFileAndAdd("feature-file", "content").
			Commit("feature").
			Checkout("original-branch").
			CreateFileAndAdd("staged-file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("feature-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and commit")).
			Confirm()

		t.ExpectToast(Equals("Disabled: Cannot squash merge while there are staged changes"))

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and leave uncommitted")).
			Confirm()

		t.ExpectToast(Equals("Disabled: Cannot squash merge while there are staged changes"))
	},
})
