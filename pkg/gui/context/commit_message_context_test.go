package context

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	guiTypes "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type commitMessageTestGuiCommon struct {
	guiTypes.IGuiCommon
	views guiTypes.Views
}

func (self commitMessageTestGuiCommon) Views() guiTypes.Views {
	return self.views
}

func TestCommitMessageGenerationSubtitle(t *testing.T) {
	userConfig := config.GetDefaultConfig()
	cmn := &common.Common{
		Log: logrus.NewEntry(logrus.New()),
		Tr:  i18n.EnglishTranslationSet(),
		Fs:  afero.NewMemMapFs(),
	}
	cmn.SetUserConfig(userConfig)

	commitMessageView := gocui.NewView("commitMessage", 0, 0, 80, 3, gocui.OutputNormal)
	commitDescriptionView := gocui.NewView("commitDescription", 0, 4, 80, 10, gocui.OutputNormal)
	commitMessageView.Editable = true
	commitDescriptionView.Editable = true

	ctx := NewCommitMessageContext(&ContextCommon{
		Common: cmn,
		IGuiCommon: commitMessageTestGuiCommon{views: guiTypes.Views{
			CommitMessage:     commitMessageView,
			CommitDescription: commitDescriptionView,
		}},
	})

	ctx.SetPanelState(NoCommitIndex, "Commit summary", "Commit description", false, "", nil, nil, false, "")
	assert.Contains(t, commitDescriptionView.Subtitle, "Press <tab> to toggle focus")

	cancelCalled := false
	ctx.StartGeneratingCommitMessage(func() {
		cancelCalled = true
	})

	assert.False(t, commitMessageView.Editable)
	assert.False(t, commitDescriptionView.Editable)
	assert.Contains(t, commitDescriptionView.Subtitle, "Generating commit message")
	assert.Contains(t, commitDescriptionView.Subtitle, "<esc>")

	assert.True(t, ctx.CancelGenerateCommitMessage())
	assert.True(t, cancelCalled)
	assert.Contains(t, commitDescriptionView.Subtitle, "Canceling commit message generation")

	ctx.StopGeneratingCommitMessage()
	assert.True(t, commitMessageView.Editable)
	assert.True(t, commitDescriptionView.Editable)
	assert.Contains(t, commitDescriptionView.Subtitle, "Press <tab> to toggle focus")
}
