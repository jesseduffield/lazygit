package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepBothHalvesOfAChangeSelected = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A change selected on the one row a renderer draws it as is selected on both rows of a renderer that splits it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			// Git's own diff, which has a row for each half of a change. It announces
			// the metadata protocol so that its output is taken at its word; the rows
			// are placed by parsing it as the diff it still looks like.
			{Name: "unified", Command: `printf '\033]1717;1\007'; cat`},
			// A renderer that puts the two halves of a change beside each other on one
			// row, which nothing but the records it states could place. It ignores the
			// diff it is given and draws this one.
			{Name: "columns", Command: `printf '\033]1717;1\007'; ` +
				`printf '\033]1717;1;f;;;file1\007file1\n'; ` +
				`printf '\033]1717;1;h;1;;file1\007@@\n'; ` +
				`printf '\033]1717;1;c;1;;file1\007one    one\n'; ` +
				`printf '\033]1717;1;d;2;2;file1\007two    \033]1717;1;a;2;;file1\007TWO\n'; ` +
				`printf '\033]1717;1;c;3;;file1\007three  three\n'; ` +
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
				Contains("+TWO"),
			).
			// The change is one row here, and selecting it selects that row: both
			// halves are on it.
			Press(keys.Universal.CycleDiffRenderers).
			Tap(func() {
				t.ExpectToast(Equals("Diff renderer: columns (2 of 2)"))
			}).
			SelectedLines(
				Contains("two    TWO"),
			).
			// Split apart again, the same change is the same two lines it was.
			Press(keys.Universal.CycleDiffRenderers).
			Tap(func() {
				t.ExpectToast(Equals("Diff renderer: unified (1 of 2)"))
			}).
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			)
	},
})
