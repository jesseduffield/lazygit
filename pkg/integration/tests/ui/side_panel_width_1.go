package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SidePanelWidth1 = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that setting sidePanelWidth to 1.0 doesn't crash when navigating to commit files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.SidePanelWidth = 1.0
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1.txt", "content")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().
			Press(keys.Universal.NextScreenMode).
			Press(keys.Universal.GoInto)

		t.Views().CommitFiles().
			Press(keys.Universal.GoInto)

		t.Views().PatchBuilding().Content(Contains("+content"))
	},
})
