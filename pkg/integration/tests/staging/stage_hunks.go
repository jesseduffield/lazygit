package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageHunks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and unstage various hunks of a file in the staging panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1a\n2a\n3a\n4a\n5a\n6a\n7a\n8a")
		shell.Commit("one")

		shell.UpdateFile("file1", "1a\n2a\n3b\n4a\n5a\n6b\n7a\n8a")

		// hunk looks like:
		// diff --git a/file1 b/file1
		// index 3653080..a6388b6 100644
		// --- a/file1
		// +++ b/file1
		// @@ -1,6 +1,6 @@
		//  1a
		//  2a
		// -3a
		// +3b
		//  4a
		//  5a
		// -6a
		// +6b
		//  7a
		//  8a
		// \ No newline at end of file
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			SelectedLines(
				Contains("-3a"),
			).
			Press(keys.Universal.NextBlock).
			SelectedLines(
				Contains("-6a"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-6a"),
				Contains("+6b"),
			).
			// when in hunk mode, pressing up/down moves us up/down by a hunk
			SelectPreviousItem().
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			).
			SelectNextItem().
			SelectedLines(
				Contains("-6a"),
				Contains("+6b"),
			).
			// stage the second hunk
			PressPrimaryAction().
			ContainsLines(
				Contains("-3a"),
				Contains("+3b"),
			).
			Tap(func() {
				t.Views().StagingSecondary().
					ContainsLines(
						Contains("-6a"),
						Contains("+6b"),
					)
			}).
			Press(keys.Universal.TogglePanel)

		t.Views().StagingSecondary().
			IsFocused().
			// after toggling panel, we're back to only having selected a single line
			SelectedLines(
				Contains("-6a"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("+6b"),
			).
			PressPrimaryAction().
			IsEmpty()

		t.Views().Staging().
			IsFocused().
			SelectedLines(
				Contains("-3a"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains(`-3a`),
				Contains(`+3b`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			}).
			Content(DoesNotContain("-3a").DoesNotContain("+3b")).
			SelectedLines(
				Contains("-6a"),
				Contains("+6b"),
			)
	},
})
