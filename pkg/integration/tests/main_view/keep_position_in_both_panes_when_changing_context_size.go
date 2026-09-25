package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionInBothPanesWhenChangingContextSize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Changing the diff's context size keeps the place in the lower pane too, not only in the upper one",
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

		// Four staged changes, far enough apart that they stay four hunks as the
		// context size grows, and one unstaged one to split the file's diff across
		// both panes.
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

		// The lower pane holds the staged changes; getting to the last of them scrolls
		// it, so there is a position to lose.
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
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			SelectedLines(
				Contains("-line35"),
			).
			SelectedLineIdx(42).
			OriginY(21)
	},
})
