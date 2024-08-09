package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwitchTabWithPanelJumpKeys = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switch tab with the panel jump keys",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().Focus().
			Press(keys.Universal.JumpToBlock[2])

		t.Views().Branches().IsFocused().
			Press(keys.Universal.JumpToBlock[2])

		t.Views().Remotes().IsFocused().
			Press(keys.Universal.JumpToBlock[2])

		t.Views().Tags().IsFocused().
			Press(keys.Universal.JumpToBlock[2])

		t.Views().Branches().IsFocused().
			Press(keys.Universal.JumpToBlock[1])

		// When jumping to a panel from a different one, keep its current tab:
		t.Views().Worktrees().IsFocused()
	},
})
