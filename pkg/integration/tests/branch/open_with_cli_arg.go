package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenWithCliArg = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open straight to branches panel using a CLI arg",
	ExtraCmdArgs: "branch",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().IsFocused()
	},
})
