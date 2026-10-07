package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EmptyDirFallback = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start in a directory without repos and get the existing notARepository behavior",
	ExtraEnvVars: ceilingEnv,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().NotARepository = "create"
	},
	SetupRepo: func(shell *Shell) {
		shell.Chdir("..")
		shell.CreateDir("empty")
		shell.Chdir("empty")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused()

		t.Views().Status().
			Content(Contains("empty"))
	},
})
