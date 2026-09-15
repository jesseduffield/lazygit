package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FixupKeepMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fixup a commit, keeping its commit message",
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
			NavigateToLine(Contains("Second Commit")).
			Press(keys.Commits.MarkCommitAsFixup).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Fixup")).
					Select(Contains("use this commit's message")).
					Confirm()
			}).
			Lines(
				Contains("Third Commit"),
				Contains("Second Commit").IsSelected(),
			)

		t.Views().Main().
			// The resulting commit should have the message from the fixup commit,
			// not the target commit
			Content(Contains("Second Commit")).
			Content(DoesNotContain("First Commit")).
			Content(Contains("+File1 Content")).
			Content(Contains("+File2 Content"))
	},
})
