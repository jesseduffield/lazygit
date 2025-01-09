package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StashStagedPartialFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stash staged changes when a file is partially staged",
	ExtraCmdArgs: []string{},
	GitVersion:   AtLeast("git version 2.35.0"),
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file-staged", "line1\nline2\nline3\nline4\n")
		shell.Commit("initial commit")
		shell.UpdateFile("file-staged", "line1\nline2 mod\nline3\nline4 mod\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			PressEnter()

		t.Views().Staging().
			Content(
				Contains(" line1\n-line2\n+line2 mod\n line3\n-line4\n+line4 mod"),
			).
			PressPrimaryAction().
			PressPrimaryAction().
			Content(
				Contains(" line1\n line2 mod\n line3\n-line4\n+line4 mod"),
			).
			PressEscape()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.ViewStashOptions)

		t.ExpectPopup().Menu().Title(Equals("Stash options")).Select(MatchesRegexp("Stash staged changes$")).Confirm()

		t.ExpectPopup().Prompt().Title(Equals("Stash changes")).Type("my stashed file").Confirm()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("my stashed file"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file-staged").IsSelected(),
			)
		t.Views().Main().
			Content(
				Contains(" line1\n-line2\n+line2 mod\n line3\n line4"),
			)

		t.Views().Files().
			Lines(
				Contains("file-staged"),
			)

		t.Views().Staging().
			Content(
				Contains(" line1\n line2\n line3\n-line4\n+line4 mod"),
			)
	},
})
