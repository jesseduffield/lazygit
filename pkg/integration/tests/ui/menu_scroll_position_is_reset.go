package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MenuScrollPositionIsReset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A menu that is opened after a scrolled down one starts at the top again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.Views().Menu().
			IsFocused().
			// The first line is a section header, so the first item is at index 1
			SelectedLineIdx(1).
			OriginY(0).
			Press(keys.Universal.GotoBottom).
			OriginYAtLeast(1).
			PressEscape()

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.Views().Menu().
			IsFocused().
			SelectedLineIdx(1).
			OriginY(0)
	},
})
