package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildFromMainViewWithUnsupportedPager = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a custom patch from a commit's focused main view under a pager that restructures the diff without emitting metadata; it falls back to the raw diff so the selection is still toggleable",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		// `cat -n` numbers every line, which the buffer parser can't resolve, and cat
		// emits no metadata, so the focused main view must fall back to the raw diff.
		cfg.GetUserConfig().Git.Pagers = []config.PagingConfig{
			{Pager: "cat -n"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\n")
		shell.Commit("update")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			)

		// While browsing, the main view shows the pager's output (numbered lines).
		t.Views().Main().Content(Contains("  +THREE"))

		// Focusing the main view falls back to the raw diff, so the change line is
		// resolved and the selection is toggleable into a custom patch.
		t.Views().Commits().Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Secondary().
			ContainsLines(
				Contains("-three"),
				Contains("+THREE"),
			)
	},
})
