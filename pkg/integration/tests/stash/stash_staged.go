package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashStaged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stash staged changes",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-staged", "content")
		shell.CreateFileAndAdd("file-unstaged", "content")
		shell.EmptyCommit("initial commit")
		shell.UpdateFileAndAdd("file-staged", "new content")
		shell.UpdateFile("file-unstaged", "new content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Stash().
			IsEmpty()

		t.Views().Files().
			Lines(
				Contains("file-staged"),
				Contains("file-unstaged"),
			).
			Press(keys.Files.ViewStashOptions)

		t.ExpectPopup().Menu().Title(Equals("Stash options")).Select(MatchesRegexp("stash staged changes$")).Confirm()

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Lines(
				Contains("my stashed file"),
			)

		t.Views().Files().
			Lines(
				Contains("file-unstaged"),
			)

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
