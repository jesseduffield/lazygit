package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectThemeRecolorsViews = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting a theme that sets the default text color and the author colors recolors the commits view, and the staging view that the theme menu is opened from",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "always"
		writeThemeFile(cfg, "vivid", "gui:\n  theme:\n    defaultFgColor:\n      - '#ff00ff'\n    authorColors:\n      '*': '#00ff00'\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "one\ntwo\nthree\n").
			Commit("add file").
			UpdateFile("file", "one\nTWO\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("add file"),
			).
			DoesNotContainColoredText("#00ff00", "CI").
			DoesNotContainColoredText("#00ff00", "○").
			DoesNotContainColoredText("#ff00ff", "add file")

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			Content(Contains(" three")).
			DoesNotContainColoredText("#ff00ff", "three").
			Press(keys.Universal.SelectTheme)

		t.ExpectPopup().Menu().
			Title(Equals("Select theme")).
			Select(Contains("vivid")).
			Confirm()

		t.ExpectToast(Equals("Theme: vivid"))

		t.Views().Staging().
			IsFocused().
			ContainsColoredText("#ff00ff", "three")
		t.Views().Commits().
			ContainsColoredText("#00ff00", "CI").
			ContainsColoredText("#00ff00", "○").
			ContainsColoredText("#ff00ff", "add file")
	},
})
