package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExcludeFileInWorktree = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add a file to .git/info/exclude in a worktree",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit1")
		shell.AddWorktree("HEAD", "../linked-worktree", "mybranch")
		shell.CreateFile("../linked-worktree/toExclude", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().
			Focus().
			Lines(
				Contains("repo (main)").IsSelected(),
				Contains("linked-worktree"),
			).
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Files().
			Focus().
			Lines(
				Contains("toExclude"),
			).
			Press(keys.Files.IgnoreFile).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Ignore or exclude file")).Select(Contains("Add to .git/info/exclude")).Confirm()
			}).
			IsEmpty()

		t.FileSystem().FileContent("../repo/.git/info/exclude", Contains("toExclude"))
	},
})
