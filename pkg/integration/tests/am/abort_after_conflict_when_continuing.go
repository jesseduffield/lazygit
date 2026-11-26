package am

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AbortAfterConflictWhenContinuing = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Continue applying patches with git am, and abort when the next patch conflicts too",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setUpConflictingAm(shell, true, true)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Shell().UpdateFileAndAdd("file", "resolved\n")

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Applying patches")).
			Select(Contains("continue")).
			Confirm()

		t.ExpectPopup().Menu().
			Title(Equals("Conflicts!")).
			Select(Contains("Abort applying patches")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Commits().
			Lines(
				Contains("conflicting change"),
				Contains("base"),
			)
	},
})
