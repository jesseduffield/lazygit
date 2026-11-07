package am

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SkipConflictingPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Skip a patch that conflicts while applying patches with git am",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setUpConflictingAm(shell, false, true)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("(applying patches) repo → master"))

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Applying patches")).
			Select(Contains("skip")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Commits().
			Lines(
				Contains("add file2"),
				Contains("conflicting change"),
				Contains("base"),
			)
	},
})
