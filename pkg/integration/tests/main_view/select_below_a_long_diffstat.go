package main_view

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectBelowALongDiffstat = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view over a commit whose diff begins below a long diffstat still shows a selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A short terminal, so that the diffstat below fills more than the screenful the
	// first paint reveals, and the diff itself is longer than the initial read of it.
	Width:       100,
	Height:      20,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		for i := range 40 {
			shell.CreateFileAndAdd(fmt.Sprintf("file%02d", i+1), "one\ntwo\nthree\n")
		}
		shell.Commit("first commit")
		for i := range 40 {
			shell.UpdateFileAndAdd(fmt.Sprintf("file%02d", i+1), "one\nTWO\nthree\n")
		}
		shell.Commit("touch every file")
		shell.EmptyCommit("nothing to see here")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A commit with nothing to select leaves the pane showing no selection.
		t.Views().Commits().
			Focus().
			Lines(
				Contains("nothing to see here").IsSelected(),
				Contains("touch every file"),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			PressEscape()

		// The commit below it has plenty to select, even though none of it is among
		// the diffstat the first paint shows.
		t.Views().Commits().
			IsFocused().
			SelectNextItem()

		t.Views().Main().Content(Contains("40 files changed"))

		t.Views().Commits().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive()
	},
})
