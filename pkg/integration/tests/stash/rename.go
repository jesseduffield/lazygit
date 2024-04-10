package stash

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var Rename = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try to rename the stash.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("blah").
			CreateFileAndAdd("file-1", "change to stash1").
			Stash("foo").
			CreateFileAndAdd("file-2", "change to stash2").
			Stash("bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			Focus().
			Lines(
				Contains("On master: bar"),
				Contains("On master: foo"),
			).
			SelectNextItem().
			Press(keys.Stash.RenameStash).
			Tap(func() {
				t.ExpectPopup().Prompt().Title(Equals("Rename stash: stash@{1}")).Type(" baz").Confirm()
			}).
			SelectedLine(Contains("On master: foo baz"))
	},
})
