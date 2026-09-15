package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDirtyOnly = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing space on a submodule that only has dirty content (no new commit) can't stage anything, so we explain that with an error instead of silently doing nothing.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		// Dirty working-tree content, but no new commit: there's nothing the
		// parent repo can stage.
		shell.CreateFile("my_submodule_path/dirty_file", "dirty content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			).
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("Nothing to stage")).
					Confirm()
			}).
			// The status is unchanged: nothing got staged.
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			).
			// Pressing "stage all" must behave the same way.
			Press(keys.Files.ToggleStagedAll).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("Nothing to stage")).
					Confirm()
			}).
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			)
	},
})
