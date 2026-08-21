package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectBelowALongCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view over a commit whose diff begins below a very long commit message still shows a selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A short terminal, so that the message below is longer than a render is asked to
	// read, and the diff under it is only reached by reading on.
	Width:       100,
	Height:      20,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		body := make([]string, 1000)
		for i := range body {
			body[i] = fmt.Sprintf("message line %d", i+1)
		}
		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("a commit with a great deal to say\n\n" + strings.Join(body, "\n"))

		shell.EmptyCommit("nothing to see here")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A commit with nothing to select leaves the pane showing no selection.
		t.Views().Commits().
			Focus().
			Lines(
				Contains("nothing to see here").IsSelected(),
				Contains("a commit with a great deal to say"),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsHidden().
			PressEscape()

		// The commit below it has a change, a thousand lines further down than a render
		// reads by itself. The pane reads on until it knows, rather than taking the
		// answer from the commit before it or waiting for the user to scroll.
		t.Views().Commits().
			IsFocused().
			SelectNextItem().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			// The change is in the pane, a thousand lines below where the render stopped
			// reading of its own accord.
			Content(Contains("+TWO")).
			SelectionIsActive()
	},
})
