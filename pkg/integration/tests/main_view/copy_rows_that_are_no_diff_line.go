package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyRowsThatAreNoDiffLine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copying rows that stand for no line of the diff says so rather than copying nothing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		cfg.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
		// A renderer that states which line of the file each row of its diff shows, and
		// ends with a row of its own that shows none. It ignores its input and prints
		// this one.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007'; ` +
				`printf '\033]1717;1;f;;;file1\007file1\n'; ` +
				`printf '\033]1717;1;c;1;;file1\007 one\n'; ` +
				`printf '\033]1717;1;d;2;2;file1\007-two\n'; ` +
				`printf '\033]1717;1;a;2;;file1\007+TWO\n'; ` +
				`printf '\033]1717;1;c;3;;file1\007 three\n'; ` +
				`printf -- '--- that was the diff ---\n'; ` +
				`cat >/dev/null`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nTWO\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			SelectNextItem().
			SelectNextItem().
			SelectNextItem().
			SelectedLines(
				Contains("that was the diff"),
			).
			Press(keys.Universal.CopyToClipboard)

		t.ExpectToast(Equals("Nothing in the selection could be found in the diff"))
		t.FileSystem().PathNotPresent("clipboard")
	},
})
