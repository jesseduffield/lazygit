package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AdvanceToNextHunkAfterStagingShiftsLineNumbers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After staging a hunk that adds lines, the selection advances to the next hunk even though staging shifted the later hunks' old-side line numbers",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n7\n8\n")
		shell.Commit("one")

		// Three change blocks: a modification near the top, an inserted line in the
		// middle (which changes the line count), and a deletion-led modification below
		// it. The middle block is the one we stage.
		shell.UpdateFile("file1", "1\nX\n3\n4\nNEW\n5\n6\nY\n8\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The first block (2 -> X) is selected on focus; move down to the inserted line.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-2"),
				Contains("+X"),
			).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("+NEW"),
			).
			// Staging the inserted line bumps the old-side line numbers of the block
			// below it, so matching the next hunk by its old-side number would miss and
			// the selection would fall back to the earlier hunk.
			PressPrimaryAction().
			SelectedLines(
				Contains("-7"),
				Contains("+Y"),
			)
	},
})
