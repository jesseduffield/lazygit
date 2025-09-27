package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableSwitchTabWithPanelJumpKeys = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that the tab does not change by default when jumping to an already focused panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Focus().
			Press(keys.Universal.JumpToBlock[1])
		t.Views().Files().IsFocused().
			Press(keys.Universal.JumpToBlock[1])

		// Despite jumping to an already focused panel,
		// the tab should not change from the base files view
		t.Views().Files().IsFocused()
	},
})
