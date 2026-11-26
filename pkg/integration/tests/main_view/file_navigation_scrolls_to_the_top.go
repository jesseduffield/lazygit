package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FileNavigationScrollsToTheTop = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Going to a file of a diff brings it to the top of the view, unless it is on screen already",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       20,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// More context than fits on screen, so that a long file's change is further
		// down than a screenful from the header naming the file.
		config.GetUserConfig().Git.DiffContextSize = 30
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 100)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		long := strings.Join(lines, "\n") + "\n"
		longChanged := strings.Replace(long, "line100", "LINE100", 1)

		// Two long files with two short ones between them: the short ones are on screen
		// together, and there is enough diff below them to scroll to.
		shell.CreateFileAndAdd("aaa.txt", long)
		shell.CreateFileAndAdd("bbb.txt", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("ccc.txt", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("ddd.txt", long)
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("aaa.txt", longChanged)
		shell.UpdateFileAndAdd("bbb.txt", "one\nTWO\nthree\n")
		shell.UpdateFileAndAdd("ccc.txt", "one\nTWO\nthree\n")
		shell.UpdateFileAndAdd("ddd.txt", longChanged)
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		// A file below the viewport becomes the top of it, so that as much of the file
		// as possible is on screen.
		t.Views().Main().
			IsFocused().
			SelectionIsActive().
			Press(keys.Main.NextFile).
			TopVisibleLine(Contains("diff --git a/bbb.txt b/bbb.txt")).
			SelectedLines(
				Contains("diff --git a/bbb.txt b/bbb.txt"),
			).
			// The next file is on screen already, so the view stays where it is and
			// only the selection moves.
			Press(keys.Main.NextFile).
			TopVisibleLine(Contains("diff --git a/bbb.txt b/bbb.txt")).
			SelectedLines(
				Contains("diff --git a/ccc.txt b/ccc.txt"),
			).
			// Going back to a file above the viewport brings that one to the top.
			Press(keys.Main.PrevFile).
			Press(keys.Main.PrevFile).
			TopVisibleLine(Contains("diff --git a/aaa.txt b/aaa.txt")).
			SelectedLines(
				Contains("diff --git a/aaa.txt b/aaa.txt"),
			).
			// In hunk mode the selection is the file's first change rather than the row
			// the file begins at, and with this much context that change is further down
			// than a screenful. The selection has to be on screen, so the alignment gives
			// way and the selection is scrolled into view as any other jump's is.
			Press(keys.Main.ToggleSelectHunk).
			Press(keys.Main.NextFile).
			Press(keys.Main.NextFile).
			Press(keys.Main.NextFile).
			SelectedLines(
				Contains("-line100"),
				Contains("+LINE100"),
			).
			SelectedLineIsVisible()
	},
})
