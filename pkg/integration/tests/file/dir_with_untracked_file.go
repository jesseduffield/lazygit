package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DirWithUntrackedFile = NewIntegrationTest(NewIntegrationTestArgs{
	// notably, we currently _don't_ actually see the untracked file in the diff. Not sure how to get around that.
	Description:  "When selecting a directory that contains an untracked file, we should not get an error",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.ShowFileTree = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("dir")
		shell.CreateFile("dir/file", "foo")
		shell.GitAddAll()
		shell.Commit("first commit")
		shell.CreateFile("dir/untracked", "bar")
		shell.UpdateFile("dir/file", "baz")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("first commit"),
			)

		t.Views().Main().
			Content(DoesNotContain("error: Could not access")).
			// we show baz because it's a modified file but we don't show bar because it's untracked
			// (though it would be cool if we could show that too)
			Content(Contains("baz"))
	},
})
