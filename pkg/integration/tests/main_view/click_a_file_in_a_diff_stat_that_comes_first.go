package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickAFileInADiffStatThatComesFirst = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jump to a file by clicking its name in a diffstat the renderer's handshake runs into",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		// A renderer that announces the protocol and then passes the diff on as it came.
		// The handshake has no newline after it, so it runs into the first line the
		// renderer is given — which for the diff of a range of commits is the first
		// entry of the diffstat, there being no commit above it.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007'; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("aaa.txt", "one\n")
		shell.Commit("one")
		shell.CreateFileAndAdd("zzz.txt", "one\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("two")).
			Press(keys.Universal.ToggleRangeSelect).
			SelectNextItem().
			SelectedLines(
				Contains("two"),
				Contains("one"),
			)

		// Below the line lazygit writes to say what the diff is of, the diff opens with
		// the diffstat, and the handshake runs into its first entry. That is the one
		// clicked here.
		t.Views().Main().
			TopLines(
				Contains("Showing diff for range"),
				Equals(""),
				Contains("aaa.txt"),
				Contains("zzz.txt"),
				Contains("2 files changed"),
			).
			Click(2, 2)

		t.Views().Commits().IsFocused()

		t.Views().Main().
			TopVisibleLine(Contains("diff --git a/aaa.txt b/aaa.txt"))
	},
})
