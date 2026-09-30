package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BrokenThemeSelectionAtStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A file with the selected theme that can't be parsed at startup is reported, and lazygit starts without a theme",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", `
gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
		// The unclosed flow sequence makes this invalid YAML, so the name of
		// the selected theme is unknown
		writeSelectedThemeFile(cfg, "name: [pink\n")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectToast(Equals("Couldn't load the selected theme"))

		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			DoesNotContainColoredText("#ff00ff", "master")
	},
})
