package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ForceRemoveWorktreeWithSubmodules = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Force remove a worktree that contains submodules",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.CloneIntoSubmodule("submodule", "submodule")
		shell.Commit("Add submodule")
		shell.AddWorktree("mybranch", "../linked-worktree", "newbranch")
		shell.RunCommand([]string{"git", "-C", "../linked-worktree", "submodule", "update", "--init"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			).
			NavigateToLine(Contains("linked-worktree")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove worktree")).
					Content(Equals("Are you sure you want to remove worktree 'linked-worktree'?")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Remove worktree")).
					Content(Equals("'linked-worktree' contains modified or untracked files, or submodules (or all of these). Are you sure you want to remove it?")).
					Confirm()
			}).
			Lines(
				Contains("repo (main)").IsSelected(),
			)
	},
})
