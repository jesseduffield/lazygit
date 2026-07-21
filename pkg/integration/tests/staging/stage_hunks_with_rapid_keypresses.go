package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// The second space is pressed before the refresh triggered by the first one
// has updated the staging panel. That refresh is what moves the selection to
// the next hunk, so the second press must not be handled until it has landed;
// handling it earlier would try to stage the first hunk a second time.
var StageHunksWithRapidKeypresses = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage two hunks with two space presses in rapid succession",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
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
			PressRapidly(keys.Universal.Select, keys.Universal.Select)

		t.Views().StagingSecondary().
			IsFocused().
			ContainsLines(
				Contains("+1b"),
				Contains("+2b"),
			).
			ContainsLines(
				Contains("+3b"),
				Contains("+4b"),
			)
	},
})
