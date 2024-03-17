package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var InteractiveRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Interactive rebase",
	ExtraCmdArgs: []string{"log"},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateRepoHistory()

		shell.NewBranch("feature/demo")

		shell.CreateNCommitsWithRandomMessages(10)

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("feature/demo", "origin/feature/demo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Interactive rebase")
		t.Wait(1000)

		t.Views().Commits().
			IsFocused().
			Press(keys.Commits.StartInteractiveRebase).
			PressFast(keys.Universal.RangeSelectDown).
			PressFast(keys.Universal.RangeSelectDown).
			Press(keys.Commits.MarkCommitAsFixup).
			PressFast(keys.Commits.MoveDownCommit).
			PressFast(keys.Commits.MoveDownCommit).
			Delay().
			SelectNextItem().
			SelectNextItem().
			Press(keys.Universal.Remove).
			SelectNextItem().
			Press(keys.Commits.SquashDown).
			Press(keys.Universal.CreateRebaseOptionsMenu).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Contains("Rebase options")).
					Select(Contains("continue")).
					Confirm()
			}).
			SetCaptionPrefix("Push to remote").
			Press(keys.Universal.NextScreenMode).
			Press(keys.Universal.Push).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Contains("Force push")).
					Content(AnyString()).
					Confirm()
			})
	},
})
