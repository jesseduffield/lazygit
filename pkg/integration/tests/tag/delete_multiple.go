package tag

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteMultiple = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete a range selection of tags remotely, locally and remotely, and locally",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CloneIntoRemote("origin")

		// The tags are sorted by date, newest first. Give each tag's commit a
		// date of its own so that the order is fixed.
		for i := 1; i <= 8; i++ {
			shell.EmptyCommitWithDate(fmt.Sprintf("commit %02d", i), fmt.Sprintf("2024-01-%02d 10:00:00", i))
			shell.CreateLightweightTag(fmt.Sprintf("tag-%02d", i), "HEAD")
		}

		for _, tag := range []string{"tag-03", "tag-04", "tag-05", "tag-06"} {
			shell.PushBranch("origin", "refs/tags/"+tag) // abusing PushBranch to push a tag
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			Lines(
				Contains("tag-08").IsSelected(),
				Contains("tag-07"),
				Contains("tag-06"),
				Contains("tag-05"),
				Contains("tag-04"),
				Contains("tag-03"),
				Contains("tag-02"),
				Contains("tag-01"),
			).
			NavigateToLine(Contains("tag-06")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected tags?")).
					Select(Contains("Delete remote tags")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Remote from which to remove the selected tags:")).
					InitialText(Equals("origin")).
					Confirm()

				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete selected tags?")).
					Content(Equals("Are you sure you want to delete the selected tags from 'origin'?")).
					Confirm()

				t.ExpectToast(Equals("Remote tags deleted"))
			}).
			// The local tags are still there, so the selection stays as it is
			Lines(
				Contains("tag-08"),
				Contains("tag-07"),
				Contains("tag-06").IsSelected(),
				Contains("tag-05").IsSelected(),
				Contains("tag-04"),
				Contains("tag-03"),
				Contains("tag-02"),
				Contains("tag-01"),
			).
			Tap(func() {
				t.Git().
					RemoteTagDeleted("origin", "tag-06").
					RemoteTagDeleted("origin", "tag-05")
			}).
			NavigateToLine(Contains("tag-04")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected tags?")).
					Select(Contains("Delete local and remote tags")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Equals("Remote from which to remove the selected tags:")).
					InitialText(Equals("origin")).
					Confirm()

				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete selected tags?")).
					Content(Equals("Are you sure you want to delete the selected tags from both your machine and from 'origin'?")).
					Confirm()
			}).
			Lines(
				Contains("tag-08"),
				Contains("tag-07"),
				Contains("tag-06"),
				Contains("tag-05"),
				Contains("tag-02").IsSelected(),
				Contains("tag-01"),
			).
			Tap(func() {
				t.Git().
					RemoteTagDeleted("origin", "tag-04").
					RemoteTagDeleted("origin", "tag-03")
			}).
			NavigateToLine(Contains("tag-07")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected tags?")).
					Select(Contains("Delete local tags")).
					Confirm()
			}).
			Lines(
				Contains("tag-08"),
				Contains("tag-05").IsSelected(),
				Contains("tag-02"),
				Contains("tag-01"),
			)
	},
})
