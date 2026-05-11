package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PreventDiscardingFileChanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check that it is not allowed to discard changes to a file of a stash",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("file", "content")
		shell.GitAddAll()
		shell.Stash("stash one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash one").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file").IsSelected(),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Confirmation().
			Title(Equals("Error")).
			Content(Contains("Changes can only be discarded from local commits")).
			Confirm()
	},
})
