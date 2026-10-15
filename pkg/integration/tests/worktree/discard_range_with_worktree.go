package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiscardRangeWithWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Discarding a range of files that includes a worktree inside the repo is not supported, because a worktree has its own menu for removing it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("mybranch")
		shell.CreateFileAndAdd("README.md", "hello world")
		shell.Commit("initial commit")
		shell.AddWorktree("mybranch", "nested-worktree", "newbranch")
		shell.CreateFile("untracked-file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? nested-worktree (worktree)").IsSelected(),
				Equals("?? untracked-file"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: Multiselection not supported for worktrees"))
			})
	},
})
