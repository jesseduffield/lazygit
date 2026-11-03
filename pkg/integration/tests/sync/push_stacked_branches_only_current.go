package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesOnlyCurrent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the menu shortcut to push only the current branch in a rebased stack",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createRebasedStackOfBranches(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Menu().
			Title(Equals("Push")).
			ContainsLines(Contains("o Push only 'branch3'"))

		t.Views().Menu().Press(config.Keybinding{"o"})

		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Equals("Your branch has diverged from the remote branch. Press <esc> to cancel, or <enter> to force push.")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ✓"),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑2"),
				Contains("local-only"),
				Contains("master ✓"),
			)
	},
})
