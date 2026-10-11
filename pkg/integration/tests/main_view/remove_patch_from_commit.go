package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemovePatchFromCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Remove a whole-file custom patch from its original commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
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
		t.Common().SelectPatchOption(Contains("Remove patch from original commit"))

		t.Views().Files().IsEmpty()
		t.Views().Main().
			IsFocused().
			Content(Contains("+file2 content"))
		t.Views().Commits().Lines(
			Contains("first commit").IsSelected(),
		)
	},
})
