package popup

import (
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type RealPopupHandler struct {
	*common.Common
	index int
	sync.Mutex
	createPopupPanelFn  func(types.CreatePopupPanelOpts) error
	onErrorFn           func() error
	closePopupFn        func() error
	createMenuFn        func(types.CreateMenuOptions) error
	withWaitingStatusFn func(message string, f func() error) error
	toastFn             func(message string)
	getPromptInputFn    func() string
}

var _ types.IPopupHandler = &RealPopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(types.CreatePopupPanelOpts) error,
	onErrorFn func() error,
	closePopupFn func() error,
	createMenuFn func(types.CreateMenuOptions) error,
	withWaitingStatusFn func(message string, f func() error) error,
	toastFn func(message string),
	getPromptInputFn func() string,
) *RealPopupHandler {
	return &RealPopupHandler{
		Common:              common,
		index:               0,
		createPopupPanelFn:  createPopupPanelFn,
		onErrorFn:           onErrorFn,
		closePopupFn:        closePopupFn,
		createMenuFn:        createMenuFn,
		withWaitingStatusFn: withWaitingStatusFn,
		toastFn:             toastFn,
		getPromptInputFn:    getPromptInputFn,
	}
}

func (self *RealPopupHandler) Menu(opts types.CreateMenuOptions) error {
	return self.createMenuFn(opts)
}

func (self *RealPopupHandler) Toast(message string) {
	self.toastFn(message)
}

func (self *RealPopupHandler) WithWaitingStatus(message string, f func() error) error {
	return self.withWaitingStatusFn(message, f)
}

func (self *RealPopupHandler) Error(err error) error {
	if err == gocui.ErrQuit {
		return err
	}

	return self.ErrorMsg(err.Error())
}

func (self *RealPopupHandler) ErrorMsg(message string) error {
	self.Lock()
	self.index++
	self.Unlock()

	// Need to set bold here explicitly; otherwise it gets cancelled by the red colouring.
	coloredMessage := style.FgRed.SetBold().Sprint(strings.TrimSpace(message))
	if err := self.onErrorFn(); err != nil {
		return err
	}

	return self.Alert(self.Tr.Error, coloredMessage)
}

func (self *RealPopupHandler) Alert(title string, message string) error {
	return self.Confirm(types.ConfirmOpts{Title: title, Prompt: message})
}

func (self *RealPopupHandler) Confirm(opts types.ConfirmOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(types.CreatePopupPanelOpts{
		Title:               opts.Title,
		Prompt:              opts.Prompt,
		HandleConfirm:       opts.HandleConfirm,
		HandleClose:         opts.HandleClose,
		HandlersManageFocus: opts.HandlersManageFocus,
	})
}

func (self *RealPopupHandler) Prompt(opts types.PromptOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(types.CreatePopupPanelOpts{
		Title:               opts.Title,
		Prompt:              opts.InitialContent,
		Editable:            true,
		HandleConfirmPrompt: opts.HandleConfirm,
		HandleClose:         opts.HandleClose,
		FindSuggestionsFunc: opts.FindSuggestionsFunc,
		Mask:                opts.Mask,
	})
}

func (self *RealPopupHandler) WithLoaderPanel(message string, f func() error) error {
	index := 0
	self.Lock()
	self.index++
	index = self.index
	self.Unlock()

	err := self.createPopupPanelFn(types.CreatePopupPanelOpts{
		Prompt:    message,
		HasLoader: true,
	})
	if err != nil {
		self.Log.Error(err)
		return nil
	}

	go utils.Safe(func() {
		if err := f(); err != nil {
			self.Log.Error(err)
		}

		self.Lock()
		if index == self.index {
			_ = self.closePopupFn()
		}
		self.Unlock()
	})

	return nil
}

// returns the content that has currently been typed into the prompt. Useful for
// asynchronously updating the suggestions list under the prompt.
func (self *RealPopupHandler) GetPromptInput() string {
	return self.getPromptInputFn()
}

type TestPopupHandler struct {
	OnErrorMsg func(message string) error
	OnConfirm  func(opts types.ConfirmOpts) error
	OnPrompt   func(opts types.PromptOpts) error
}

var _ types.IPopupHandler = &TestPopupHandler{}

func (self *TestPopupHandler) Error(err error) error {
	return self.ErrorMsg(err.Error())
}

func (self *TestPopupHandler) ErrorMsg(message string) error {
	return self.OnErrorMsg(message)
}

func (self *TestPopupHandler) Alert(title string, message string) error {
	panic("not yet implemented")
}

func (self *TestPopupHandler) Confirm(opts types.ConfirmOpts) error {
	return self.Confirm(opts)
}

func (self *TestPopupHandler) Prompt(opts types.PromptOpts) error {
	return self.OnPrompt(opts)
}

func (self *TestPopupHandler) WithLoaderPanel(message string, f func() error) error {
	return f()
}

func (self *TestPopupHandler) WithWaitingStatus(message string, f func() error) error {
	return f()
}

func (self *TestPopupHandler) Menu(opts types.CreateMenuOptions) error {
	panic("not yet implemented")
}

func (self *TestPopupHandler) Toast(message string) {
	panic("not yet implemented")
}

func (self *TestPopupHandler) GetPromptInput() string {
	panic("not yet implemented")
}
