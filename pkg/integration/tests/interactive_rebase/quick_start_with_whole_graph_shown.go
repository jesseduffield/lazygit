package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var QuickStartWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Quick-start an interactive rebase while the whole graph is shown, so that commits of another branch come first",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.Order = "date-order"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		setupInterleavedSiblingBranches(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("dev-a-2").IsSelected(),
				Contains("dev-b-2"),
				Contains("dev-b-1"),
				Contains("dev-a-1"),
				Contains("base"),
			).
			Press(keys.Commits.StartInteractiveRebase).
			TopLines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("dev-b-2"),
				Contains("pick").Contains("dev-b-1"),
				Contains("─── Commits"),
			)

		t.Common().AbortRebase()
	},
})
