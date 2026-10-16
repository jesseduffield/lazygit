package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageNestedRepository = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a nested repository, or all files while there is one, is not supported, because git would add it as an embedded repository",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")
		shell.RunCommand([]string{"git", "rm", "--cached", "modules/sub"})
		shell.RunCommand([]string{"git", "rm", ".gitmodules"})
		shell.Commit("remove submodule")

		shell.CreateFile("file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? file").IsSelected(),
				Equals("?? modules/sub (removed submodule)"),
			).
			Press(keys.Files.ToggleStagedAll).
			Tap(func() {
				t.ExpectToast(Equals("Disabled: Can't stage 'modules/sub' because it is a git repository of its own. Discard it with d, or ignore it with i."))
			}).
			NavigateToLine(Contains("modules/sub")).
			PressPrimaryAction().
			Tap(func() {
				t.ExpectToast(Equals("Disabled: Can't stage 'modules/sub' because it is a git repository of its own. Discard it with d, or ignore it with i."))
			}).
			// Staging all files works if the nested repository is filtered out
			FilterOrSearch("file").
			Lines(
				Equals("?? file").IsSelected(),
			).
			Press(keys.Files.ToggleStagedAll).
			Lines(
				Equals("A  file").IsSelected(),
			)
	},
})
