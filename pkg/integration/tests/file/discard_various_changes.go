package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardVariousChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding all possible permutations of changed files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		createAllPossiblePermutationsOfChangedFiles(shell)
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		type statusFile struct {
			status string
			label  string
		}

		discardOneByOne := func(files []statusFile) {
			for _, file := range files {
				t.Views().Files().
					IsFocused().
					SelectedLine(Contains(file.status + " " + file.label)).
					Press(keys.Universal.Remove)

				t.ExpectPopup().Menu().
					Title(Equals("Discard changes")).
					Select(Contains("Discard all changes")).
					Confirm()
			}
		}

		discardOneByOne([]statusFile{
			{status: "UA", label: "added-them-changed-us.txt"},
			{status: "AA", label: "both-added.txt"},
			{status: "DD", label: "both-deleted.txt"},
			{status: "UU", label: "both-modded.txt"},
			{status: "AU", label: "changed-them-added-us.txt"},
			{status: "UD", label: "deleted-them.txt"},
			{status: "DU", label: "deleted-us.txt"},
		})

		t.ExpectPopup().Confirmation().
			Title(Equals("Continue")).
			Content(Contains("All merge conflicts resolved. Continue?")).
			Cancel()

		discardOneByOne([]statusFile{
			{status: "AM", label: "added-changed.txt"},
			{status: "MD", label: "change-delete.txt"},
			{status: "D ", label: "delete-change.txt"},
			{status: "D ", label: "deleted-staged.txt"},
			{status: " D", label: "deleted.txt"},
			{status: "MM", label: "double-modded.txt"},
			{status: "M ", label: "modded-staged.txt"},
			{status: " M", label: "modded.txt"},
			{status: "A ", label: "new-staged.txt"},
			{status: "??", label: "new.txt"},
			// the menu title only includes the new file
			{status: "R ", label: "renamed.txt → renamed2.txt"},
		})

		t.Views().Files().IsEmpty()
	},
})
