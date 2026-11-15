package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CheckForUpdateInDevBuild = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checking for an update in a build that has no release version explains why it isn't possible",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus().
			Press(keys.Status.CheckForUpdate)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(
				Contains("This is a development build (").
					Contains("), so lazygit can't check for updates."),
			).
			Confirm()
	},
})
