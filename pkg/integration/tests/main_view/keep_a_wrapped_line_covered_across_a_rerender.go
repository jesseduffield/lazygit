package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepAWrappedLineCoveredAcrossARerender = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A selection over a line too long for the view still covers all of it after a re-render",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        80,
	Height:       20,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		long := strings.Repeat("word ", 40)
		lines := make([]string, 20)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		before := strings.Join(lines[:10], "\n") + "\n"
		after := strings.Join(lines[10:], "\n") + "\n"

		shell.CreateFileAndAdd("file1", before+long+"\n"+after)
		shell.Commit("one")

		shell.UpdateFile("file1", before+"CHANGED "+long+"\n"+after)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The changed line is far too long for the view, so each half of the change
		// is drawn as several view lines, and hunk mode selects all of them.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-word word"),
				Contains("+CHANGED word"),
			).
			SelectedViewLineRange(8, 16).
			// The same two lines of the diff, wrapped the same way, are still covered
			// to their ends once the diff has been rendered again.
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 4"))
			}).
			SelectedLines(
				Contains("-word word"),
				Contains("+CHANGED word"),
			).
			SelectedViewLineRange(9, 17)
	},
})
