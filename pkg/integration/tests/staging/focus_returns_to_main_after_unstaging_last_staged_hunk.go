package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusReturnsToMainAfterUnstagingLastStagedHunk = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unstaging the last staged hunk from the secondary half collapses the split; focus returns to the main half on the now-unstaged change",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		shell.Commit("one")

		// One staged hunk...
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nnine\nten\neleven\n")
		// ...plus an unstaged change, so the main view splits into staged/unstaged.
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nSIX\nseven\neight\nnine\nten\neleven\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The unstaged half is focused first; switch to the staged half.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-six"),
				Contains("+SIX"),
			).
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			// Unstage the only staged hunk, emptying the staged side and collapsing the split.
			PressPrimaryAction()

		// Focus returns to the main half, which now shows the unstaged diff, landing
		// on the hunk we just unstaged.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			)
	},
})
