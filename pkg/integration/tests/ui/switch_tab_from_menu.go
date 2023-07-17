package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwitchTabFromMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switch tab via the options menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().
			Press(keys.Universal.OptionMenuAlt1)

		t.ExpectPopup().Menu().Title(Equals("Keybindings")).
			// Looping back around to the end to side-step the worktrees view which is
			// only present on recent git versions
			Select(Contains("Previous tab")).
			Confirm()

		t.Views().Submodules().IsFocused()
	},
})
