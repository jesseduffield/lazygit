package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditHunkInFocusedDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit the hunk around the selection in an editor, and stage what comes back",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
		// Stand in for the editor: record the line it was pointed at, outside the
		// repo so that the files panel keeps saying what the test is about, then
		// write a patch that stages something neither side of the diff says. That
		// is the point of editing a hunk.
		cfg.GetUserConfig().OS.EditAtLineAndWait = "echo {{line}} > ../edit-line && " +
			"printf '%s\\n' '--- a/file1' '+++ b/file1' '@@ -1,3 +1,3 @@' " +
			"' one' '-two' '+TWO_EDITED' ' three' > {{filename}}"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\nTWO\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-two")).
			Press(keys.Main.EditSelectHunk)

		// The patch is written with a two-line header, so the deletion the cursor was
		// on is its fifth line.
		t.FileSystem().FileContent("../edit-line", Equals("5\n"))

		// What the editor wrote went into the index, leaving the working tree as it
		// was: the file is changed on both sides now, differently.
		t.Views().Files().Lines(
			Contains("MM").Contains("file1"),
		)
		t.Views().Secondary().ContainsLines(
			Contains("-two"),
			Contains("+TWO_EDITED"),
		)
		t.Views().Main().ContainsLines(
			Contains("-TWO_EDITED"),
			Contains("+TWO"),
		)
	},
})
