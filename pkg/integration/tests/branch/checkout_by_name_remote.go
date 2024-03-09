package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckoutByNameRemote = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a remote branch by name, both using the full name and the local name.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		// create an origin/foo remote branch
		shell.CloneIntoRemote("origin")
		shell.NewBranch("foo")
		shell.PushBranch("origin", "foo")
		// delete the local version of the branch because we need to test checking it out from scratch
		shell.Checkout("master")
		shell.DeleteBranch("foo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			// maximising window so that we can see the tracked branch
			Press(keys.Universal.NextScreenMode).
			Press(keys.Branches.CheckoutBranchByName).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Branch name:")).
					Type("foo").
					SuggestionLines(
						Contains("foo"),
						Contains("origin/foo"),
					).
					ConfirmFirstSuggestion()
			}).
			Lines(
				Contains("foo").
					// we have not checked out origin/foo...
					DoesNotContain("origin/foo").
					// ... but we are tracking it (formatted as '<remote> <branch>')
					Contains("origin foo"),
				Contains("master"),
			).
			Press(keys.Branches.CheckoutBranchByName).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Branch name:")).
					Type("origin/foo").
					SuggestionLines(
						Contains("origin/foo"),
					).
					ConfirmFirstSuggestion()
			}).
			Lines(
				Contains("HEAD detached at origin/foo"),
				Contains("foo"),
				Contains("master"),
			)
	},
})
