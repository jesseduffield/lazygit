package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExternalRemoveCurrentWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Recover when the current linked worktree is deleted externally",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
		shell.Chdir("../linked-worktree")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Confirm we're in the linked worktree
		t.Views().Status().
			Lines(
				Contains("repo(linked-worktree) → newbranch"),
			)

		// Simulate external deletion (e.g. another terminal runs rm -rf)
		t.Shell().RunShellCommand("rm -rf ../linked-worktree")

		// Trigger a refresh so lazygit detects the deleted CWD
		t.GlobalPress(keys.Universal.Refresh)

		t.ExpectToast(Contains("Worktree deleted externally"))

		// Lazygit should auto-switch to the main worktree
		t.Views().Status().
			Lines(
				Contains("repo → mybranch"),
			)
	},
})
