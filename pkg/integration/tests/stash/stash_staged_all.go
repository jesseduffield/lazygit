package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashStagedAll = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stash staged (all files and content) changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-staged", "content")
		shell.EmptyCommit("initial commit")
		shell.UpdateFileAndAdd("file-staged", "more content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("file-staged"),
			).
			Press(keys.Files.ViewStashOptions)

		t.ExpectPopup().Menu().Title(Equals("Stash options")).Select(MatchesRegexp("Stash staged changes$")).Confirm()

		// in the previous implementation if you give a name to the stash entry it whould work
		// that's why here I'm specifically not giving one
		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("").Confirm()

		t.Views().Stash().
			LineCount(EqualsInt(1)).    // I check that there's only one line here
			Lines(MatchesRegexp("WIP")) // I didn't specify the file so cannot check the content.

		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file-staged").IsSelected(),
			)
	},
})
