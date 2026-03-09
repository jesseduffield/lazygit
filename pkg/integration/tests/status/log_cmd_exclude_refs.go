package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LogCmdExcludeRefs = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Exclude configured refs from the all-branches log command",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AllBranchesLogCmds = []string{"git log --all --pretty=%s"}
		config.GetUserConfig().Git.AllBranchesLogExcludeRefs = []string{"refs/jj/*"}
		config.GetUserConfig().Gui.StatusPanelView = "allBranchesLog"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base commit").
			NewBranch("hidden").
			EmptyCommit("hidden jj commit").
			RunCommand([]string{"git", "update-ref", "refs/jj/test", "HEAD"}).
			Checkout("master").
			RunCommand([]string{"git", "branch", "-D", "hidden"}).
			EmptyCommit("visible main commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus()

		t.Views().Main().
			Content(Contains("visible main commit").DoesNotContain("hidden jj commit"))
	},
})
