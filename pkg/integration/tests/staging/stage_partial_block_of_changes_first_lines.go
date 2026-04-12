package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StagePartialBlockOfChangesFirstLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage only the first few lines of a block of consecutive changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n7\n8\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "1\n2b\n3b\n4b\n5b\n6b\n7b\n8\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			ContainsLines(
				Contains(" 1"),
				Contains("-2"),
				Contains("-3"),
				Contains("-4"),
				Contains("-5"),
				Contains("-6"),
				Contains("-7"),
				Contains("+2b"),
				Contains("+3b"),
				Contains("+4b"),
				Contains("+5b"),
				Contains("+6b"),
				Contains("+7b"),
				Contains(" 8"),
			).
			SelectedLines(Contains("-2")).
			PressPrimaryAction().
			SelectedLines(Contains("-3")).
			PressPrimaryAction().
			NavigateToLine(Contains("+2b")).
			PressPrimaryAction().
			SelectedLines(Contains("+3b")).
			PressPrimaryAction()

		t.Views().StagingSecondary().
			ContainsLines(
				Contains(" 1"),
				Contains("-2"),
				Contains("-3"),
				Contains("+2b"),
				Contains("+3b"),
				Contains(" 4"),
				Contains(" 5"),
				Contains(" 6"),
			)
	},
})
