package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// note: this is required to simulate the clipboard during CI
func expectClipboard(t *TestDriver, matcher *TextMatcher) {
	defer t.Shell().DeleteFile("clipboard")

	t.FileSystem().FileContent("clipboard", matcher)
}

var CopyRenamedFileDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the diff of a renamed file to the clipboard; the diff shows the rename rather than a delete and add",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "line1\nline2\nline3\nline4\nline5\n")
		shell.Commit("first commit")

		shell.RenameFileInGit("original", "renamed")
		shell.UpdateFileAndAdd("renamed", "line1\nline2 changed\nline3\nline4\nline5\n")
		shell.Commit("rename with modification")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("rename with modification").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("original → renamed").IsSelected(),
			).
			Press(keys.Files.CopyFileInfoToClipboard)

		t.ExpectPopup().Menu().
			Title(Equals("Copy to clipboard")).
			Select(Contains("Diff of selected file")).
			Confirm()

		t.ExpectToast(Contains("File diff copied to clipboard"))

		expectClipboard(t,
			Contains("rename from original").
				Contains("rename to renamed").
				Contains("-line2").
				Contains("+line2 changed"),
		)
	},
})
