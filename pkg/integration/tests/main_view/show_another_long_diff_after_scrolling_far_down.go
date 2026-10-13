package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowAnotherLongDiffAfterScrollingFarDown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Selecting another commit shows its diff when the main view is scrolled far down a long diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       20,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Longer than the 5000 lines that a render reads at most before it waits for
		// the user to scroll further.
		lines := make([]string, 6000)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%04d", i+1)
		}
		long := strings.Join(lines, "\n") + "\n"

		shell.CreateFileAndAdd("first.txt", long)
		shell.Commit("first commit")
		shell.CreateFileAndAdd("second.txt", long)
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("second commit")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.GotoBottom).
			SelectedLines(Contains("+line6000")).
			Press(keys.Universal.Return)

		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			SelectedLine(Contains("first commit"))

		t.Views().Main().
			/* EXPECTED:
			Content(Contains("diff --git a/first.txt b/first.txt"))
			ACTUAL: */
			Content(Contains("diff --git a/second.txt b/second.txt"))
	},
})
