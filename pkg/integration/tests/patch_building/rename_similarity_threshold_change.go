package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenameSimilarityThresholdChange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Changing the rename similarity threshold refreshes the commit files panel, but is disabled while building a patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("add original")

		shell.RenameFileInGit("original", "renamed")
		shell.UpdateFileAndAdd("renamed", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("change name and contents")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("change name and contents").IsSelected(),
				Contains("add original"),
			).
			PressEnter()

		// At the default threshold of 50% the 50%-similar change is not detected
		// as a rename.
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /"),
				Equals("  D original"),
				Equals("  A renamed"),
			).
			// Lowering the threshold turns it into a rename; the panel refreshes.
			Press(keys.Universal.DecreaseRenameSimilarityThreshold).
			Tap(func() {
				t.ExpectToast(Equals("Changed rename similarity threshold to 45%"))
			}).
			Lines(
				Equals("R original → renamed"),
			).
			// Start building a patch from the renamed file.
			PressPrimaryAction().
			Tap(func() {
				t.Views().Information().Content(Contains("Building patch"))

				// Changing the threshold is now disabled: the patch builder
				// can't cope with the rename turning into a delete and add.
				t.Views().CommitFiles().
					Press(keys.Universal.IncreaseRenameSimilarityThreshold)
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("Cannot change the rename similarity threshold while in patch building mode")).
					Confirm()
			}).
			// The file is unchanged: still a rename, still in the patch.
			Lines(
				Contains("original → renamed").IsSelected(),
			)
	},
})
