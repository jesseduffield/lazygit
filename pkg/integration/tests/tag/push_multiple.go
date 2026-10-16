package tag

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PushMultiple = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Push a range selection of tags",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CloneIntoRemote("origin")

		// The tags are sorted by date, newest first. Give each tag's commit a
		// date of its own so that the order is fixed.
		for i := 1; i <= 4; i++ {
			shell.EmptyCommitWithDate(fmt.Sprintf("commit %02d", i), fmt.Sprintf("2024-01-%02d 10:00:00", i))
			shell.CreateLightweightTag(fmt.Sprintf("tag-%02d", i), "HEAD")
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("tag-04").IsSelected(),
				Contains("tag-03"),
				Contains("tag-02"),
				Contains("tag-01"),
			).
			NavigateToLine(Contains("tag-03")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Branches.PushTag).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Equals("Remote to push the selected tags to:")).
					InitialText(Equals("origin")).
					Confirm()
			}).
			Lines(
				Contains("tag-04"),
				Contains("tag-03").IsSelected(),
				Contains("tag-02").IsSelected(),
				Contains("tag-01"),
			)

		t.Git().RemoteTagNames("origin", []string{"tag-02", "tag-03"})
	},
})
