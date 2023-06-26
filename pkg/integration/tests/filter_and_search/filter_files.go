package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Basic file filtering by text",
	ExtraCmdArgs: []string{},
	Skip:         true, // Skipping until we have implemented file view filtering
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("folder1")
		shell.CreateFile("folder1/apple-grape", "apple-grape")
		shell.CreateFile("folder1/apple-orange", "apple-orange")
		shell.CreateFile("folder1/grape-orange", "grape-orange")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Contains(`folder1`).IsSelected(),
				Contains(`apple-grape`),
				Contains(`apple-orange`),
				Contains(`grape-orange`),
			).
			Press(keys.Files.ToggleTreeView).
			Lines(
				Contains(`folder1/apple-grape`).IsSelected(),
				Contains(`folder1/apple-orange`),
				Contains(`folder1/grape-orange`),
			).
			FilterOrSearch("apple").
			Lines(
				Contains(`folder1/apple-grape`).IsSelected(),
				Contains(`folder1/apple-orange`),
			).
			Press(keys.Files.ToggleTreeView).
			// filter still applies when we toggle tree view
			Lines(
				Contains(`folder1`),
				Contains(`apple-grape`).IsSelected(),
				Contains(`apple-orange`),
			).
			Press(keys.Files.ToggleTreeView).
			Lines(
				Contains(`folder1/apple-grape`).IsSelected(),
				Contains(`folder1/apple-orange`),
			).
			NavigateToLine(Contains(`folder1/apple-orange`)).
			Press(keys.Universal.Return).
			Lines(
				Contains(`folder1/apple-grape`),
				// selection is retained after escaping filter mode
				Contains(`folder1/apple-orange`).IsSelected(),
				Contains(`folder1/grape-orange`),
			).
			Tap(func() {
				t.Views().Search().IsInvisible()
			}).
			Press(keys.Files.ToggleTreeView).
			Lines(
				Contains(`folder1`),
				Contains(`apple-grape`),
				Contains(`apple-orange`).IsSelected(),
				Contains(`grape-orange`),
			).
			FilterOrSearch("folder1/grape").
			Lines(
				// first item is always selected after filtering
				Contains(`folder1`).IsSelected(),
				Contains(`grape-orange`),
			)
	},
})
