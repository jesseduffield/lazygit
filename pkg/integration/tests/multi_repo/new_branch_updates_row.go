package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchUpdatesRow = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check out a new branch in the open repo and see its row in the repos panel follow",
	ExtraEnvVars: ceilingEnv,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    setupWorkspace,
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Contains("New branch name")).
			Type("topic").
			Confirm()

		t.Views().Repos().
			Lines(
				Contains("  * alpha").Contains("topic"),
				Contains("beta").Contains("master*"),
				Contains("org/gamma").Contains("master"),
			)
	},
})
