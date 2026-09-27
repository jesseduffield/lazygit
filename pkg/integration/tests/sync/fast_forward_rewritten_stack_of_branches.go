package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FastForwardRewrittenStackOfBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fast-forward a whole stack of branches whose upstream branches were rewritten",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		createStackRewrittenOnTheRemote(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("branch1 ↓1↑1"),
				Contains("branch2 ↓2↑2"),
				Contains("branch3 ↓3↑3"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			SelectNextItem().
			Lines(
				Contains("master"),
				Contains("branch1 ↓1↑1").IsSelected(),
				Contains("branch2 ↓2↑2").IsSelected(),
				Contains("branch3 ↓3↑3").IsSelected(),
			).
			Press(keys.Branches.FastForward).
			Lines(
				Contains("master"),
				Contains("branch1 ✓").IsSelected(),
				Contains("branch2 ✓").IsSelected(),
				Contains("branch3 ✓").IsSelected(),
			)
	},
})
