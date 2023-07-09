package popup

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FakePopupHandler struct {
	OnErrorMsg func(message string) error
	OnConfirm  func(opts types.ConfirmOpts) error
	OnPrompt   func(opts types.PromptOpts) error
}

var _ types.IPopupHandler = &FakePopupHandler{}

func (self *FakePopupHandler) Error(err error) error {
	return self.ErrorMsg(err.Error())
}

func (self *FakePopupHandler) ErrorMsg(message string) error {
	return self.OnErrorMsg(message)
}

func (self *FakePopupHandler) Alert(title string, message string) error {
	panic("not yet implemented")
}

func (self *FakePopupHandler) Confirm(opts types.ConfirmOpts) error {
	return self.OnConfirm(opts)
}

func (self *FakePopupHandler) Prompt(opts types.PromptOpts) error {
	return self.OnPrompt(opts)
}

func (self *FakePopupHandler) WithLoaderPanel(message string, f func(gocui.Task) error) error {
	return f(gocui.NewFakeTask())
}

func (self *FakePopupHandler) WithWaitingStatus(message string, f func(gocui.Task) error) error {
	return f(gocui.NewFakeTask())
}

func (self *FakePopupHandler) Menu(opts types.CreateMenuOptions) error {
	panic("not yet implemented")
}

func (self *FakePopupHandler) Toast(message string) {
	panic("not yet implemented")
}

func (self *FakePopupHandler) GetPromptInput() string {
	panic("not yet implemented")
}
