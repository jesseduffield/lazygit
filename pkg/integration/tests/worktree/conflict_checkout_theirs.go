package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var ConflictCheckoutTheirs = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Conflicting file can be resolved by checkout out their version",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFile(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Press(keys.Files.OpenMergeTool).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Merge conflict options")).
					Select(Contains("Checkout theirs")).
					Confirm()
			}).
			Press(keys.Universal.Refresh).
			Tap(func() {
				t.Common().ContinueOnConflictsResolved()
			}).
			IsEmpty()

		// t.Views().Files().
		// 	IsEmpty()
	},
})
