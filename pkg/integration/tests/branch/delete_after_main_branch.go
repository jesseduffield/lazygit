package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteAfterMainBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a local branch after deleting one of the configured main branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("first commit").
			NewBranch("main").
			Checkout("master").
			NewBranch("dev").
			Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("main")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'main'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Lines(
				Contains("master"),
				Contains("dev").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'dev'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Lines(
				Contains("master").IsSelected(),
			)
	},
})
