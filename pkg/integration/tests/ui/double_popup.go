package ui

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DoublePopup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a popup from within another popup and assert you can escape back to the side panels",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			// arbitrarily bringing up a popup
			PressPrimaryAction()

		t.ExpectPopup().Alert().
			Title(Contains("Error")).
			Content(Contains("You have already checked out this branch"))

		t.GlobalPress(keys.Universal.OpenRecentRepos)

		t.ExpectPopup().Menu().Title(Contains("recent repositories")).Cancel()

		t.Views().Branches().IsFocused()

		t.Views().Files().Focus()
	},
})
