package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Tests that after staging individual lines from a consecutive changes block,
// the cursor advances to the correct next change. The file has two separate
// hunks so that we can verify the cursor crosses hunk boundaries correctly.
var SelectNextLineAfterStagingInTwoHunkDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After staging lines from a two-hunk diff, the cursor advances correctly",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		// Use 7 context lines between the two change blocks so that git creates
		// two separate hunks.
		shell.CreateFileAndAdd("file1", "1\n2\na\nb\nc\nd\ne\nf\ng\n3\n4\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "1b\n2b\na\nb\nc\nd\ne\nf\ng\n3b\n4b\n")
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
				Contains("-1"),
				Contains("-2"),
				Contains("+1b"),
				Contains("+2b"),
				Contains(" a"),
				Contains(" b"),
				Contains(" c"),
				Contains("@@"),
				Contains(" e"),
				Contains(" f"),
				Contains(" g"),
				Contains("-3"),
				Contains("-4"),
				Contains("+3b"),
				Contains("+4b"),
			).
			NavigateToLine(Contains("-2")).
			PressPrimaryAction().
			SelectedLine(Contains("+1b")).
			PressPrimaryAction().
			SelectedLine(Contains("+2b")).
			PressPrimaryAction().
			SelectedLine(Contains("-3"))
	},
})
