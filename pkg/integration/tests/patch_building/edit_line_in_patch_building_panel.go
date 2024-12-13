package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditLineInPatchBuildingPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit a line in the patch building panel; make sure we end up on the right line",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.EditAtLine = "echo {{filename}}:{{line}} > edit-command"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file.txt", "4\n5\n6\n")
		shell.Commit("01")
		shell.UpdateFileAndAdd("file.txt", "1\n2a\n2b\n3\n4\n5\n6\n")
		shell.Commit("02")
		shell.UpdateFile("file.txt", "1\n2\n3\n4\n5\n6\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("02").IsSelected(),
				Contains("01"),
			).
			Press(keys.Universal.NextItem).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("A file.txt").IsSelected(),
			).
			PressEnter()

		t.Views().PatchBuilding().
			IsFocused().
			Content(Contains("+4\n+5\n+6")).
			NavigateToLine(Contains("+5")).
			Press(keys.Universal.Edit)

		t.FileSystem().FileContent("edit-command", Contains("file.txt:5\n"))
	},
})
