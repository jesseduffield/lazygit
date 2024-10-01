package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteRemoteBranchWithDifferentName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a remote branch that has a different name than the local branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("origin")
		shell.NewBranch("mybranch-local")
		shell.PushBranchAndSetUpstream("origin", "mybranch-local:mybranch-remote")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("mybranch-local ✓"),
			).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'mybranch-local'?")).
					Select(Contains("Delete remote branch")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete branch 'mybranch-remote'?")).
					Content(Equals("Are you sure you want to delete the remote branch 'mybranch-remote' from 'origin'?")).
					Confirm()
			}).
			Lines(
				Contains("master"),
				Contains("mybranch-local (upstream gone)").IsSelected(),
			)
	},
})
