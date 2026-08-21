package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit what you staged without leaving the focused main view, but not while looking at a commit's diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nADDED\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Committing acts on the working tree, so it is offered over the working tree's
		// diff and does nothing over a commit's.
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Files.CommitChanges).
			// Nothing happened: the commit message panel would have taken the focus.
			IsFocused()

		t.Views().Files().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADDED"),
			).
			PressPrimaryAction()

		// Staging it left nothing unstaged, so the diff — and the focus with it — is in
		// the pane the staged side has.
		t.Views().Secondary().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("staged from the diff").
			Confirm()

		// The commit left the staged side with nothing in it, so the pane it was made
		// from is gone and the focus is in the one that is still there.
		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsFocused().
			Content(Contains("No changed files"))

		t.Views().Commits().Lines(
			Contains("staged from the diff"),
			Contains("one"),
		)
	},
})
