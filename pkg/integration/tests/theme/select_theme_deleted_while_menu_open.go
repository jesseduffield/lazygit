package theme

import (
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectThemeDeletedWhileMenuOpen = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme whose file was deleted after the menu was opened shows an error and selects nothing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		themesDir := filepath.Join(config.ConfigDir(), "themes")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("pink")).
			Tap(func() {
				if err := os.Remove(filepath.Join(themesDir, "pink.yml")); err != nil {
					t.Fail(err.Error())
				}
			}).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Equals("Theme 'pink' not found in " + themesDir)).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master"),
			).
			DoesNotContainColoredText("#ff00ff", "master")
		t.FileSystem().
			PathNotPresent(filepath.Join(config.ConfigDir(), "selected_theme.yml"))
	},
})
