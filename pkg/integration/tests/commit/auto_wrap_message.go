package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AutoWrapMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit, and test how the commit message body is auto-wrapped",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Use a ridiculously small width so that we don't have to use so much test data
		config.GetUserConfig().Git.Commit.AutoWrapWidth = 20
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction(). // stage file
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("subject").
			SwitchToDescription().
			Type("Lorem ipsum dolor sit amet, consectetur adipiscing elit.").
			// See how it automatically inserted line feeds to wrap the text:
			Content(Equals("Lorem ipsum dolor \nsit amet, \nconsectetur \nadipiscing elit.")).
			SwitchToSummary().
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("subject"),
			).
			Focus().
			Tap(func() {
				t.Views().Main().Content(Contains(
					"subject\n    \n    Lorem ipsum dolor\n    sit amet,\n    consectetur\n    adipiscing elit."))
			}).
			Press(keys.Commits.RenameCommit)

		// Test that when rewording, the hard line breaks are turned back into
		// soft ones, so that we can insert text at the beginning and have the
		// paragraph reflow nicely.
		t.ExpectPopup().CommitMessagePanel().
			InitialText(Equals("subject")).
			SwitchToDescription().
			Content(Equals("Lorem ipsum dolor \nsit amet, \nconsectetur \nadipiscing elit.")).
			GoToBeginning().
			Type("More text. ").
			Content(Equals("More text. Lorem \nipsum dolor sit \namet, consectetur \nadipiscing elit."))
	},
})
