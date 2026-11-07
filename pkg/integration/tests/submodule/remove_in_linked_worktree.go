package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveInLinkedWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a submodule in a linked worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		shell.AddWorktree("HEAD", "../linked-worktree", "mybranch")
		shell.Chdir("../linked-worktree")
		shell.RunCommand([]string{"git", "submodule", "update", "--init"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Git keeps the git dirs of submodules separately for each worktree
		mainGitDirSubmodulePath := "../repo/.git/modules/my_submodule_name"
		linkedGitDirSubmodulePath := "../repo/.git/worktrees/linked-worktree/modules/my_submodule_name"
		t.FileSystem().
			PathPresent(mainGitDirSubmodulePath).
			PathPresent(linkedGitDirSubmodulePath)

		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule_name").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Remove submodule")).
					Content(Equals("Are you sure you want to remove submodule 'my_submodule_name' and its corresponding directory? This is irreversible.")).
					Confirm()
			}).
			IsEmpty()

		t.FileSystem().
			PathPresent(mainGitDirSubmodulePath).
			PathNotPresent(linkedGitDirSubmodulePath)
	},
})
