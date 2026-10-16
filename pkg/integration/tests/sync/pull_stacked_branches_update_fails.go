package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullStackedBranchesUpdateFails = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull a stack that was rebased on the remote, when updating the branches below the current one fails; the current branch is not pulled then",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)
		shell.Checkout("branch3")
		shell.SetConfig("pull.rebase", "true")

		// Somebody deleted branch1 on the remote, which we haven't fetched yet
		shell.RunCommand([]string{"git", "-C", "../origin", "branch", "-D", "branch1"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Pull)

		t.ExpectPopup().Menu().
			Title(Equals("Pull")).
			Select(Contains("Pull all these branches in addition to the current one")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("couldn't find remote ref")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ↓3↑3"),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑2"),
				Contains("master"),
			)
	},
})
