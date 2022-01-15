package gui

import (
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type PopupHandler interface {
	ErrorMsg(message string) error
	Error(err error) error
	Ask(opts askOpts) error
	Prompt(opts promptOpts) error
	WithLoaderPanel(message string, f func() error) error
	WithWaitingStatus(message string, f func() error) error
	Menu(opts createMenuOptions) error
}

type RealPopupHandler struct {
	*common.Common
	index int
	sync.Mutex
	createPopupPanelFn  func(createPopupPanelOpts) error
	onErrorFn           func() error
	closePopupFn        func() error
	createMenuFn        func(createMenuOptions) error
	withWaitingStatusFn func(message string, f func() error) error
}

var _ PopupHandler = &RealPopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(createPopupPanelOpts) error,
	onErrorFn func() error,
	closePopupFn func() error,
	createMenuFn func(createMenuOptions) error,
	withWaitingStatusFn func(message string, f func() error) error,
) *RealPopupHandler {
	return &RealPopupHandler{
		Common:              common,
		index:               0,
		createPopupPanelFn:  createPopupPanelFn,
		onErrorFn:           onErrorFn,
		closePopupFn:        closePopupFn,
		createMenuFn:        createMenuFn,
		withWaitingStatusFn: withWaitingStatusFn,
	}
}

func (self *RealPopupHandler) Menu(opts createMenuOptions) error {
	return self.createMenuFn(opts)
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

	coloredMessage := style.FgRed.Sprint(strings.TrimSpace(message))
	if err := self.onErrorFn(); err != nil {
		return err
	}

	return self.Ask(askOpts{
		title:  self.Tr.Error,
		prompt: coloredMessage,
	})
}

func (self *RealPopupHandler) Ask(opts askOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.prompt,
		handleConfirm:       opts.handleConfirm,
		handleClose:         opts.handleClose,
		handlersManageFocus: opts.handlersManageFocus,
	})
}

func (self *RealPopupHandler) Prompt(opts promptOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.initialContent,
		editable:            true,
		handleConfirmPrompt: opts.handleConfirm,
		findSuggestionsFunc: opts.findSuggestionsFunc,
	})
}

func (self *RealPopupHandler) WithLoaderPanel(message string, f func() error) error {
	index := 0
	self.Lock()
	self.index++
	index = self.index
	self.Unlock()

	err := self.createPopupPanelFn(createPopupPanelOpts{
		prompt:    message,
		hasLoader: true,
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

type TestPopupHandler struct {
	onError  func(message string) error
	onAsk    func(opts askOpts) error
	onPrompt func(opts promptOpts) error
}

func (self *TestPopupHandler) Error(message string) error {
	return self.onError(message)
}

func (self *TestPopupHandler) Ask(opts askOpts) error {
	return self.onAsk(opts)
}

func (self *TestPopupHandler) Prompt(opts promptOpts) error {
	return self.onPrompt(opts)
}

func (self *TestPopupHandler) WithLoaderPanel(message string, f func() error) error {
	return f()
}
