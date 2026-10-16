package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PatchMarksFollowARendererSwitch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Switching diff renderers mid-build leaves the marks on the lines that are in the custom patch, wherever the new rendering puts them",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Two renderers that announce the metadata protocol — so that focusing the main
		// view keeps their output rather than falling back to git's own — and pass the
		// diff through under a banner of their own. The second one's banner is a line
		// longer, so every line of the diff it renders is a line further down than the
		// first one's.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "one", Command: `printf '\033]1717;1\007RENDERED BY ONE\n'; cat`},
			{Name: "two", Command: `printf '\033]1717;1\007RENDERED BY TWO\nAND ONE MORE LINE\n'; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(Contains("RENDERED BY ONE")).
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction().
			MarkedLines(
				Contains("-two"),
			).
			Press(keys.Universal.CycleDiffRenderers)

		t.ExpectToast(Equals("Diff renderer: two (2 of 2)"))

		// The marks are of lines of the diff, not of rows of the rendering, so the new
		// rendering has them on the same line of the file.
		t.Views().Main().
			Content(Contains("AND ONE MORE LINE")).
			MarkedLines(
				Contains("-two"),
			)
	},
})
