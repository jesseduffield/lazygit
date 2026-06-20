package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageFromMainViewWithConformingPager = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focus the main view under a pager that announces the diff-metadata protocol (a handshake); its output is trusted, so the diff is not re-rendered raw and staging works on it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		// A fake conforming pager: it emits the version-only OSC 1717 handshake as its
		// first output (announcing it speaks the protocol), prints a marker line so the
		// test can tell its output apart from the raw diff, then passes the diff through
		// unchanged. The probe finds the handshake and trusts the pager, so focusing
		// does not fall back to the raw diff. (The passed-through diff is structurally
		// intact, so the selection still resolves via the buffer parser.)
		cfg.GetUserConfig().Git.Pagers = []config.PagingConfig{
			{Pager: "printf '\\033]1717;1\\007CONFORMING-PAGER\\n'; cat"},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			)

		// The handshake is swallowed (it leaves no visible bytes), but the pager's
		// marker line proves its output is what's shown.
		t.Views().Main().Content(Contains("CONFORMING-PAGER"))

		t.Views().Files().Press(keys.Universal.FocusMainView)

		// Focusing did not re-render raw — the marker is still there — and the selection
		// resolved on the pager's (structure-preserving) output.
		t.Views().Main().
			IsFocused().
			Content(Contains("CONFORMING-PAGER")).
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
			})
	},
})
