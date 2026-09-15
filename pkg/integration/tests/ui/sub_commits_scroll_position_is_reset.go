package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SubCommitsScrollPositionIsReset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Viewing the commits of a branch again after scrolling down starts at the top again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(40)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			OriginY(0).
			Press(keys.Universal.GotoBottom).
			OriginYAtLeast(1).
			PressEscape()

		t.Views().Branches().
			IsFocused().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLineIdx(0).
			OriginY(0)
	},
})
