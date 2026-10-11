package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepPositionWhenIgnoringWhitespaceRemovesIt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignoring whitespace where that takes the line you were on out of the diff lands on the nearest line it kept",
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
		shell.CreateFileAndAdd("file2", "one\ntwo\nthree\n")
		shell.Commit("one")

		// Real changes at lines 5, 15 and 25, and a whitespace-only one at 35, far
		// enough apart to be hunks of their own.
		for _, i := range []int{5, 15, 25} {
			lines[i-1] = strings.ToUpper(lines[i-1])
		}
		lines[34] = " " + lines[34]
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")

		// Nothing but reindentation, so ignoring whitespace leaves no diff at all.
		shell.UpdateFile("file2", "  one\n  two\n  three\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			SelectNextItem().
			SelectedLine(Contains("file1")).
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
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			// That hunk was a whitespace change and nothing else, so ignoring
			// whitespace takes it — and the context around it — out of the diff
			// entirely. The nearest line the diff kept is the last line of the hunk
			// above, so that is where the selection lands; it goes back on the row it
			// was on itself, which leaves everything above it exactly where it was.
			SelectedLines(
				Contains(" line28"),
			).
			SelectedLineIdx(30).
			OriginY(14).
			// The whole diff can go this way, and then there is nothing to land on.
			Press(keys.Universal.ToggleWhitespaceInDiffView).
			PressEscape()

		t.Views().Files().
			IsFocused().
			SelectNextItem().
			SelectedLine(Contains("file2")).
			Press(keys.Universal.ToggleWhitespaceInDiffView)

		t.Views().Main().
			Content(Equals(""))
	},
})
