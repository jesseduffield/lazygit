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
				Contains("UA").Contains("added-them-changed-us.txt").IsSelected(),
				Contains("AA").Contains("both-added.txt"),
				Contains("DD").Contains("both-deleted.txt"),
				Contains("UU").Contains("both-modded.txt"),
				Contains("AU").Contains("changed-them-added-us.txt"),
				Contains("UD").Contains("deleted-them.txt"),
				Contains("DU").Contains("deleted-us.txt"),
			).
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
					Content(Contains("All merge conflicts resolved. Continue?")).
					Cancel()
			}).
			Lines(
				Contains("AM").Contains("added-changed.txt").IsSelected(),
				Contains("MD").Contains("change-delete.txt"),
				Contains("D ").Contains("delete-change.txt"),
				Contains("D ").Contains("deleted-staged.txt"),
				Contains(" D").Contains("deleted.txt"),
				Contains("MM").Contains("double-modded.txt"),
				Contains("M ").Contains("modded-staged.txt"),
				Contains(" M").Contains("modded.txt"),
				Contains("A ").Contains("new-staged.txt"),
				Contains("??").Contains("new.txt"),
				Contains("R ").Contains("renamed.txt → renamed2.txt"),
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
