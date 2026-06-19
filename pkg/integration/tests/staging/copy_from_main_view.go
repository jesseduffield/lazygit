package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a hunk from the focused main view to the clipboard, with the raw diff's +/- prefix stripped",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
		// Emulate the clipboard by writing to a file.
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		// Two consecutive added lines, so the selected hunk is a homogeneous addition.
		shell.UpdateFile("file1", "one\nADD1\nADD2\ntwo\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
				Contains("+ADD2"),
			).
			Press(keys.Universal.CopyToClipboard)

		// With no pager the main view shows the raw diff, so a homogeneous addition is
		// copied with the '+' column stripped — ready to paste into code.
		t.FileSystem().FileContent("clipboard", Equals("ADD1\nADD2\n"))
	},
})
