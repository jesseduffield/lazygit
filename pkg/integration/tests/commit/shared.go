package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func setupForAmendTests(shell *Shell) {
	shell.EmptyCommit("base commit")
	shell.NewBranch("branch")
	shell.Checkout("master")
	shell.CreateFileAndAdd("file1", "master")
	shell.Commit("file1 changed in master")
	shell.Checkout("branch")
	shell.UpdateFileAndAdd("file2", "two")
	shell.Commit("commit two")
	shell.CreateFileAndAdd("file1", "branch")
	shell.Commit("file1 changed in branch")
	shell.UpdateFileAndAdd("file3", "three")
	shell.Commit("commit three")
}

func doTheRebaseForAmendTests(t *TestDriver, keys config.KeybindingConfig) {
	t.Views().Commits().
		Focus().
		Lines(
			Contains("commit three").IsSelected(),
			Contains("file1 changed in branch"),
			Contains("commit two"),
			Contains("base commit"),
		)
	t.Views().Branches().
		Focus().
		NavigateToLine(Contains("master")).
		Press(keys.Branches.RebaseBranch).
		Tap(func() {
			t.ExpectPopup().Menu().
				Title(Equals("Rebase 'branch'")).
				Select(Contains("Simple rebase")).
				Confirm()
			t.Common().AcknowledgeConflicts()
		})

	t.Views().Commits().
		Lines(
			Contains("pick").Contains("commit three"),
			Contains("conflict").Contains("<-- YOU ARE HERE --- file1 changed in branch"),
			Contains("commit two"),
			Contains("file1 changed in master"),
			Contains("base commit"),
		)

	t.Views().Files().
		Focus().
		PressEnter()

	t.Views().MergeConflicts().
		IsFocused().
		SelectNextItem(). // choose "incoming"
		PressPrimaryAction()

	t.ExpectPopup().Confirmation().
		Title(Equals("Continue")).
		Content(Contains("All merge conflicts resolved. Continue?")).
		Cancel()
}

func checkCommitContainsChange(t *TestDriver, commitSubject string, change string) {
	t.Views().Commits().
		Focus().
		NavigateToLine(Contains(commitSubject))
	t.Views().Main().
		Content(Contains(change))
}
