package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionWhenSwitchingDiffRenderers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching to another diff renderer keeps the line you were looking at where it was",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Both announce the metadata protocol, so that their output is taken at its
		// word and shown as it is; a renderer that says nothing about what it renders
		// is replaced by git's own diff as soon as the main view is focused.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "plain", Command: `printf '\033]1717;1\007'; cat`},
			// The same diff, three lines further down the view. (Lines before the
			// diff's own header aren't part of it, so it still reads the same.)
			{Name: "banner", Command: `printf '\033]1717;1\007rendered for you\n\n\n'; cat`},
		}
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

		t.Views().Main().
			IsFocused().
			Press(keys.Main.NextHunk).
			Press(keys.Main.NextHunk).
			Press(keys.Main.NextHunk).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(35).
			OriginY(14).
			Press(keys.Universal.CycleDiffRenderers).
			Tap(func() {
				t.ExpectToast(Equals("Diff renderer: banner (2 of 2)"))
			}).
			// The banner pushed the whole diff three lines down, and the view came
			// along with it: the same line on the same screen row (38 - 17 = 21).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(38).
			OriginY(17).
			Press(keys.Universal.CycleDiffRenderers).
			Tap(func() {
				t.ExpectToast(Equals("Diff renderer: plain (1 of 2)"))
			}).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(35).
			OriginY(14)
	},
})
