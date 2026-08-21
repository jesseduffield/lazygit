package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionInBothPanesWhenIgnoringWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignoring whitespace keeps the place in the lower pane too, not only in the upper one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 60)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		// Staged: real changes at lines 5, 25 and 45, and a whitespace-only one at 15,
		// whose hunk goes when whitespace stops counting.
		lines[4] = strings.ToUpper(lines[4])
		lines[14] = " " + lines[14]
		lines[24] = strings.ToUpper(lines[24])
		lines[44] = strings.ToUpper(lines[44])
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
		shell.GitAddAll()

		// And one unstaged change, to split the file's diff across both panes.
		lines[59] = strings.ToUpper(lines[59])
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
				Contains("-line45"),
			).
			SelectedLineIdx(35).
			OriginY(14).
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			// The hunk above this one held nothing but a whitespace change, so it is
			// gone and has taken nine lines of the lower pane's diff with it — leaving
			// the line we were on where it was on the screen.
			SelectedLines(
				Contains("-line45"),
			).
			SelectedLineIdx(26).
			OriginY(5)
	},
})
