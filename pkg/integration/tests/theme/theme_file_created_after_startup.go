package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ThemeFileCreatedAfterStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selected theme whose file only appears after startup is marked as not loaded in the theme menu, and choosing it applies it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		selectThemeAtStartup(cfg, "pink")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Without a focus event, nothing loads the new file
		writeThemeFileWhileRunning(t, "pink", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")

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
				Contains("( ) (none)").IsSelected(),
				Contains("(•) pink (not loaded)"),
				Contains("Cancel"),
			).
			Select(Contains("pink (not loaded)")).
			Confirm()

		t.ExpectToast(Equals("Theme: pink"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")
	},
})
