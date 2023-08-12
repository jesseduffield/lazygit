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
		shell.CreateFile("my-file.txt", "myfile content")
		shell.CreateFile("my-other-file.rb", "my-other-file content")

		shell.CreateNCommitsWithRandomMessages(60)
		shell.NewBranch("feature/demo")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("feature/demo", "origin/feature/demo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Interactive rebase")
		t.Wait(1000)

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("Add TypeScript types to User module")).
			Press(keys.Universal.Edit).
			SelectPreviousItem().
			Press(keys.Universal.Remove).
			SelectPreviousItem().
			Press(keys.Commits.SquashDown).
			SelectPreviousItem().
			Press(keys.Commits.MarkCommitAsFixup).
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
