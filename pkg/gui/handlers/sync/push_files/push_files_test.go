package push_files

import (
	"testing"

	commandsMocks "github.com/jesseduffield/lazygit/pkg/commands/mocks"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files/mocks"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/mock"
)

type MockGui struct {
	*mocks.Gui
}

func (m *MockGui) GetTr() *i18n.TranslationSet {
	tr := i18n.EnglishTranslationSet()
	return &tr
}

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
			testObj := &MockGui{Gui: new(mocks.Gui)}
			testObj.On("PopupPanelFocused").Return(false)
			testObj.On("CurrentBranch").Return(&models.Branch{Pushables: "0", Pullables: "0", Name: "mybranch"})
			testObj.On("WithPopupWaitingStatus", "Pushing...", mock.AnythingOfType("func() error")).Return(nil).Run(func(args mock.Arguments) {
				callback := args.Get(1).(func() error)
				_ = callback()
			})

			testGitCommandObj := new(commandsMocks.IGitCommand)

			testObj.On("GetGitCommand").Return(testGitCommandObj)

			testGitCommandObj.On("WithSpan", "Push").Return(testGitCommandObj)
			testGitCommandObj.On("Push", "mybranch", false, "", "", mock.AnythingOfType("func(string) string")).Return(nil)

			testObj.On("HandleCredentialsPopup", nil).Return()
			testObj.On("RefreshSidePanels", RefreshOptions{Mode: ASYNC}).Return(nil)

			handler := &PushFilesHandler{
				Gui: testObj,
			}

			if err := handler.Run(); (err != nil) != tt.wantErr {
				t.Errorf("PushFilesHandler.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
