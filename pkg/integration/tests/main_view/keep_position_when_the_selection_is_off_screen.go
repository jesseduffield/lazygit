package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionWhenTheSelectionIsOffScreen = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A re-render keeps the lines that are on screen where they are, not a selection scrolled away from",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Half a diff per scroll, to leave the selection well behind in two presses.
		cfg.GetUserConfig().Gui.ScrollHeight = 15
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 60)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		for _, i := range []int{5, 15, 25, 35, 45, 55} {
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
			SelectedLines(
				Contains("-line05"),
			).
			// Reading on past the selection leaves it far behind, off the top of the
			// view.
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			OriginY(30).
			Press(keys.Universal.DecreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			// A context line less on either side of the four hunks above what is on
			// screen pulls it nine lines up the diff, and the view follows it there: the
			// lines the user was reading are still on the rows they were on.
			OriginY(21).
			// The selection is where it always was, on its own line of the diff, rather
			// than having been dragged back into view.
			SelectedLines(
				Contains("-line05"),
			).
			SelectedLineIdx(7)
	},
})
