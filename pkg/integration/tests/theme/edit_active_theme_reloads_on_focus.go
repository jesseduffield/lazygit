package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditActiveThemeReloadsOnFocus = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Editing the file of the theme selected in the menu and refocusing the window applies the edit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")

		writeThemeFileWhileRunning(t, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#00ff00'\n")
		t.FocusIn()

		t.Views().Branches().
			ContainsColoredText("#00ff00", "master")
	},
})
