package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardFromACommitOnlyWhereItCanBeRewritten = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding lines is refused over a diff that belongs to no commit we may rewrite, and over the custom patch itself",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFile("file1", "one\nTWO\nthree\n")
		shell.Stash("a stashed change")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A stash entry is no commit of ours to rewrite, so its lines can go into a
		// patch but can't be taken out of what they are part of.
		t.Views().Stash().
			Focus().
			Lines(
				Contains("a stashed change").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Changes can only be discarded from local commits")).
			Confirm()

		// The pane previewing the patch shows the patch's own lines, which are not the
		// commit's to discard; space takes them back out of the patch instead.
		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-two"),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Lines shown here are the custom patch's")).
			Confirm()
	},
})
