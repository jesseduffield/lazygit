package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MoveToIndexWhenCommitNotInCurrentBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch into the index after checking out a branch that doesn't contain the patch's commit; expect a friendly error rather than a crash",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("base commit")
		shell.NewBranch("feature")
		shell.CreateFileAndAdd("file1", "file1 content\n")
		shell.Commit("feature commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Build a custom patch from a commit that only exists on the feature branch
		t.Views().Commits().
			Focus().
			Lines(
				Contains("feature commit").IsSelected(),
				Contains("base commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		// Check out a branch that does not contain the patch's commit, so the
		// commit the patch was built from is no longer in the commits list.
		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("master")).
			PressPrimaryAction()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base commit").IsSelected(),
			)

		// Previously this panicked with "index out of range [-1]" because the
		// patch's commit could not be found in the current commits list. It
		// should now report a friendly error instead.
		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("Cannot find the commit this custom patch was created from")).
			Confirm()
	},
})
