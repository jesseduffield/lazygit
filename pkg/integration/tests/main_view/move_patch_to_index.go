package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchToIndex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to the index",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit").IsSelected(),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("▼ /").IsSelected(),
				Contains("file1"),
				Contains("file2"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+file1 content")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().Content(Contains("+file1 content"))

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Views().Files().Lines(
			Contains("A").Contains("file1"),
		)
		t.Views().Main().
			IsFocused().
			Content(Contains("+file2 content"))
		t.Views().Commits().Lines(
			Contains("first commit").IsSelected(),
		)

		t.Views().Files().Focus()
		t.Views().Secondary().Content(Contains("file1 content"))
	},
})
