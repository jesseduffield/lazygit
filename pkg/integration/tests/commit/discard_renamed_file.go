package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardRenamedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discard a renamed file from an old commit; both the new and the old path are handled so the rename is undone",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("original", "line1\nline2\nline3\nline4\nline5\n")
		shell.CreateFileAndAdd("other", "other content\n")
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
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Discard file changes")).
			Content(Contains("Are you sure you want to discard changes to the selected file(s) from this commit?")).
			Confirm()

		// The rename is undone: the commit no longer touches any file. (If only
		// the new path were discarded, the commit would still delete the old
		// path and show "D original" here instead.)
		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("(none)"),
			).
			PressEscape()

		// The working tree is clean; the original file is back at HEAD.
		t.Views().Files().
			IsEmpty()
	},
})
