package config

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoteNamedStar = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Having a config remote.*",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			SetConfig("remote.*.prune", "true").
			CreateNCommits(2)
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// here we're just asserting that we haven't panicked upon starting lazygit
		t.Views().Commits().
			Lines(
				Anything(),
				Anything(),
			)
	},
})
