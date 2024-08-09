package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowDivergenceFromBaseBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show divergence from base branch in the status panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowDivergenceFromBaseBranch = "arrowAndNumber"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
		shell.CloneIntoRemote("origin")
		shell.NewBranch("feature")
		shell.HardReset("HEAD^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.NextBlock)

		t.Views().Status().
			Content(Equals("↓1 repo → feature"))
	},
})
