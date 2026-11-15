package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyCommitMessageFromPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the commit message from the commit message panel to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > ../clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("subject").
			SwitchToDescription().
			Type("This line is longer than the auto-wrap width, so the clipboard gets it with a line break").
			AddNewline().
			Type("second line").
			SwitchToSummary().
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Copy commit message to clipboard")).
			Confirm()

		t.ExpectToast(Equals("Commit message copied to clipboard"))

		t.FileSystem().FileContent("../clipboard", Equals(
			"subject\n\nThis line is longer than the auto-wrap width, so the clipboard gets it\nwith a line break\nsecond line"))

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("subject"))
	},
})
