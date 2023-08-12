package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendOldCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amend old commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
		config.UserConfig.Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(60)
		shell.NewBranch("feature/demo")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("feature/demo", "origin/feature/demo")

		shell.UpdateFile("navigation/site_navigation.go", "package navigation\n\nfunc Navigate() {\n\tpanic(\"unimplemented\")\n}")
		shell.CreateFile("docs/README.md", "my readme content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Amend an old commit")
		t.Wait(1000)

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("site_navigation.go")).
			PressPrimaryAction()

		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("Improve accessibility of site navigation")).
			Wait(500).
			Press(keys.Commits.AmendToCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Amend commit")).
					Wait(1000).
					Content(AnyString()).
					Confirm()

				t.Wait(1000)
			}).
			Press(keys.Universal.Push).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Force push")).
					Content(AnyString()).
					Wait(1000).
					Confirm()
			})
	},
})
