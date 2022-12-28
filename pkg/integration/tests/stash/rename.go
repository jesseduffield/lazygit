package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Rename = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to rename the stash.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("blah").
			CreateFileAndAdd("file-1", "change to stash1").
			StashWithMessage("foo").
			CreateFileAndAdd("file-2", "change to stash2").
			StashWithMessage("bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			Focus().
			Lines(
				Equals("On master: bar"),
				Equals("On master: foo"),
			).
			SelectNextItem().
			Press(keys.Stash.RenameStash).
			Tap(func() {
				t.ExpectPopup().Prompt().Title(Equals("Rename stash: stash@{1}")).Type(" baz").Confirm()
			}).
			SelectedLine(Equals("On master: foo baz"))
	},
})
