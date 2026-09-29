package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectTheme = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a theme from the theme menu, then go back to no theme",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
		writeThemeFile(cfg, "blue", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#0000ff'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("master"),
			).
			DoesNotContainColoredText("#ff00ff", "master")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("(•) (none)").IsSelected(),
				Contains("( ) blue"),
				Contains("( ) pink"),
				Contains("Cancel"),
			).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")
		t.Views().Status().
			ContainsColoredText("#ff00ff", "master")
		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Equals("name: pink\n"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("( ) (none)").IsSelected(),
				Contains("( ) blue"),
				Contains("(•) pink"),
				Contains("Cancel"),
			).
			Confirm()

		t.ExpectToast(Equals("Theme: (none)"))

		t.Views().Branches().
			Lines(
				Contains("master"),
			).
			DoesNotContainColoredText("#ff00ff", "master")
		t.Views().Status().
			Content(Contains("master")).
			DoesNotContainColoredText("#ff00ff", "master")
		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Equals("name: \"\"\n"))
	},
})
