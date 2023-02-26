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
			SelectedLines(
				Contains("-3a"),
			).
			Press(keys.Universal.NextBlock).
			SelectedLines(
				Contains("-13a"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("@@ -10,6 +10,6 @@"),
				Contains(" 10a"),
				Contains(" 11a"),
				Contains(" 12a"),
				Contains("-13a"),
				Contains("+13b"),
				Contains(" 14a"),
				Contains(" 15a"),
				Contains(`\ No newline at end of file`),
			).
			// when in hunk mode, pressing up/down moves us up/down by a hunk
			SelectPreviousItem().
			SelectedLines(
				Contains(`@@ -1,6 +1,6 @@`),
				Contains(` 1a`),
				Contains(` 2a`),
				Contains(`-3a`),
				Contains(`+3b`),
				Contains(` 4a`),
				Contains(` 5a`),
				Contains(` 6a`),
			).
			SelectNextItem().
			SelectedLines(
				Contains("@@ -10,6 +10,6 @@"),
				Contains(" 10a"),
				Contains(" 11a"),
				Contains(" 12a"),
				Contains("-13a"),
				Contains("+13b"),
				Contains(" 14a"),
				Contains(" 15a"),
				Contains(`\ No newline at end of file`),
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
						Contains("-13a"),
						Contains("+13b"),
					)
			}).
			Press(keys.Universal.TogglePanel)

		t.Views().StagingSecondary().
			IsFocused().
			// after toggling panel, we're back to only having selected a single line
			SelectedLines(
				Contains("-13a"),
			).
			PressPrimaryAction().
			SelectedLines(
				Contains("+13b"),
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
				Contains(`@@ -1,6 +1,6 @@`),
				Contains(` 1a`),
				Contains(` 2a`),
				Contains(`-3a`),
				Contains(`+3b`),
				Contains(` 4a`),
				Contains(` 5a`),
				Contains(` 6a`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			}).
			Content(DoesNotContain("-3a").DoesNotContain("+3b")).
			SelectedLines(
				Contains("@@ -10,6 +10,6 @@"),
				Contains(" 10a"),
				Contains(" 11a"),
				Contains(" 12a"),
				Contains("-13a"),
				Contains("+13b"),
				Contains(" 14a"),
				Contains(" 15a"),
				Contains(`\ No newline at end of file`),
			)
	},
})
