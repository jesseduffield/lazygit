package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitAllWithNestedRepository = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Committing with nothing staged doesn't offer to stage all files if that would stage a nested repository",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CreateFile("file", "content")
		shell.RunCommand([]string{"git", "init", "nested"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? file").IsSelected(),
				Equals("?? nested (nested repository)"),
			).
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Equals("Can't stage 'nested' because it is a git repository of its own. Discard it with d, or ignore it with i.")).
					Confirm()
			}).
			// Once something is staged, committing works as usual
			PressPrimaryAction().
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Type("add file").
					Confirm()
			}).
			Lines(
				Equals("?? nested (nested repository)").IsSelected(),
			)

		t.Views().Commits().
			Lines(
				Contains("add file"),
				Contains("first commit"),
			)
	},
})
