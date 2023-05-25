package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Restore part of a stash entry via applying a custom patch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("myfile", "content")
		shell.CreateFile("myfile2", "content")
		shell.GitAddAll()
		shell.Stash("stash one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash one").IsSelected(),
			).
			PressEnter().
			Tap(func() {
				t.Views().CommitFiles().
					IsFocused().
					Lines(
						Contains("myfile").IsSelected(),
						Contains("myfile2"),
					).
					PressPrimaryAction()

				t.Views().Information().Content(Contains("Building patch"))

				t.Views().
					CommitFiles().
					Press(keys.Universal.CreatePatchOptionsMenu)

				t.ExpectPopup().Menu().
					Title(Equals("Patch options")).
					Select(MatchesRegexp(`Apply patch$`)).Confirm()
			})

		t.Views().Files().Lines(
			Contains("myfile"),
		)
	},
})
