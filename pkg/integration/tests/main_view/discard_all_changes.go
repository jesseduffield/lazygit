package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardAllChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard all changes of a file from the focused main view, then land on the next file's diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.CreateFileAndAdd("file2", "1\n2\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
		shell.UpdateFile("file2", "1\n2\n3\n4\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("   M file1"),
				Equals("   M file2"),
			).
			SelectNextItem().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+three")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			}).
			SelectedLines(Contains("+four")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			})

		t.Views().Files().Lines(
			Equals(" M file2"),
		)
		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+3"))
	},
})
