package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var LockErrorPartwayThroughSkip = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Skip in a rebase while another git process takes index.lock before the next commit is committed",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "original\n").
			Commit("base").
			UpdateFileAndAdd("file", "changed\n").
			Commit("commit to edit").
			CreateFileAndAdd("file-y", "y\n").
			Commit("add file-y").
			CreateFileAndAdd("file-z", "z\n").
			Commit("add file-z")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("add file-z").IsSelected(),
				Contains("add file-y"),
				Contains("commit to edit"),
				Contains("base"),
			).
			NavigateToLine(Contains("commit to edit")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("add file-z"),
				Contains("pick").Contains("add file-y"),
				Contains("─── Commits"),
				Contains("commit to edit").IsSelected(),
				Contains("base"),
			)

		// After amending the commit we stopped at, git has to make a new commit
		// when it picks "add file-y", rather than fast-forwarding to it. The
		// lock then stops it after it has applied the changes, but before it
		// has committed them.
		t.Shell().RunCommand([]string{"git", "commit", "--amend", "-m", "edited commit"})
		shared.HoldIndexLockBeforeNextPickIsCommitted(t.Shell())

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Rebase options")).
			Select(Contains("skip")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("index.lock': File exists.")).
			Confirm()

		t.Views().Files().
			Lines(
				Contains("A").Contains("file-y"),
			)
	},
})
