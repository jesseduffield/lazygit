package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectBrokenTheme = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme file that contains more than theme settings shows an error and keeps the current theme",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
		writeThemeFile(cfg, "broken", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#00ff00'\ngit:\n  autoFetch: false\n")
		selectThemeAtStartup(cfg, "pink")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("( ) (none)").IsSelected(),
				Contains("( ) broken"),
				Contains("(•) pink"),
				Contains("Cancel"),
			).
			Select(Contains("broken")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("broken.yml` couldn't be loaded")).
			Confirm()

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master").
			DoesNotContainColoredText("#00ff00", "master")
		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Contains("name: pink"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("( ) (none)").IsSelected(),
				Contains("( ) broken"),
				Contains("(•) pink"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
