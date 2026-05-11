package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Suggestions = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checking out a branch with name suggestions",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("my commit message").
			NewBranch("new-branch").
			NewBranch("new-branch-2").
			NewBranch("new-branch-3").
			NewBranch("branch-to-checkout").
			NewBranch("other-new-branch-2").
			NewBranch("other-new-branch-3")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Branches.CheckoutBranchByName)

		// we expect the first suggestion to be the branch we want because it most
		// closely matches what we typed in
		t.ExpectPopup().Prompt().
			Title(Equals("Branch name:")).
			Type("branch-to").
			SuggestionTopLines(Contains("branch-to-checkout")).
			ConfirmFirstSuggestion()

		t.Git().CurrentBranchName("branch-to-checkout")
	},
})
