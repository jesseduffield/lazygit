package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickAFileInTheDiffStat = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jump to a file of a commit's diff by clicking the line that names it in the diffstat",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       30,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 600)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		// A long file at either end, so that the file jumped to is far below the
		// diffstat and has a diff under it to scroll past.
		shell.CreateFileAndAdd("aaa.txt", strings.Join(lines, "\n")+"\n")
		shell.CreateFileAndAdd("dir/bbb.txt", "one\n")
		shell.CreateFileAndAdd("zzz.txt", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("one"))

		// The click below is at a line of the diffstat, so the diff has to open with
		// the lines this expects.
		t.Views().Main().
			TopLines(
				Contains("commit"),
				Contains("Author:"),
				Contains("Date:"),
				Equals(""),
				Contains("one"),
				Equals("---"),
				Contains("aaa.txt"),
				Contains("dir/bbb.txt"),
				Contains("zzz.txt"),
				Contains("3 files changed"),
			).
			Click(2, 8)

		// The panel keeps the focus, and the diff goes to the file clicked.
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("one"))

		t.Views().Main().
			TopVisibleLine(Contains("diff --git a/zzz.txt b/zzz.txt"))

		t.Views().Commits().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// With the pane focused there is a selection to move, and the click moves it to
		// the file, exactly as picking the file from the menu would.
		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Universal.GotoTop)

		// The first file's diff begins on screen already, so the view stays where it
		// is; a jump only scrolls as far as it must once there is a selection to point
		// at the file with.
		t.Views().Main().
			TopVisibleLine(Contains("commit")).
			Click(2, 6).
			SelectedLines(
				Contains("diff --git a/aaa.txt b/aaa.txt"),
			).
			TopVisibleLine(Contains("commit")).
			// The last file's is far below, so that one is scrolled to.
			Click(2, 8).
			SelectedLines(
				Contains("diff --git a/zzz.txt b/zzz.txt"),
			).
			TopVisibleLine(Contains("diff --git a/zzz.txt b/zzz.txt"))
	},
})
