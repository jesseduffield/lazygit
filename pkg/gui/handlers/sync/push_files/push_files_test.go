package push_files

import (
	"testing"

	commandsMocks "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files/push_filesfakes"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

var tr = i18n.EnglishTranslationSet()

func TestPushFilesHandler_Run(t *testing.T) {
	type fields struct {
		Gui Gui
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mockGui := &push_filesfakes.FakeGui{}
			mockGui.GetTrStub = func() *i18n.TranslationSet { return &tr }
			mockGui.PopupPanelFocusedReturns(false)
			mockGui.CurrentBranchReturns(&models.Branch{Pushables: "0", Pullables: "0", Name: "mybranch"})
			mockGui.WithPopupWaitingStatusStub = func(message string, f func() error) error {
				assert.Equal(t, "Pushing...", message)
				return f()
			}

			testGitCommandObj := &commandsMocks.FakeIGitCommand{}

			mockGui.GetGitCommandReturns(testGitCommandObj)

			testGitCommandObj.WithSpanReturns(testGitCommandObj)
			testGitCommandObj.PushStub = func(branchName string, force bool, upstream, args string, promptUserForCredential func(string) string) error {
				assert.Equal(t, "mybranch", branchName)
				assert.Equal(t, false, force)
				assert.Equal(t, "", upstream)
				return nil
			}

			handler := &PushFilesHandler{
				Gui: mockGui,
			}

			if err := handler.Run(); (err != nil) != tt.wantErr {
				t.Errorf("PushFilesHandler.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
