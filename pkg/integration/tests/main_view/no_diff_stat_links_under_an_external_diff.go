package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NoDiffStatLinksUnderAnExternalDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The files named in the diffstat are not linked under a diff renderer that says nothing about its rows",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		// An external diff whose output has nothing to say about which line of which
		// file each row shows. git writes the diffstat itself, so the names are there
		// to be clicked, but nothing could find the file they name.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "opaque", Type: "extDiff", Command: `sh -c 'echo EXT'`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("aaa.txt", "one\n")
		shell.CreateFileAndAdd("zzz.txt", "one\n")
		shell.Commit("one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("one"))

		// The click below is at a line of the diffstat, so the diff has to open with
		// the lines this expects.
		t.Views().Main().
			TopLines(
				Contains("commit"),
				Contains("Author:"),
				Contains("Date:"),
				Equals(""),
				Contains("one"),
				Equals("---"),
				Contains("aaa.txt"),
				Contains("zzz.txt"),
				Contains("2 files changed"),
			).
			Click(2, 7)

		// The name is no link, so the click is an ordinary one, which focuses the pane
		// it lands in. Were it a link, it would have been followed instead, and would
		// have had to report that it found no such file — the test fails on the toast
		// that leaves unacknowledged.
		t.Views().Main().IsFocused()
	},
})
