package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickWorkingTreeStateToOpenRebaseOptionsMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on the working tree state in the status side panel to open the rebase options menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.Edit)

		t.Views().Status().
			Content(Contains("(rebasing) repo")).
			Click(1, 0)

		t.ExpectPopup().Menu().Title(Equals("Rebase options"))
	},
})
