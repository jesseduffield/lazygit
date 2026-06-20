package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageFromMainViewWithUnsupportedPager = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focus the main view under a pager that restructures the diff without emitting metadata; it falls back to the raw diff so the selection is still stageable",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		// `cat -n` prepends a line number to every line, which pushes the diff's +/-
		// markers off the start of the line: the buffer parser can't resolve it, and
		// cat emits no metadata, so the focused main view must fall back to the raw diff.
		cfg.GetUserConfig().Git.Pagers = []config.PagingConfig{
			{Pager: "cat -n"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		// Two separate change blocks, so staging one leaves the file split (the other
		// stays unstaged) and the staged side renders into the secondary view.
		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// While browsing, the main view shows the pager's output: `cat -n` has numbered
		// every line (the first being the `diff --git` header), so the diff isn't in
		// its raw form.
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			)
		t.Views().Main().Content(Contains("1  diff --git a/file1 b/file1"))

		// Focusing the main view falls back to the raw diff, so the change line is
		// resolved and selected (the pager's line numbers are gone).
		t.Views().Files().Press(keys.Universal.FocusMainView)

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
			PressPrimaryAction().
			Tap(func() {
				t.Views().Secondary().
					ContainsLines(
						Contains("-three"),
						Contains("+THREE"),
					)
			}).
			// The other block stays unstaged, also shown raw.
			ContainsLines(
				Contains("-nine"),
				Contains("+NINE"),
			)
	},
})
