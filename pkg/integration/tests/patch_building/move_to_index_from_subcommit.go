package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToIndexFromSubcommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit of the current branch to the index",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.CreateFileAndAdd("file2", "file2 content\n")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master").IsSelected(),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
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
			SelectNextItem().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().Secondary().Content(Contains("+file1 content"))

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Views().Files().
			Lines(
				Contains("A").Contains("file1"),
			).
			Focus()

		t.Views().Main().
			Content(Contains("file1 content"))
	},
})
