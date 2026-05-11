package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowDivergenceFromUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show divergence from upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content1")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file", "content2")
		shell.Commit("two")
		shell.CreateFileAndAdd("file3", "content3")
		shell.Commit("three")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("master", "origin/master")

		shell.HardReset("HEAD^^")
		shell.CreateFileAndAdd("file4", "content4")
		shell.Commit("four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("four"),
				Contains("one"),
			)

		t.Views().Branches().
			Focus().
			Lines(Contains("master")).
			Press(keys.Branches.SetUpstream)

		t.ExpectPopup().Menu().Title(Contains("Upstream")).Select(Contains("View divergence from upstream")).Confirm()

		t.Views().SubCommits().
			IsFocused().
			Title(Contains("Commits (master <-> origin/master)")).
			Lines(
				DoesNotContainAnyOf("↓", "↑").Contains("--- Remote ---"),
				Contains("↓").Contains("three"),
				Contains("↓").Contains("two"),
				DoesNotContainAnyOf("↓", "↑").Contains("--- Local ---"),
				Contains("↑").Contains("four"),
			)
	},
})
