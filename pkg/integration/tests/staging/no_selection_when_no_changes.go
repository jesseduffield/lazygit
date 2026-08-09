package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoSelectionWhenNoChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view when there are no changes shows no selection, and navigating doesn't conjure one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			IsEmpty().
			Press(keys.Universal.FocusMainView)

		// There's nothing to act on, so the focused main view shows the placeholder with
		// no selection — and a navigation key just scrolls rather than conjuring one.
		t.Views().Main().
			IsFocused().
			Content(Contains("No changed files")).
			SelectionIsHidden().
			Press(keys.Universal.GotoTop)

		t.Views().Main().
			SelectionIsHidden()
	},
})
