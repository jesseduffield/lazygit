package misc

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RecentReposBranchColumn = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The branch column of the recent repositories menu shows what each repo has checked out",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip: false,
	SetupConfig: func(cfg *config.AppConfig) {
		// the first entry is the repo we're in, so it isn't offered
		current, _ := filepath.Abs(".")
		onBranch, _ := filepath.Abs("../on-branch")
		detached, _ := filepath.Abs("../detached")
		submodule, _ := filepath.Abs("sub")
		cfg.GetAppState().RecentRepos = []string{current, onBranch, detached, submodule}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneNonBare("on-branch")
		shell.CloneNonBare("detached")
		shell.RunCommand([]string{"git", "-C", "../detached", "checkout", "--detach"})
		shell.CloneIntoSubmodule("submodule", "sub")
		shell.GitAddAll()
		shell.Commit("add submodule")
		shell.RunCommand([]string{"git", "-C", "sub", "checkout", "--detach"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Lines(
				Contains("on-branch").Contains("master").IsSelected(),
				Contains("detached").MatchesRegexp(`HEAD detached at [0-9a-f]{8}`),
				Contains("sub").MatchesRegexp(`HEAD detached at [0-9a-f]{8}`),
				Contains("Cancel"),
			).
			Cancel()
	},
})
