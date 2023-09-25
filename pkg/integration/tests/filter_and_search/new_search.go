package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// This is a regression test to ensure https://github.com/jesseduffield/lazygit/issues/2971
// doesn't happen again

var NewSearch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Start a new search and verify the search begins from the current cursor position, not from the current search match",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// need to create some branches, each with their own commits
		shell.EmptyCommit("Add foo")
		shell.EmptyCommit("Remove foo")
		shell.EmptyCommit("Add bar")
		shell.EmptyCommit("Remove bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains(`Remove bar`).IsSelected(),
				Contains(`Add bar`),
				Contains(`Remove foo`),
				Contains(`Add foo`),
			).
			FilterOrSearch("Add").
			SelectedLine(Contains(`Add bar`)).
			SelectPreviousItem().
			SelectedLine(Contains(`Remove bar`)).
			FilterOrSearch("Remove").
			SelectedLine(Contains(`Remove bar`))
	},
})
