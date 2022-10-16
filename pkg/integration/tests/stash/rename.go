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
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToStashWindow()
		assert.CurrentViewName("stash")

		assert.MatchSelectedLine(Equals("On master: bar"))
		input.NextItem()
		assert.MatchSelectedLine(Equals("On master: foo"))
		input.PressKeys(keys.Stash.RenameStash)
		assert.InPrompt()
		assert.MatchCurrentViewTitle(Equals("Rename stash: stash@{1}"))

		input.Type(" baz")
		input.Confirm()

		assert.MatchSelectedLine(Equals("On master: foo baz"))
	},
})
