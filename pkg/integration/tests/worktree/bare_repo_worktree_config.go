package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This case is identical to dotfile_bare_repo.go, except
// that it invokes lazygit with $GIT_DIR set but not
// $GIT_WORK_TREE. Instead, the repo uses the core.worktree
// config to identify the main worktre.

var BareRepoWorktreeConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in the worktree of a vcsh-style bare repo and add a file and commit",
	ExtraCmdArgs: []string{"--git-dir={{.actualPath}}/.bare"},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - .bare
		//  - . (a worktree at the same path as .bare)
		//
		//
		// 'repo' is the repository/directory that all lazygit tests start in

		shell.CreateFileAndAdd("a/b/c/blah", "blah\n")
		shell.Commit("initial commit")

		shell.CreateFileAndAdd(".gitignore", ".bare/\n/repo\n")
		shell.Commit("add .gitignore")

		shell.Chdir("..")

		// configure this "fake bare"" repo using the vcsh convention
		// of core.bare=false and core.worktree set to the actual
		// worktree path (a homedir root). This allows $GIT_DIR
		// alone to make this repo "self worktree identifying"
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "init", "--shared=false"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "config", "core.bare", "false"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "config", "core.worktree", ".."})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "remote", "add", "origin", "./repo"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "checkout", "-b", "main"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "config", "branch.main.remote", "origin"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "config", "branch.main.merge", "refs/heads/master"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "fetch", "origin", "master"})
		shell.RunCommand([]string{"git", "--git-dir=./.bare", "-c", "merge.ff=true", "merge", "origin/master"})

		// we no longer need the original repo so remove it
		shell.DeleteFile("repo")

		shell.UpdateFile("a/b/c/blah", "updated content\n")
		shell.Chdir("a/b/c")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("main"),
			)

		t.Views().Commits().
			Lines(
				Contains("add .gitignore"),
				Contains("initial commit"),
			)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains(" M a/b/c/blah"), // shows as modified
			).
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("Add blah").
			Confirm()

		t.Views().Files().
			IsEmpty()

		t.Views().Commits().
			Lines(
				Contains("Add blah"),
				Contains("add .gitignore"),
				Contains("initial commit"),
			)
	},
})
