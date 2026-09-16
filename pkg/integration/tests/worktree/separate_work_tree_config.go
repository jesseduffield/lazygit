package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This case is like bare_repo_worktree_config.go, except that lazygit isn't
// told where the git dir is: it is started in the directory containing it, and
// finds it the way git does. The work tree is somewhere else entirely, so git
// can't find its way back from there, and every command we run has to be told
// where the repo is.

var SeparateWorkTreeConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in the git dir of a repo whose work tree is elsewhere, and add a file and commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - repo (holds the .git dir, and nothing else; lazygit starts here)
		//  - worktree (holds the files)
		//
		// 'repo' is the repository/directory that all lazygit tests start in

		shell.CreateFileAndAdd("blah", "original content\n")
		shell.Commit("initial commit")

		// point the repo at a work tree outside of it (core.worktree is
		// relative to the .git dir), and fill that work tree from HEAD
		shell.CreateDir("../worktree")
		shell.SetConfig("core.worktree", "../../worktree")
		shell.RunCommand([]string{"git", "reset", "--hard"})

		// the copy of the file we committed from is not in the work tree, so
		// git no longer knows anything about it
		shell.DeleteFile("blah")

		shell.UpdateFile("../worktree/blah", "updated content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
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
