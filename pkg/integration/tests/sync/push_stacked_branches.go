package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushStackedBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a rebased stack of branches in one go, force-pushing all of them after a single confirmation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createRebasedStackOfBranches(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("branch3 ↓3↑3"),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑2"),
				Contains("local-only").DoesNotContain("↑"),
				Contains("master ✓"),
			)

		t.Views().Files().IsFocused().Press(keys.Universal.Push)

		t.ExpectPopup().Menu().
			Title(Equals("Push")).
			ContainsLines(
				Contains("  branch2 ↓2↑2"),
				Contains("  branch1 ↓1↑1"),
			).
			Select(Contains("Push all these branches in addition to the current one")).
			Confirm()

		t.ExpectPopup().Confirmation().
			Title(Equals("Force push")).
			Content(Contains("The following branches have diverged from their remote branches:").
				Contains("branch3").Contains("branch2").Contains("branch1")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ✓"),
				Contains("branch1 ✓"),
				Contains("branch2 ✓"),
				Contains("local-only").DoesNotContain("✓"),
				Contains("master ✓"),
			)

		t.Views().Remotes().
			Focus().
			Lines(
				Contains("origin"),
			).
			PressEnter()

		t.Views().RemoteBranches().
			IsFocused().
			NavigateToLine(Contains("branch1")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("one-rebased"),
				Contains("base"),
			)
	},
})
