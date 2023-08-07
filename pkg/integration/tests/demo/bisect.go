package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Bisect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Interactive rebase",
	ExtraCmdArgs: []string{"log"},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Gui.NerdFontsVersion = "3"
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
		t.SetCaptionPrefix("Git bisect")
		t.Wait(1000)

		markCommitAsBad := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as bad`)).Confirm()
		}

		markCommitAsGood := func() {
			t.Views().Commits().
				Press(keys.Commits.ViewBisectOptions)

			t.ExpectPopup().Menu().Title(Equals("Bisect")).Select(MatchesRegexp(`Mark .* as good`)).Confirm()
		}

		t.Views().Commits().
			IsFocused().
			Tap(func() {
				markCommitAsBad()

				t.Views().Information().Content(Contains("Bisecting"))
			}).
			SelectedLine(Contains("<-- bad")).
			NavigateToLine(Contains("Add TypeScript types to User module")).
			Tap(markCommitAsGood).
			SelectedLine(Contains("Add loading indicators to improve UX").Contains("<-- current")).
			Tap(markCommitAsBad).
			SelectedLine(Contains("Fix broken links on the help page").Contains("<-- current")).
			Tap(markCommitAsGood).
			SelectedLine(Contains("Add end-to-end tests for checkout flow").Contains("<-- current")).
			Tap(markCommitAsBad).
			Tap(func() {
				t.Wait(2000)

				t.ExpectPopup().Alert().Title(Equals("Bisect complete")).Content(MatchesRegexp("(?s).*Do you want to reset")).Confirm()
			}).
			SetCaptionPrefix("Inspect problematic commit").
			Wait(500).
			Press(keys.Universal.PrevScreenMode).
			IsFocused().
			Content(Contains("Add end-to-end tests for checkout flow")).
			Wait(500).
			PressEnter()

		t.Views().Information().Content(DoesNotContain("Bisecting"))
	},
})
