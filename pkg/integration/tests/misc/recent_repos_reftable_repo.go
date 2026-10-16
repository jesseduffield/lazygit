package misc

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RecentReposReftableRepo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The recent repositories menu shows the branch of a repo that keeps its refs in a reftable",
	ExtraCmdArgs: []string{},
	ExtraEnvVars: map[string]string{
		"SHOW_RECENT_REPOS": "true",
	},
	Skip:       false,
	GitVersion: AtLeast("2.45.0"),
	SetupConfig: func(cfg *config.AppConfig) {
		// the first entry is the repo we're in, so it isn't offered
		current, _ := filepath.Abs(".")
		reftable, _ := filepath.Abs("../reftable")
		unborn, _ := filepath.Abs("../reftable-unborn")
		cfg.GetAppState().RecentRepos = []string{current, reftable, unborn}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.RunCommand([]string{"git", "clone", "--ref-format=reftable", ".", "../reftable"})
		// A branch without a commit is the one thing git can't answer for with
		// rev-parse
		shell.RunCommand([]string{"git", "init", "--ref-format=reftable", "../reftable-unborn"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.ExpectPopup().Menu().
			Title(Equals("Recent repositories")).
			Lines(
				Contains("reftable").Contains("master").IsSelected(),
				Contains("reftable-unborn").Contains("master"),
				Contains("Cancel"),
			).
			Cancel()
	},
})
