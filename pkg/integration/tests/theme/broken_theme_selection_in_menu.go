package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BrokenThemeSelectionInMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When the file with the selected theme couldn't be parsed, the theme menu says why, and selecting a theme replaces the file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
		// The unclosed flow sequence makes this invalid YAML
		writeSelectedThemeFile(cfg, "name: [pink\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectToast(Equals("Couldn't load the selected theme"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("(•) (none)").IsSelected(),
				Contains("( ) pink"),
				Contains("Cancel"),
			).
			Tooltip(Contains("selected_theme.yml` couldn't be parsed")).
			Select(Contains("pink")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")
		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Equals("name: pink\n"))
	},
})
