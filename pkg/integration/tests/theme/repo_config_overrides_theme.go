package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RepoConfigOverridesTheme = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Settings in a repo's own config take precedence over the theme selected in the menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n      other: '#ff00ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("initial commit").
			NewBranch("other").
			Checkout("master")
		shell.CreateFile(".git/lazygit.yml", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#00ff00'\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			ContainsColoredText("#00ff00", "master")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "other").
			ContainsColoredText("#00ff00", "master").
			DoesNotContainColoredText("#ff00ff", "master")
	},
})
