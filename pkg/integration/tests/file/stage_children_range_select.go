package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageChildrenRangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a range of files/folders and their children using range select",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("foo", "")
		shell.CreateFile("foobar", "")
		shell.CreateFile("baz/file", "")
		shell.CreateFile("bazbam/file", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("▼ baz").IsSelected(),
				Contains("  ??").Contains("file"),
				Contains("▼ bazbam"),
				Contains("  ??").Contains("file"),
				Contains("??").Contains("foo"),
				Contains("??").Contains("foobar"),
			).
			// Select everything
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("foobar")).
			// Stage
			PressPrimaryAction().
			Lines(
				Contains("▼ baz").IsSelected(),
				Contains("  A ").Contains("file").IsSelected(),
				Contains("▼ bazbam").IsSelected(),
				Contains("  A ").Contains("file").IsSelected(),
				Contains("A ").Contains("foo").IsSelected(),
				Contains("A ").Contains("foobar").IsSelected(),
			)
	},
})
