package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectedRangeWhenChangingContextSize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A range selection still covers the same lines of the diff after the context size changes",
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

		for _, i := range []int{5, 15, 25, 35} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// A range from a change down into the context below it, so that the cursor is
		// on the last line of the selection and the other end is three lines above.
		t.Views().Main().
			IsFocused().
			Press(keys.Main.NextHunk).
			Press(keys.Main.NextHunk).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("-line25"),
				Contains("+LINE25"),
				Contains(" line26"),
				Contains(" line27"),
			).
			// Both ends are still lines of the diff with more context around the
			// change, so the selection still covers the same four.
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			SelectedLines(
				Contains("-line25"),
				Contains("+LINE25"),
				Contains(" line26"),
				Contains(" line27"),
			).
			// With a single line of context, the line the cursor was on is no longer in
			// the diff. The end that survived stays put and the cursor lands on the
			// nearest line that is left, so the selection shrinks with the diff.
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			}).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 1"))
			}).
			SelectedLines(
				Contains("-line25"),
				Contains("+LINE25"),
				Contains(" line26"),
			).
			// The other way round: a range extended upwards, so that it is the far end
			// that the shrinking context takes away. There is no guessing which line
			// inherits it, so what is left is the line the cursor is on.
			PressEscape().
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 3"))
			}).
			SelectedLines(
				Contains(" line26"),
			).
			Press(keys.Universal.NextItem).
			Press(keys.Universal.RangeSelectUp).
			Press(keys.Universal.RangeSelectUp).
			Press(keys.Universal.RangeSelectUp).
			SelectedLines(
				Contains("-line25"),
				Contains("+LINE25"),
				Contains(" line26"),
				Contains(" line27"),
			).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 1"))
			}).
			SelectedLines(
				Contains("-line25"),
			)
	},
})
