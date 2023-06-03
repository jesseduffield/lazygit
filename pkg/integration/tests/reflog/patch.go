package reflog

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Patch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Build a patch from a reflog commit and apply it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.CreateFileAndAdd("file1", "content1")
		shell.CreateFileAndAdd("file2", "content2")
		shell.Commit("three")
		shell.HardReset("HEAD^^")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().ReflogCommits().
			Focus().
			Lines(
				Contains("reset: moving to HEAD^^").IsSelected(),
				Contains("commit: three"),
				Contains("commit: two"),
				Contains("commit (initial): one"),
			).
			SelectNextItem().
			Lines(
				Contains("reset: moving to HEAD^^"),
				Contains("commit: three").IsSelected(),
				Contains("commit: two"),
				Contains("commit (initial): one"),
			).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
				Contains("file2"),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		t.Views().
			CommitFiles().
			Press(keys.Universal.CreatePatchOptionsMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Patch options")).
			Select(MatchesRegexp(`Apply patch$`)).Confirm()

		t.Views().Files().Lines(
			Contains("file1"),
		)
	},
})
