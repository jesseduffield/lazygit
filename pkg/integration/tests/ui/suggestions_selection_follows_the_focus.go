package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SuggestionsSelectionFollowsTheFocus = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The suggestions list only shows a selection while it, rather than the prompt, has the focus",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("one").
			NewBranch("branch-to-checkout")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Branches.CheckoutBranchByName)

		t.ExpectPopup().Prompt().
			Title(Equals("Branch name:")).
			Type("branch-to").
			SuggestionTopLines(Contains("branch-to-checkout"))

		t.Views().Suggestions().SelectionIsHidden()

		t.Views().Prompt().Press(keys.Universal.TogglePanel)
		t.Views().Suggestions().
			IsFocused().
			SelectionIsActive().
			Press(keys.Universal.TogglePanel)

		t.Views().Prompt().IsFocused()
		t.Views().Suggestions().SelectionIsHidden()

		t.Views().Prompt().Press(keys.Universal.Return)
		t.Views().Branches().IsFocused()
	},
})
