package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// The clipboard is emulated by a file, so that this works on CI too.
func expectClipboard(t *TestDriver, matcher *TextMatcher) {
	defer t.Shell().DeleteFile("clipboard")

	t.FileSystem().FileContent("clipboard", matcher)
}

var CopySelectedDiffLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the selected diff lines from the focused main view, as the diff reads rather than as the renderer drew it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = true
		// Emulate the clipboard by writing to a file.
		cfg.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
		// A renderer that decorates every line of a diff's body, so that what is on
		// screen is not what the diff says. It announces the metadata protocol, so that
		// its output is taken at its word rather than replaced by git's own; and it
		// reads the +/- column, so it wants its input uncoloured.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{
				Command: `printf '\033]1717;1\007'; ` +
					`awk '/^@@/ { body = 1 } body && /^[-+ ]/ { print $0 " <<<"; next } { print }'`,
				ColorArg: "never",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("two")

		shell.UpdateFile("file1", "one\nTWO\nADD1\nADD2\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1 <<<"),
				Contains("+ADD2 <<<"),
			).
			Press(keys.Universal.CopyToClipboard)

		// The renderer's decoration is nowhere in what was copied, and a selection that
		// is all additions loses its '+' column, ready to be pasted into code.
		expectClipboard(t, Equals("ADD1\nADD2\n"))

		// A commit's diff is copied the same way, through the panel that produced it.
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("two")).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two <<<"),
				Contains("+TWO <<<"),
			).
			Press(keys.Universal.CopyToClipboard)

		// Both kinds of line are in the selection, so the columns stay: what comes out
		// is the diff itself.
		expectClipboard(t, Equals("-two\n+TWO\n"))

		// The pane showing the custom patch is a diff too, of the patch's own lines.
		t.Views().Main().
			PressPrimaryAction().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-two <<<"),
				Contains("+TWO <<<"),
			).
			Press(keys.Universal.CopyToClipboard)

		expectClipboard(t, Equals("-two\n+TWO\n"))
	},
})
