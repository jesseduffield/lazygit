package filter_and_search

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilterBranchesRegexp = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter local branches with re: prefix so ^ anchors to start of name",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial")
		shell.NewBranch("main")
		shell.EmptyCommit("on-main-branch")
		shell.Checkout("master")
		shell.NewBranch("amain")
		shell.EmptyCommit("on-amain-branch")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Content(Contains(`amain`)).
			Content(Contains(`main`)).
			Content(Contains(`master`)).
			FilterOrSearch("re:^main").
			Lines(
				Contains(`main`).IsSelected(),
			).
			PressEscape().
			Content(Contains(`amain`)).
			Content(Contains(`main`)).
			Content(Contains(`master`))
	},
})
