package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SymlinkIntoRepoSubdir = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in a symlink into a repo's subdirectory",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - repo/a/b/c (main worktree with subdirs)
		//  - link (symlink to repo/a/b/c)
		//
		shell.CreateFileAndAdd("a/b/c/blah", "blah\n")
		shell.Commit("initial commit")
		shell.UpdateFile("a/b/c/blah", "updated content\n")

		shell.Chdir("..")
		shell.RunCommand([]string{"ln", "-s", "repo/a/b/c", "link"})

		shell.Chdir("link")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
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
				Contains("initial commit"),
			)
	},
})
