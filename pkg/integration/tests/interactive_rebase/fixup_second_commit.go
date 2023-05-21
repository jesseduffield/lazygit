package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FixupSecondCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fixup the second commit into the first (initial)",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file1.txt", "File1 Content\n").Commit("First Commit").
			CreateFileAndAdd("file2.txt", "Fixup Content\n").Commit("Fixup Commit Message").
			CreateFileAndAdd("file3.txt", "File3 Content\n").Commit("Third Commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("Third Commit"),
				Contains("Fixup Commit Message"),
				Contains("First Commit"),
			).
			NavigateToLine(Contains("Fixup Commit Message")).
			Press(keys.Commits.MarkCommitAsFixup).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Fixup")).
					Content(Equals("Are you sure you want to 'fixup' this commit? It will be merged into the commit below")).
					Confirm()
			}).
			Lines(
				Contains("Third Commit"),
				Contains("First Commit").IsSelected(),
			)

		t.Views().Main().
			// Make sure that the resulting commit message doesn't contain the
			// message of the fixup commit; compare this to
			// squash_down_second_commit.go, where it does.
			Content(Contains("First Commit")).
			Content(DoesNotContain("Fixup Commit Message")).
			Content(Contains("+File1 Content")).
			Content(Contains("+Fixup Content"))
	},
})
