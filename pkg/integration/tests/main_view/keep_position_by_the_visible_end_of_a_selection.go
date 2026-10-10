package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionByTheVisibleEndOfASelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A re-render keeps the place by the end of a selected hunk that is on screen when its other end isn't",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
		// One line per scroll, so that the test can put the top of the view exactly
		// where it wants it.
		cfg.GetUserConfig().Gui.ScrollHeight = 1
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 60)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		// A first change tall enough to be scrolled halfway out of the view, and more
		// of them below it, so that a context-size change moves the lines further down
		// the diff by more than it moves the first change.
		for _, i := range []int{10, 11, 12, 13, 14, 15, 30, 45} {
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
				Contains("-line10"),
				Contains("-line11"),
				Contains("-line12"),
				Contains("-line13"),
				Contains("-line14"),
				Contains("-line15"),
				Contains("+LINE10"),
				Contains("+LINE11"),
				Contains("+LINE12"),
				Contains("+LINE13"),
				Contains("+LINE14"),
				Contains("+LINE15"),
			).
			SelectedLineIdx(8).
			// Scroll past the start of the selected block, leaving its last lines on
			// screen and the cursor above the top of the view.
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			OriginY(14).
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			// The block's last line was the fifth row of the screen, and one context
			// line more above the block puts it a line further down the diff: the view
			// follows it, rather than the middle visible line, which the hunks below
			// have pushed further still.
			OriginY(15).
			SelectedLines(
				Contains("-line10"),
				Contains("-line11"),
				Contains("-line12"),
				Contains("-line13"),
				Contains("-line14"),
				Contains("-line15"),
				Contains("+LINE10"),
				Contains("+LINE11"),
				Contains("+LINE12"),
				Contains("+LINE13"),
				Contains("+LINE14"),
				Contains("+LINE15"),
			)
	},
})
