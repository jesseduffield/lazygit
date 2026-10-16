package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesOnlyCurrent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Decline pushing the branches below the current one in a rebased stack, pushing only the current branch",
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
			Select(Contains("Push only 'branch3'")).
			Confirm()

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
