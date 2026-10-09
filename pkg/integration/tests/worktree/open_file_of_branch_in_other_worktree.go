package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var OpenFileOfBranchInOtherWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open a file of a commit of a branch that is checked out in another worktree, and copy its absolute path",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.Open = "echo {{filename}} > open-command"
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("file.txt", "content")
		shell.Commit("add file")
		shell.Checkout("master")
		shell.AddWorktreeCheckout("mybranch", "../linked-worktree")
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
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file.txt").IsSelected(),
			).
			Press(keys.Universal.OpenFile).
			Tap(func() {
				/* EXPECTED:
				t.FileSystem().FileContent("open-command", Contains("/linked-worktree/file.txt\n"))
				ACTUAL: */
				t.FileSystem().FileContent("open-command", Contains("/repo/file.txt\n"))
			}).
			Press(keys.Files.CopyFileInfoToClipboard).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Copy to clipboard")).
					Select(Contains("Absolute path")).
					Confirm()

				t.ExpectToast(Equals("File path copied to clipboard"))

				/* EXPECTED:
				t.FileSystem().FileContent("clipboard", Contains("/linked-worktree/file.txt"))
				ACTUAL: */
				t.FileSystem().FileContent("clipboard", Contains("/repo/file.txt"))
			})
	},
})
