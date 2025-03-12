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
				Equals("▼ /").IsSelected(),
				Equals("  ▼ baz"),
				Equals("    ?? file"),
				Equals("  ▼ bazbam"),
				Equals("    ?? file"),
				Equals("  ?? foo"),
				Equals("  ?? foobar"),
			).
			// Select everything
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("foobar")).
			// Stage
			PressPrimaryAction().
			Lines(
				Equals("▼ /").IsSelected(),
				Equals("  ▼ baz").IsSelected(),
				Equals("    A  file").IsSelected(),
				Equals("  ▼ bazbam").IsSelected(),
				Equals("    A  file").IsSelected(),
				Equals("  A  foo").IsSelected(),
				Equals("  A  foobar").IsSelected(),
			)
	},
})
