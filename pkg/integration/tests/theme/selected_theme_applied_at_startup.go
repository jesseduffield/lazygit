package theme

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedThemeAppliedAtStartup = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The theme that was selected in an earlier session is applied at startup",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		writeThemeFile(cfg, "pink", `
gui:
  theme:
    branchColorPatterns:
      master: '#ff00ff'
`)
		selectThemeAtStartup(cfg, "pink")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			ContainsColoredText("#ff00ff", "master")
	},
})
