package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Stage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and unstage a submodule that has both a new commit and dirty content. The new commit can be staged, but the dirty content can't, so unstaging must still work; this must hold for both the stage (space) and stage-all (a) keybindings.",
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

		// Give the submodule a new commit, which is a change that the parent
		// repo can stage, as well as some dirty working-tree content, which
		// the parent repo can never stage. This is what gets us a "MM" status
		// once the new commit is staged.
		shell.RunCommand([]string{"git", "-C", "my_submodule_path", "commit", "--allow-empty", "-m", "submodule commit"})
		shell.CreateFile("my_submodule_path/dirty_file", "dirty content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Focus().
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			).
			// Staging the submodule stages the new commit, but the dirty
			// content remains unstaged, leaving us at "MM".
			PressPrimaryAction().
			Lines(
				Equals("MM my_submodule_path (submodule)").IsSelected(),
			).
			// Pressing again must unstage the submodule, taking us back to
			// " M" rather than trying (and failing) to stage the dirty content.
			PressPrimaryAction().
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			).
			// The same has to hold for the stage-all keybinding, which shares
			// the same decision logic: it stages the new commit...
			Press(keys.Files.ToggleStagedAll).
			Lines(
				Equals("MM my_submodule_path (submodule)").IsSelected(),
			).
			// ...and then unstages it again rather than getting stuck on the
			// dirty content.
			Press(keys.Files.ToggleStagedAll).
			Lines(
				Equals(" M my_submodule_path (submodule)").IsSelected(),
			)
	},
})
