package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Regression test for https://github.com/jesseduffield/lazygit/issues/1118:
// when lazygit is launched against a dotfile-style bare repo via
// --git-dir/--work-tree (the yadm/vcsh setup), entering a submodule and
// pressing escape must return to the parent repo. The parent repo is only
// discoverable through the GIT_DIR/GIT_WORK_TREE env vars, so escape fails
// with "not a git repository" if those aren't restored on the way back.
var EnterDotfileBareRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter a submodule of a dotfile bare repo (--git-dir/--work-tree) and escape back out",
	ExtraCmdArgs: []string{"--git-dir={{.actualPath}}/.bare", "--work-tree={{.actualPath}}/repo"},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// directory structure (like worktree/dotfile_bare_repo.go):
		// <root>
		//  - .bare (the git dir; only reachable via GIT_DIR env var)
		//  - repo (the worktree; has no .git entry at all)
		//  - my_submodule_name (bare clone serving as the submodule's remote)

		// create a repo to act as the submodule's remote, using the default
		// .git dir that every test repo starts with
		shell.EmptyCommit("submodule initial commit")
		shell.Clone("my_submodule_name")

		// now turn the test repo into a dotfile-style bare repo
		shell.DeleteFile(".git")
		shell.RunCommand([]string{"git", "init", "--bare", "../.bare"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "checkout", "-b", "mybranch"})
		shell.CreateFile("blah", "blah\n")
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "add", "blah"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "commit", "-m", "initial commit"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "-c", "protocol.file.allow=always", "submodule", "add", "--name", "my_submodule_name", "../my_submodule_name", "my_submodule_path"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "commit", "-m", "add submodule"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		assertInParentRepo := func() {
			t.Views().Status().Content(Contains("mybranch"))
		}
		assertInSubmodule := func() {
			t.Views().Status().Content(Contains("(my_submodule_name)"))
		}

		assertInParentRepo()

		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule_name").IsSelected(),
			).
			// enter the submodule
			PressEnter()

		assertInSubmodule()

		t.Views().Files().IsFocused().
			// return to the parent repo
			PressEscape()

		assertInParentRepo()

		t.Views().Submodules().IsFocused()
	},
})
