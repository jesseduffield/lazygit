package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageHunks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and unstage various hunks of a file in the staging panel",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// need to be working with a few lines so that git perceives it as two separate hunks
		shell.CreateFileAndAdd("file1", "1a\n2a\n3a\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13a\n14a\n15a")
		shell.Commit("one")

		shell.UpdateFile("file1", "1a\n2a\n3b\n4a\n5a\n6a\n7a\n8a\n9a\n10a\n11a\n12a\n13b\n14a\n15a")

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
		//  6a
		// @@ -10,6 +10,6 @@
		//  10a
		//  11a
		//  12a
		// -13a
		// +13b
		//  14a
		//  15a
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
			SelectedLine(Contains("-3a")).
			Press(keys.Universal.NextBlock).
			SelectedLine(Contains("-13a")).
			Press(keys.Main.ToggleSelectHunk).
			// when in hunk mode, pressing up/down moves us up/down by a hunk
			SelectPreviousItem().
			SelectedLine(Contains("-3a")).
			SelectNextItem().
			SelectedLine(Contains("-13a")).
			// stage the second hunk
			PressPrimaryAction().
			Content(Contains("-3a\n+3b")).
			Tap(func() {
				t.Views().StagingSecondary().
					Content(Contains("-13a\n+13b"))
			}).
			Press(keys.Universal.TogglePanel)

		t.Views().StagingSecondary().
			IsFocused().
			SelectedLine(Contains("-13a")).
			// after toggling panel, we're back to only having selected a single line
			PressPrimaryAction().
			SelectedLine(Contains("+13b")).
			PressPrimaryAction().
			IsEmpty()

		t.Views().Staging().
			IsFocused().
			SelectedLine(Contains("-3a")).
			Press(keys.Main.ToggleSelectHunk).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Actions().ConfirmDiscardLines()
			}).
			SelectedLine(Contains("-13a")).
			Content(DoesNotContain("-3a").DoesNotContain("+3b"))
	},
})
