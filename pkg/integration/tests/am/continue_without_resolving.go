package am

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ContinueWithoutResolving = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Continue applying patches with git am after a patch didn't apply, without doing anything about it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		setUpConflictingAm(shell, false, false)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Content(Equals("(applying patches) repo → master"))

		t.Views().Files().IsEmpty()

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Applying patches")).
			Select(Contains("continue")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("No changes - did you forget to use 'git add'?")).
			Confirm()

		t.Views().Status().Content(Equals("(applying patches) repo → master"))
	},
})
