package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CancelCommitDuringSquashMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open the commit panel during a squash merge started outside lazygit, and close it without changing git's message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(cfg *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one").
			NewBranch("feature-branch").
			CreateFileAndAdd("feature-file", "content").
			Commit("feature").
			Checkout("master").
			RunCommand([]string{"git", "merge", "--squash", "feature-branch"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Equals("A  feature-file"),
			).
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("Squashed commit of the following:")).
			Cancel()

		t.FileSystem().PathNotPresent(".git/LAZYGIT_PENDING_COMMIT")
	},
})
