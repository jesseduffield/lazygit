package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardVariousChangesRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding all possible permutations of changed files via range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		createAllPossiblePermutationsOfChangedFiles(shell)
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  UA added-them-changed-us.txt"),
				Equals("  AA both-added.txt"),
				Equals("  DD both-deleted.txt"),
				Equals("  UU both-modded.txt"),
				Equals("  AU changed-them-added-us.txt"),
				Equals("  UD deleted-them.txt"),
				Equals("  DU deleted-us.txt"),
			).
			SelectNextItem().
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("deleted-us.txt")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Continue")).
					Content(Contains("All merge conflicts resolved. Continue the merge?")).
					Cancel()
			}).
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  AM added-changed.txt"),
				Equals("  MD change-delete.txt"),
				Equals("  D  delete-change.txt"),
				Equals("  D  deleted-staged.txt"),
				Equals("   D deleted.txt"),
				Equals("  MM double-modded.txt"),
				Equals("  M  modded-staged.txt"),
				Equals("   M modded.txt"),
				Equals("  A  new-staged.txt"),
				Equals("  ?? new.txt"),
				Equals("  R  renamed.txt → renamed2.txt"),
			).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("renamed.txt")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			})

		t.Views().Files().IsEmpty()
	},
})
