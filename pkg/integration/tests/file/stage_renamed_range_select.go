package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRenamedRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a range of renamed files/folders using range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("dir1/file-a", "A's content")
		shell.CreateFileAndAdd("file-b", "B's content")
		shell.Commit("first commit")
		shell.Rename("dir1", "dir1_v2")
		shell.Rename("file-b", "file-b_v2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ dir1").IsSelected(),
				Contains("   D").Contains("file-a"),
				Contains("▼ dir1_v2"),
				Contains("  ??").Contains("file-a"),
				Contains(" D").Contains("file-b"),
				Contains("??").Contains("file-b_v2"),
			).
			// Select everything
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("file-b_v2")).
			// Stage
			PressPrimaryAction().
			Lines(
				Contains("▼ dir1_v2"),
				Contains("  R ").Contains("dir1/file-a → file-a"),
				Contains("R ").Contains("file-b → file-b_v2").IsSelected(),
			)
	},
})
