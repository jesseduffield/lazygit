package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PullStackedBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pull a stack of branches that was rebased on the remote, updating the branches below the current one too",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)
		shell.Checkout("branch3")
		shell.SetConfig("pull.rebase", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("branch3 ↓3↑3"),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑2"),
				Contains("master"),
			)

		t.Views().Files().IsFocused().Press(keys.Universal.Pull)

		t.ExpectPopup().Menu().
			Title(Equals("Pull")).
			ContainsLines(
				Contains("  branch2 ↓2↑2"),
				Contains("  branch1 ↓1↑1"),
			).
			Select(Contains("Pull all these branches in addition to the current one")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("branch3 ✓"),
				Contains("branch1 ✓"),
				Contains("branch2 ✓"),
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
				Contains("three-rewritten"),
				Contains("two-rewritten"),
				Contains("one-rewritten"),
				Contains("base"),
			)
	},
})
