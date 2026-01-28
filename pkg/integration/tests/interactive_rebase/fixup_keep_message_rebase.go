package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FixupKeepMessageRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Set fixup -C flag on a fixup commit during interactive rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1.txt", "File1 Content\n").Commit("First Commit").
			CreateFileAndAdd("file2.txt", "File2 Content\n").Commit("Second Commit").
			CreateFileAndAdd("file3.txt", "File3 Content\n").Commit("Third Commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Third Commit"),
				Contains("Second Commit"),
				Contains("First Commit"),
			).
			// Start interactive rebase from the first commit
			NavigateToLine(Contains("First Commit")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick CI Third Commit"),
				Contains("pick CI Second Commit"),
				Contains("--- Commits ---"),
				Contains("First Commit").IsSelected(),
			).
			// Mark second commit as fixup
			NavigateToLine(Contains("Second Commit")).
			Press(keys.Commits.MarkCommitAsFixup).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick  CI Third Commit"),
				Contains("fixup CI Second Commit").IsSelected(),
				Contains("--- Commits ---"),
				Contains("First Commit"),
			).
			// Now set the -C flag using the SetFixupMessage keybinding
			Press(keys.Commits.SetFixupMessage).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Set fixup message")).
					Select(Contains("use this commit's message")).
					Confirm()
			}).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick     CI Third Commit"),
				Contains("fixup -C CI Second Commit").IsSelected(),
				Contains("--- Commits ---"),
				Contains("First Commit"),
			).
			// Continue the rebase
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("Third Commit"),
				Contains("Second Commit").IsSelected(),
			)

		t.Views().Main().
			// The resulting commit should have the message from the fixup commit
			Content(Contains("Second Commit")).
			Content(DoesNotContain("First Commit")).
			Content(Contains("+File1 Content")).
			Content(Contains("+File2 Content"))
	},
})
