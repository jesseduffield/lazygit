package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesDisabled = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Disable stacked pushes and push only the current branch, retaining force-push confirmation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
		config.GetUserConfig().Git.EnableStackedPush = false
	},
	SetupRepo: func(shell *Shell) {
		createRebasedStackOfBranches(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

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
