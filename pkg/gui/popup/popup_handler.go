package popup

import (
	"context"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/common"
	gctx "github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
)

type PopupHandler struct {
	*common.Common
	index int
	deadlock.Mutex
	createPopupPanelFn  func(context.Context, types.CreatePopupPanelOpts) error
	onErrorFn           func() error
	popContextFn        func() error
	currentContextFn    func() types.Context
	createMenuFn        func(types.CreateMenuOptions) error
	withWaitingStatusFn func(message string, f func() error) error
	toastFn             func(message string)
	getPromptInputFn    func() string
}

var _ types.IPopupHandler = &PopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(context.Context, types.CreatePopupPanelOpts) error,
	onErrorFn func() error,
	popContextFn func() error,
	currentContextFn func() types.Context,
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
		popContextFn:        popContextFn,
		currentContextFn:    currentContextFn,
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

	return self.createPopupPanelFn(context.Background(), types.CreatePopupPanelOpts{
		Title:         opts.Title,
		Prompt:        opts.Prompt,
		HandleConfirm: opts.HandleConfirm,
		HandleClose:   opts.HandleClose,
	})
}

func (self *PopupHandler) Prompt(opts types.PromptOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(context.Background(), types.CreatePopupPanelOpts{
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

	ctx, cancel := context.WithCancel(context.Background())

	err := self.createPopupPanelFn(ctx, types.CreatePopupPanelOpts{
		Prompt:    message,
		HasLoader: true,
	})
	if err != nil {
		self.Log.Error(err)
		cancel()
		return nil
	}

	go utils.Safe(func() {
		if err := f(); err != nil {
			self.Log.Error(err)
		}

		cancel()

		self.Lock()
		if index == self.index && self.currentContextFn().GetKey() == gctx.CONFIRMATION_CONTEXT_KEY {
			_ = self.popContextFn()
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
