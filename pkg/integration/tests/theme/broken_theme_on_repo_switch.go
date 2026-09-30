package theme

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BrokenThemeOnRepoSwitch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selected theme that can't be loaded is reported again after switching repos",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{otherRepo}
		writeThemeFile(cfg, "broken", `
gui:
  nerdFontsVersion: "3"
`)
		selectThemeAtStartup(cfg, "broken")
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.CloneNonBare("other")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectToast(Equals("Couldn't load theme 'broken'"))

		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other").IsSelected(),
				Contains("Cancel"),
			).
			Confirm()

		t.ExpectToast(Equals("Couldn't load theme 'broken'"))
		t.Views().Status().Content(Contains("other → master"))
	},
})
