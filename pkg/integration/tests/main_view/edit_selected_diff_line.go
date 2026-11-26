package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditSelectedDiffLine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open the selected line of the main view's diff in the editor, at that line of the file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
		config.GetUserConfig().OS.Edit = "echo {{filename}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		// The addition is the third line of the file as it now stands.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+THREE")).
			Press(keys.Universal.Edit)

		// The editor is pointed at the file by absolute path.
		t.FileSystem().FileContent("edit-command", Contains("/repo/file1:3\n"))

		// A file header points at the file rather than at a line in it, so it opens the
		// file with no line to jump to.
		t.Views().Main().
			NavigateToLine(Contains("diff --git a/file1 b/file1")).
			Press(keys.Universal.Edit)

		t.FileSystem().FileContent("edit-command", Contains("/repo/file1\n"))
	},
})
