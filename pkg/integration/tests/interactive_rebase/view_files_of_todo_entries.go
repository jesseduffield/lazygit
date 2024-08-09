package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ViewFilesOfTodoEntries = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Check that files of a pick todo can be viewed, but files of an update-ref todo can't",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetAppState().GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(1, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(1, 3)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("pick").Contains("CI commit 03").IsSelected(),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI <-- YOU ARE HERE --- commit 01"),
			).
			Press(keys.Universal.GoInto)

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file03.txt"),
			).
			PressEscape()

		t.Views().Commits().
			IsFocused().
			NavigateToLine(Contains("update-ref")).
			Press(keys.Universal.GoInto)

		t.ExpectToast(Equals("Disabled: Selected item does not have files to view"))
	},
})
