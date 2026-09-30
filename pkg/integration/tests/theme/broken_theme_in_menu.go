package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BrokenThemeInMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selected theme that couldn't be loaded is marked as not loaded in the theme menu, with the error of the last attempt as its tooltip, and choosing it loads it again",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// A theme may only set colors, so nerdFontsVersion makes the whole
		// theme invalid, including its branch color
		writeThemeFile(cfg, "broken", `
gui:
  nerdFontsVersion: "3"
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
		selectThemeAtStartup(cfg, "broken")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectToast(Equals("Couldn't load theme 'broken'"))

		// Break the file in another way, so that choosing the theme again
		// fails with a different error than the one from startup
		writeThemeFileWhileRunning(t, "broken", `
gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
    authorColors:
`)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("( ) (none)").IsSelected(),
				Contains("(•) broken (not loaded)"),
				Contains("Cancel"),
			).
			Select(Contains("broken (not loaded)")).
			Tooltip(Contains("broken.yml` couldn't be loaded")).
			Tooltip(Contains("field nerdFontsVersion not found")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("broken.yml` couldn't be loaded").Contains("gui.theme.authorColors has no value")).
			Confirm()

		t.Views().Branches().
			Lines(
				Contains("master"),
			).
			DoesNotContainColoredText("#ff00ff", "master")
		t.FileSystem().
			FileContent(filepath.Join(config.ConfigDir(), "selected_theme.yml"), Equals("name: broken\n"))

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("broken (not loaded)")).
			Tooltip(Contains("gui.theme.authorColors has no value").DoesNotContain("nerdFontsVersion")).
			Tap(func() {
				writeThemeFileWhileRunning(t, "broken", "gui:\n  theme:\n    branchColorPatterns:\n      master: '#ff00ff'\n")
			}).
			Confirm()

		t.ExpectToast(Equals("Theme: broken"))

		t.Views().Branches().
			ContainsColoredText("#ff00ff", "master")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Lines(
				Contains("( ) (none)").IsSelected(),
				Contains("(•) broken").DoesNotContain("not loaded"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
