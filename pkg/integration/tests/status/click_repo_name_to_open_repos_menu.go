package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickRepoNameToOpenReposMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on the repo name in the status side panel to open the recent repositories menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo:    func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Click(1, 0)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories"))
	},
})
