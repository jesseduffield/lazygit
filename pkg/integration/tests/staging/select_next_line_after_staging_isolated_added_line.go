package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Tests that after staging an isolated addition (one that is alone in its block of changes), the
// cursor stays at the first change of the next block of changes which moves up to the same line,
// even if that block starts with a deletion.
var SelectNextLineAfterStagingIsolatedAddedLine = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "After staging an isolated added line, the cursor advances to the next hunk's first change",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n5\n6\n7\n8\n9\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "1\n2\n3\nnew\n4\n5\n6\n7b\n8\n9\n")
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
				Contains(" 2"),
				Contains(" 3"),
				Contains("+new"),
				Contains(" 4"),
				Contains(" 5"),
				Contains(" 6"),
				Contains("-7"),
				Contains("+7b"),
				Contains(" 8"),
				Contains(" 9"),
			).
			SelectedLine(Contains("+new")).
			PressPrimaryAction().
			SelectedLine(Contains("-7"))
	},
})
