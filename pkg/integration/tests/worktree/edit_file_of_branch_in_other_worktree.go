package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditFileOfBranchInOtherWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit a file of a commit of a branch that is checked out in another worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().OS.Edit = "echo {{filename}} > edit-command"
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("file.txt", "1\n2\n3\n")
		shell.Commit("add file")
		shell.Checkout("master")
		shell.AddWorktreeCheckout("mybranch", "../linked-worktree")

		// Move the lines of the file down in the linked worktree, so that a line of
		// the commit's diff is at a different line of the file there
		shell.UpdateFile("../linked-worktree/file.txt", "0\n1\n2\n3\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
				Contains("mybranch (worktree linked-worktree)"),
			).
			NavigateToLine(Contains("mybranch")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("add file").IsSelected(),
				Contains("initial commit"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+2")).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/linked-worktree/file.txt:3\n"))
			}).
			PressEscape()

		t.Views().SubCommits().
			IsFocused().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file.txt").IsSelected(),
			).
			Press(keys.Universal.Edit).
			Tap(func() {
				t.FileSystem().FileContent("edit-command", Contains("/linked-worktree/file.txt\n"))
			})
	},
})
