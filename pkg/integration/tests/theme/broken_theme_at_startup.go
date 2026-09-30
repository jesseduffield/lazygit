package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BrokenThemeAtStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selected theme that can't be loaded at startup is reported, and lazygit starts without it",
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
		shell.EmptyCommit("initial")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectToast(Equals("Couldn't load theme 'broken'"))

		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			DoesNotContainColoredText("#ff00ff", "master")
	},
})
