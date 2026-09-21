package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionWhenIgnoringWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignoring whitespace keeps the line you were looking at where it was, even when it turns into a context line",
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

		// Real changes at lines 5 and 25, whitespace-only ones at 15, 27 and 35. The
		// one at 27 shares a hunk with the change at 25, so ignoring whitespace turns
		// it into a context line rather than taking its hunk away.
		lines[4] = strings.ToUpper(lines[4])
		lines[14] = " " + lines[14]
		lines[24] = strings.ToUpper(lines[24])
		lines[26] = lines[26] + " "
		lines[34] = " " + lines[34]
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
			SelectedLines(
				Contains("-line25"),
			).
			SelectedLineIdx(26).
			OriginY(14).
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			// The hunk above this one held nothing but a whitespace change, so it is
			// gone and has taken nine lines of diff with it. This is still the line we
			// were on, on the row we were on (26 - 14 = 17 - 5).
			SelectedLines(
				Contains("-line25"),
			).
			SelectedLineIdx(17).
			OriginY(5).
			// And back again, whitespace and all.
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			SelectedLines(
				Contains("-line25"),
			).
			SelectedLineIdx(26).
			OriginY(14).
			// The whitespace-only change further down this hunk is a line of the file
			// like any other: ignoring whitespace shows it as context instead of as a
			// change, and that is still where we are.
			Press(keys.Main.NextHunk).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains("+line27"),
			).
			SelectedLineIdx(30).
			OriginY(14).
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			SelectedLines(
				Contains(" line27"),
			).
			SelectedLineIdx(20).
			OriginY(4)
	},
})
