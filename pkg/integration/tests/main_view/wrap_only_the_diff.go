package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var longLine = "start " + strings.Repeat("word ", 40) + "end"

var WrapOnlyTheDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Turning wrapping off for diffs leaves everything else the main view shows wrapped",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.WrapLinesInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\n")
		shell.Commit(longLine)

		shell.UpdateFile("file1", "one\n"+longLine+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// The option is about diffs, so the added line runs off the edge of the pane in
		// one piece.
		t.Views().Files().
			IsFocused().
			Tap(func() {
				t.Views().Main().ContainsViewLines(Contains("start").Contains("end"))
			})

		// A branch's commit log is no diff, and the option leaves it alone: the same
		// text is laid out over as many lines as it takes.
		t.Views().Branches().
			Focus()

		t.Views().Main().
			ContainsViewLines(Contains("end").DoesNotContain("start"))
	},
})
