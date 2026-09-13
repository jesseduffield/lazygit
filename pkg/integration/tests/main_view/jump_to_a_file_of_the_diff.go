package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var JumpToAFileOfTheDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jump to a file of a commit's diff by picking it from a menu of the diff's files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A short terminal, so that the first file's diff is longer than the part of the
	// diff that has been read when the menu asks which files there are.
	Width:  100,
	Height: 20,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 600)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		shell.CreateFileAndAdd("aaa.txt", strings.Join(lines, "\n")+"\n")
		shell.CreateFileAndAdd("ccc.txt", "one\n")
		shell.CreateFileAndAdd("dir/bbb.txt", "one\n")
		// Another long one at the end, so that the file jumped to below has a diff
		// under it to scroll past and ends up at the top of the view.
		shell.CreateFileAndAdd("zzz.txt", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("ccc.txt", "two\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A menu offering the one file of a single-file diff would be a menu with
		// nothing to choose, so it says what it found instead.
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("two")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Main.JumpToFile)

		t.ExpectToast(Contains("There is only one file in this diff"))

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.Return)

		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			SelectedLine(Contains("one")).
			Press(keys.Universal.FocusMainView)

		// Every file of the diff is offered, in the order the diff shows them and by
		// the path the repo knows them by, including the ones below the part of the
		// diff that has been read.
		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Main.JumpToFile)

		t.ExpectPopup().Menu().
			Title(Equals("Jump to file")).
			Lines(
				Equals("aaa.txt"),
				Equals("ccc.txt"),
				Equals("dir/bbb.txt"),
				Equals("zzz.txt"),
				Equals("Cancel"),
			).
			Select(Equals("dir/bbb.txt")).
			Confirm()

		// The file lands where stepping to it with next-file would leave it: selected,
		// and at the top of the view.
		t.Views().Main().
			IsFocused().
			TopVisibleLine(Contains("diff --git a/dir/bbb.txt b/dir/bbb.txt")).
			SelectedLines(
				Contains("diff --git a/dir/bbb.txt b/dir/bbb.txt"),
			).
			Press(keys.Main.JumpToFile)

		// The menu filters as you type, which is the point of it for a diff of many
		// files.
		t.ExpectPopup().Menu().
			Title(Equals("Jump to file")).
			Filter("ccc").
			Lines(
				Equals("ccc.txt"),
			).
			Confirm()

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("diff --git a/ccc.txt b/ccc.txt"),
			)
	},
})
