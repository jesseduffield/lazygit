package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MissingThemeInMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selected theme whose file is missing is shown as such in the theme menu, which also says where theme files go",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		selectThemeAtStartup(cfg, "gone")
	},
	SetupRepo: func(shell *Shell) {
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		themesDir := filepath.Join(config.ConfigDir(), "themes")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			ContainsLines(
				Contains("No theme files found in"),
			).
			ContainsLines(
				Contains("( ) (none)").IsSelected(),
				Contains("(•) gone (not found)"),
				Contains("Cancel"),
			).
			Select(Contains("gone (not found)")).
			Tooltip(Contains("Disabled: Theme 'gone' not found in")).
			Confirm()

		t.ExpectToast(Equals("Disabled: Theme 'gone' not found in " + themesDir))

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Cancel()

		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Contains("name: gone"))
	},
})
