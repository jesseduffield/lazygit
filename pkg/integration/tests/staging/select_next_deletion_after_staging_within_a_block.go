package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectNextDeletionAfterStagingWithinABlock = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "In line mode, staging a deletion in the middle of a block of deletions advances to the next deletion, not back to the block's first line",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "keep1\nd1\nd2\nd3\nd4\nkeep2\n")
		shell.Commit("one")
		// A block of four consecutive deleted lines (they all share a new-file line
		// number, which is what used to make the reveal jump to the first one).
		shell.UpdateFile("file1", "keep1\nkeep2\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused().Press(keys.Universal.FocusMainView)

		t.Views().Main().IsFocused().
			SelectedLines(Contains("-d1")).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			SelectedLines(Contains("-d3")).
			PressPrimaryAction().
			SelectedLines(Contains("-d4"))
	},
})
