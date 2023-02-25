package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectAllFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "All all files of a commit to a custom patch with the 'a' keybinding",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "file1 content")
		shell.CreateFileAndAdd("file2", "file2 content")
		shell.CreateFileAndAdd("file3", "file3 content")
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
				Contains("file1").IsSelected(),
				Contains("file2"),
				Contains("file3"),
			).
			Press(keys.Files.ToggleStagedAll)

		t.Views().Information().Content(Contains("building patch"))

		t.Views().Secondary().Content(
			Contains("file1").Contains("file3").Contains("file3"),
		)
	},
})
