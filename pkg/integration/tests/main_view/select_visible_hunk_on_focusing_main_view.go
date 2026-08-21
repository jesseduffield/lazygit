package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectVisibleHunkOnFocusingMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view in hunk mode picks a block that begins on screen, leaving the diff where it is",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
		// One line per scroll, so that the test can put the top of the view exactly
		// where it wants it, and enough context to scroll about within one hunk.
		cfg.GetUserConfig().Gui.ScrollHeight = 1
		cfg.GetUserConfig().Git.DiffContextSize = 20
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 80)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		for _, i := range []int{10, 20, 70} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		scrollDown := func(lines int) {
			for range lines {
				t.Views().Files().Press(keys.Universal.ScrollDownMain)
			}
		}

		t.Views().Files().IsFocused()

		// The top of the view is the second line of the first change block, so that
		// block is only half on screen; the second one begins below it, in full.
		scrollDown(15)

		t.Views().Main().
			OriginY(15).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.FocusMainView)
			}).
			IsFocused().
			SelectedLines(
				Contains("-line20"),
				Contains("+LINE20"),
			).
			OriginY(15).
			PressEscape()

		// Now nothing begins on screen: the second block starts just above the top and
		// the third change is far below. The half-visible block is what there is, so it
		// is selected — with its first line off screen, since the view stays put.
		scrollDown(11)

		t.Views().Main().
			OriginY(26).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.FocusMainView)
			}).
			IsFocused().
			SelectedLines(
				Contains("-line20"),
				Contains("+LINE20"),
			).
			SelectedLineIdx(25).
			OriginY(26).
			PressEscape()

		// And a click inside a block that begins above the viewport selects the whole
		// block without pulling the view up to its start either.
		t.Views().Main().
			Click(0, 0).
			IsFocused().
			SelectedLines(
				Contains("-line20"),
				Contains("+LINE20"),
			).
			OriginY(26)
	},
})
