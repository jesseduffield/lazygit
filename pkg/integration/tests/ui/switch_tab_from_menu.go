package ui

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
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
			Select(Contains("Next tab")).
			Confirm()

		t.Views().Worktrees().IsFocused()
	},
})
