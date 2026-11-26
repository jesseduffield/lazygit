package am

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var LockErrorPartwayThroughSkip = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Skip a conflicting patch while another git process takes index.lock before git am applies the patch after the next one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "original\n").
			Commit("base").
			NewBranch("source").
			UpdateFileAndAdd("file", "source\n").
			Commit("change file").
			CreateFileAndAdd("file2", "two\n").
			Commit("add file2").
			CreateFileAndAdd("file3", "three\n").
			Commit("add file3").
			RunShellCommand("git format-patch --stdout master..source > ../patches.mbox").
			Checkout("master").
			UpdateFileAndAdd("file", "master\n").
			Commit("conflicting change").
			RunCommandExpectError([]string{"git", "am", "--3way", "../patches.mbox"})

		// After skipping "change file", git am commits "add file2"; the lock
		// then stops it when it tries to apply "add file3".
		shared.HoldIndexLockAfterNextAppliedPatch(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Applying patches")).
			Select(Contains("skip")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("index.lock': File exists.")).
			Confirm()

		t.Views().Status().Content(Equals("(applying patches) repo → master"))
	},
})
