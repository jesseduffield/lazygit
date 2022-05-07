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

type PopupHandler struct {
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

var _ types.IPopupHandler = &PopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(types.CreatePopupPanelOpts) error,
	onErrorFn func() error,
	closePopupFn func() error,
	createMenuFn func(types.CreateMenuOptions) error,
	withWaitingStatusFn func(message string, f func() error) error,
	toastFn func(message string),
	getPromptInputFn func() string,
) *PopupHandler {
	return &PopupHandler{
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

func (self *PopupHandler) Menu(opts types.CreateMenuOptions) error {
	return self.createMenuFn(opts)
}

func (self *PopupHandler) Toast(message string) {
	self.toastFn(message)
}

func (self *PopupHandler) WithWaitingStatus(message string, f func() error) error {
	return self.withWaitingStatusFn(message, f)
}

func (self *PopupHandler) Error(err error) error {
	if err == gocui.ErrQuit {
		return err
	}

	return self.ErrorMsg(err.Error())
}

func (self *PopupHandler) ErrorMsg(message string) error {
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

func (self *PopupHandler) Alert(title string, message string) error {
	return self.Confirm(types.ConfirmOpts{Title: title, Prompt: message})
}

func (self *PopupHandler) Confirm(opts types.ConfirmOpts) error {
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

func (self *PopupHandler) Prompt(opts types.PromptOpts) error {
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

func (self *PopupHandler) WithLoaderPanel(message string, f func() error) error {
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
func (self *PopupHandler) GetPromptInput() string {
	return self.getPromptInputFn()
}
