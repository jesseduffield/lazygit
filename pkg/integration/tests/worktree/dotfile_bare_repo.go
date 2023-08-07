package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Can't think of a better name than 'dotfile' repo: I'm using that
// because that's the case we're typically dealing with.

var DotfileBareRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in the worktree of a dotfile bare repo and add a file and commit",
	ExtraCmdArgs: []string{"--git-dir={{.actualPath}}/.bare", "--work-tree={{.actualPath}}/repo"},
	Skip:         false,
	// passing this because we're explicitly passing --git-dir and --work-tree args
	UseCustomPath: true,
	SetupConfig:   func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - .bare
		//  - repo (the worktree)
		//
		// The first repo is called 'repo' because that's the
		// directory that all lazygit tests start in

		// Delete the .git dir that all tests start with by default
		shell.DeleteFile(".git")

		// Create a bare repo in the parent directory
		shell.RunCommand([]string{"git", "init", "--bare", "../.bare"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "checkout", "-b", "mybranch"})
		shell.CreateFile("blah", "original content\n")

		// Add a file and commit
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "add", "blah"})
		shell.RunCommand([]string{"git", "--git-dir=../.bare", "--work-tree=.", "commit", "-m", "initial commit"})

		shell.UpdateFile("blah", "updated content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("mybranch"),
			)

		t.Views().Commits().
			Lines(
				Contains("initial commit"),
			)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains(" M blah"), // shows as modified
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
				Contains("initial commit"),
			)
	},
})
