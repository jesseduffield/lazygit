package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RestoreDeletedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Restore a deleted local branch from the reflog",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base commit").
			NewBranch("feature").
			EmptyCommit("on feature").
			Checkout("master").
			EmptyCommit("on master").
			RunCommand([]string{"git", "branch", "-D", "feature"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			)

		t.Views().Branches().
			Press(keys.Branches.RestoreBranch).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Deleted branches")).
					ContainsLines(
						Contains("feature"),
					).
					Select(Contains("feature")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectToast(Contains("Restored branch 'feature'"))
			})

		t.Views().Branches().
			Lines(
				Contains("master").IsSelected(),
				Contains("feature"),
			)
	},
})
