package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditClickedDiffLine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Alt- or shift-click a line of the main view's diff to open it in the editor, without focusing the view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused()

		// The click points at the line itself, so the main view can stay unfocused
		// and unselected. You read a diff where it is and click into it.
		t.Views().Main().
			AltClick(0, 8).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/file1:3\n"))
			}).
			SelectionIsHidden()

		t.Views().Files().
			IsFocused()

		// Shift-click is bound to the same thing, since neither modifier reaches
		// lazygit in every terminal. A context line names a line of the file like any
		// other row.
		t.Views().Main().
			ShiftClick(0, 9).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/file1:4\n"))
			})

		// A popup taking the focus swallows clicks on the views behind it. This one
		// stays live, so a diff can still be read and clicked into while a popup is
		// up.
		t.Views().Files().
			Press(keys.Universal.Remove)

		t.Views().Menu().
			IsFocused()

		t.Views().Main().
			AltClick(0, 6).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/repo/file1:2\n"))
			})

		t.Views().Menu().
			IsFocused().
			Press(keys.Universal.Return)
	},
})
