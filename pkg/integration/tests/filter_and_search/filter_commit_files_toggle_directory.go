package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterCommitFilesToggleDirectory = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle a filtered directory for a custom patch only adds visible files",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir1")
		shell.CreateFileAndAdd("dir1/apple-grape", "apple-grape content\n")
		shell.CreateFileAndAdd("dir1/apple-orange", "apple-orange content\n")
		shell.CreateFileAndAdd("dir1/grape-orange", "grape-orange content\n")
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
				Contains("dir1").IsSelected(),
				Contains("apple-grape"),
				Contains("apple-orange"),
				Contains("grape-orange"),
			).
			// Filter to show only "apple" files (staying in tree view)
			FilterOrSearch("apple").
			Lines(
				// first item is always selected after filtering
				Contains("dir1").IsSelected(),
				Contains("apple-grape"),
				Contains("apple-orange"),
			).
			// dir1 is already selected; toggle for patch
			PressPrimaryAction().
			Lines(
				Contains("dir1").IsSelected(),
				Contains("● apple-grape"),
				Contains("● apple-orange"),
			)

		t.Views().Information().Content(Contains("Building patch"))

		// Verify only the filtered files are in the patch (not grape-orange)
		t.Views().Secondary().Content(
			Contains("apple-grape").
				Contains("apple-orange").
				DoesNotContain("grape-orange"),
		)
	},
})
