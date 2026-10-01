package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardDiffLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard a hunk of the working tree's diff from the focused main view, and unstage one from the staged half",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFileAndAdd("file1", "one\nSTAGED\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.UpdateFile("file1", "one\nSTAGED\ntwo\nthree\nfour\nfive\nsix\nseven\nUNSTAGED\neight\nnine\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("MM file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// Discarding from the unstaged side throws the change away, so it asks first.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+UNSTAGED"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Discard change")).
					Content(Contains("Are you sure you want to discard this change")).
					Confirm()
			})

		// Nothing is unstaged any more, so that pane is gone and the focus has followed
		// the file's remaining changes into the staged one.
		t.Views().Files().Lines(
			Contains("M  file1"),
		)
		t.Views().Main().IsInvisible()

		// There the same key means "I don't want this staged", which is unstaging, so it
		// doesn't ask.
		t.Views().Secondary().
			IsFocused().
			Title(Equals("Staged changes")).
			SelectedLines(
				Contains("+STAGED"),
			).
			Press(keys.Universal.Remove)

		t.Views().Files().Lines(
			Contains(" M file1"),
		)
		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsFocused().
			Title(Equals("Unstaged changes")).
			Content(Contains("+STAGED"))
	},
})
