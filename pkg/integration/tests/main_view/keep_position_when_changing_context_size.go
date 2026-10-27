package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionWhenChangingContextSize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Changing the diff's context size keeps the line you were looking at where it was",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		// Four changes, far enough apart that they stay four hunks as the context
		// size grows.
		for _, i := range []int{5, 15, 25, 35} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Main.NextHunk).
			Press(keys.Main.NextHunk).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-line35"),
			).
			// The diff is longer than the view, so getting to the last hunk scrolled
			// it: the selected line sits 21 rows down the screen.
			SelectedLineIdx(35).
			OriginY(14).
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			// A context line more on either side of each of the four hunks pushes the
			// selected line seven lines further into the diff. The view follows it, so
			// it is still the same line on the same screen row (42 - 21 = 21).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(42).
			OriginY(21).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			}).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			// And the same the other way (28 - 21 = 7).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(28).
			OriginY(7).
			// Leaving the view gives up the selection but not the scroll position, and
			// with no selection to keep, it is the middle visible line that stays put.
			PressEscape()

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			})

		// The middle visible line here is a hunk's header, and a context-size change
		// rewrites those — they name the lines the hunk covers. So the restore falls
		// back to the nearest line that does survive, the context line just below it,
		// and puts that back on the row it was on.
		t.Views().Main().
			SelectionIsHidden().
			OriginY(12)
	},
})
