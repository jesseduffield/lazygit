package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranchesCurrentWithoutUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a rebased stack of branches whose topmost branch has no upstream yet; the prompt for its upstream comes before the force-push confirmation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createRebasedStackOfBranches(shell)
		shell.NewBranch("branch4")
		shell.EmptyCommit("four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Menu().
			Title(Equals("Push")).
			ContainsLines(
				Contains("  branch3 ↓3↑3"),
				Contains("  branch2 ↓2↑2"),
				Contains("  branch1 ↓1↑1"),
			).
			Select(Contains("Push all these branches in addition to the current one")).
			Confirm()

		t.ExpectPopup().Prompt().
			Title(Equals("Enter upstream as '<remote> <branchname>'")).
			InitialText(Equals("origin branch4")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Contains("The following branches have diverged from their remote branches:").
				Contains("branch3").Contains("branch2").Contains("branch1").DoesNotContain("branch4")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch4 ✓"),
				Contains("branch1 ✓"),
				Contains("branch2 ✓"),
				Contains("branch3 ✓"),
				Contains("local-only"),
				Contains("master ✓"),
			)
	},
})
