package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenStackOfBranchesWithLocalCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to fast-forward a stack of branches of which one has a commit of its own",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)

		shell.Checkout("branch2")
		shell.EmptyCommit("mine")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑3"),
				Contains("branch3 ↓3↑3"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			SelectNextItem().
			Press(keys.Branches.FastForward)

		t.ExpectPopup().Alert().Title(Equals("Error")).
			Content(Contains("Cannot fast-forward 'branch2' because it has commits")).
			Confirm()

		// None of them was touched, not even the ones we could have forwarded
		t.Views().Branches().
			Lines(
				Contains("master"),
				Contains("branch1 ↓1↑1").IsSelected(),
				Contains("branch2 ↓2↑3").IsSelected(),
				Contains("branch3 ↓3↑3").IsSelected(),
			)
	},
})
