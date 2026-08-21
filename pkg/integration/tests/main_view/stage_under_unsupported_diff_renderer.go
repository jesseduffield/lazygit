package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageUnderUnsupportedDiffRenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff renderer that says nothing about what it renders is replaced by git's own diff when the main view is focused, so the diff can still be staged from",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// `cat -n` numbers every line, which pushes the +/- column off the start of it:
		// the diff can't be read back from the text, and cat says nothing about what it
		// is rendering, so there is no way to act on what it produces.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: "cat -n"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// While browsing, the renderer's output is what is shown.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			)
		t.Views().Main().Content(Contains("1  diff --git a/file1 b/file1"))

		// Focusing it to act on it brings git's own diff instead, and the selection goes
		// on that.
		t.Views().Files().Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(DoesNotContain("1  diff --git")).
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction()

		t.Views().Secondary().ContainsLines(
			Contains("-three"),
			Contains("+THREE"),
		)
		// The re-render after staging stays with git's diff, so the next hunk can be
		// staged as well.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-nine"),
				Contains("+NINE"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			Contains("M  file1"),
		)
	},
})
