package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StartInParent = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start in a directory that holds several repos and see them listed with their status",
	ExtraEnvVars: ceilingEnv,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    setupWorkspace,
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Repos().
			IsFocused().
			Lines(
				Contains("  * alpha").Contains("feature").IsSelected(),
				Contains("beta").Contains("master*"),
				Contains("org/gamma").Contains("master"),
			)

		t.Views().Files().
			IsEmpty()
	},
})
