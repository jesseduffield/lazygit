package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionInBothPanesWhenSwitchingDiffRenderers = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching to another diff renderer keeps the place in the lower pane too, not only in the upper one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Renderers that speak the metadata protocol, so that focusing the main view
		// keeps their rendering rather than falling back to git's own diff.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "plain", Command: `printf '\033]1717;1\007'; cat`},
			// The same diff, three lines further down the view. (Lines before the
			// diff's own header aren't part of it, so it still reads the same.)
			{Name: "banner", Command: `printf '\033]1717;1\007'; printf 'rendered for you\n\n\n'; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		// Four staged changes to have a diff worth scrolling in the lower pane, and one
		// unstaged one to split the file's diff across both panes.
		for _, i := range []int{5, 15, 25, 35} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
		shell.GitAddAll()

		lines[39] = strings.ToUpper(lines[39])
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			PressTab()

		t.Views().Secondary().
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
			// The banner pushed the whole diff three lines down, and the lower pane came
			// along with it, just as the upper one would have.
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(38).
			OriginY(17)
	},
})
