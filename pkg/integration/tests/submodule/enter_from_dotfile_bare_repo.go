package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Entering a submodule and escaping back out again, in a repo that git can only
// find because we were told where it is (--git-dir/--work-tree). Entering the
// submodule has to leave that behind, since it says where the superproject is,
// so coming back out has to bring it along again.

var EnterFromDotfileBareRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter a submodule of a dotfile bare repo and escape back out again",
	ExtraCmdArgs: []string{"--git-dir={{.actualPath}}/.bare", "--work-tree={{.actualPath}}/repo"},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - .bare (the git dir)
		//  - repo (the work tree, with no .git of its own)
		//  - my_submodule_name (the submodule's remote)
		//
		// The work tree is called 'repo' because that's the directory that all
		// lazygit tests start in

		// make a repo for the submodule to be cloned from, using the .git dir
		// that every test starts with
		shell.EmptyCommit("initial submodule commit")
		shell.Clone("my_submodule_name")

		// now turn the test repo into a dotfile-style bare repo
		shell.DeleteFile(".git")
		shell.RunCommand([]string{"git", "init", "--bare", "../.bare"})
		gitInBareRepo := []string{"git", "--git-dir=../.bare", "--work-tree=."}
		shell.RunCommand(append(gitInBareRepo, "checkout", "-b", "mybranch"))
		shell.CreateFile("blah", "blah\n")
		shell.RunCommand(append(gitInBareRepo, "add", "blah"))
		shell.RunCommand(append(gitInBareRepo, "commit", "-m", "initial commit"))
		shell.RunCommand(append(gitInBareRepo, "-c", "protocol.file.allow=always", "submodule",
			"add", "--name", "my_submodule_name", "../my_submodule_name", "my_submodule_path"))
		shell.RunCommand(append(gitInBareRepo, "commit", "-m", "add submodule"))
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		assertInParentRepo := func() {
			t.Views().Status().Content(Contains("repo"))
			t.Views().Commits().Lines(
				Contains("add submodule"),
				Contains("initial commit"),
			)
		}

		assertInParentRepo()

		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule_name").IsSelected(),
			).
			PressEnter()

		t.Views().Status().Content(Contains("my_submodule_path"))
		t.Views().Commits().Lines(
			Contains("initial submodule commit"),
		)

		t.Views().Files().IsFocused().PressEscape()

		assertInParentRepo()
		t.Views().Submodules().IsFocused()
	},
})
