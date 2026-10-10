package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwitchRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switch to another repo from the repos panel",
	ExtraEnvVars: ceilingEnv,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    setupWorkspace,
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Repos().
			IsFocused().
			NavigateToLine(Contains("beta")).
			PressEnter()

		t.Views().Repos().
			IsFocused().
			Lines(
				DoesNotContain("  *").Contains("alpha"),
				Contains("  * beta").IsSelected(),
				DoesNotContain("  *").Contains("org/gamma"),
			)

		t.Views().Files().
			Lines(
				Contains("untracked.txt"),
			)
	},
})
