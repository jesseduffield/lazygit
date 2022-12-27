package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Delete = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to delete the checked out branch first (to no avail), and then delete another branch.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("blah").
			NewBranch("branch-one").
			NewBranch("branch-two")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToBranchesView()

		assert.CurrentView().Lines(
			MatchesRegexp(`\*.*branch-two`),
			MatchesRegexp(`branch-one`),
			MatchesRegexp(`master`),
		)

		input.Press(keys.Universal.Remove)
		input.Alert().Title(Equals("Error")).Content(Contains("You cannot delete the checked out branch!")).Confirm()

		input.NextItem()

		input.Press(keys.Universal.Remove)
		input.Confirmation().
			Title(Equals("Delete Branch")).
			Content(Contains("Are you sure you want to delete the branch 'branch-one'?")).
			Confirm()

		assert.CurrentView().Name("localBranches").
			Lines(
				MatchesRegexp(`\*.*branch-two`),
				MatchesRegexp(`master`).IsSelected(),
			)
	},
})
