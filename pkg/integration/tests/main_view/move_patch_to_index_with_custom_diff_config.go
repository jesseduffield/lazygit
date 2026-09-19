package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchToIndexWithCustomDiffConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Moving a patch to the index works even if diff.noprefix or diff.external are set",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.Commit("first commit")

		shell.SetConfig("diff.noprefix", "true")
		shell.SetConfig("diff.external", "echo")
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
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+file1 content")).
			PressPrimaryAction()

		t.Views().Secondary().Content(Contains("+file1 content"))
		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Views().CommitFiles().Lines(Equals("(none)"))
		t.Views().Files().Lines(
			Contains("A").Contains("file1"),
		)
	},
})
