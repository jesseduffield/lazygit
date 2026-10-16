package main_view

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenderTheDiffBesideThePatchMarks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff that the custom patch's marks are shown over is rendered to the width they leave it, so that a line as wide as the view still fits beside them, whether or not the view has the focus",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file", "one\nTWO\nthree\n")
		// git draws the diffstat line of a file with this many changes as wide as it is
		// told the view is.
		shell.CreateFileAndAdd("many", strings.Repeat("line\n", 300))
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			MarkedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			ContainsViewLines(
				Contains("many | 300"),
				Contains("2 files changed"),
			)

		// The commits panel renders the diff again as it takes the focus back, and the
		// marks are still shown over it.
		t.Views().Main().Press(keys.Universal.Return)

		t.Views().Commits().IsFocused()
		t.Views().Main().
			MarkedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			ContainsViewLines(
				Contains("many | 300"),
				Contains("2 files changed"),
			)
	},
})
