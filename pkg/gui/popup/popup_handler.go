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

type IPopupHandler interface {
	ErrorMsg(message string) error
	Error(err error) error
	Ask(opts AskOpts) error
	Prompt(opts PromptOpts) error
	WithLoaderPanel(message string, f func() error) error
	WithWaitingStatus(message string, f func() error) error
	Menu(opts CreateMenuOptions) error
	Toast(message string)
	GetPromptInput() string
}

type CreateMenuOptions struct {
	Title      string
	Items      []*MenuItem
	HideCancel bool
}

type CreatePopupPanelOpts struct {
	HasLoader           bool
	Editable            bool
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleConfirmPrompt func(string) error
	HandleClose         func() error

	// when HandlersManageFocus is true, do not return from the confirmation context automatically. It's expected that the handlers will manage focus, whether that means switching to another context, or manually returning the context.
	HandlersManageFocus bool

	FindSuggestionsFunc func(string) []*types.Suggestion
}

type AskOpts struct {
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleClose         func() error
	HandlersManageFocus bool
}

type PromptOpts struct {
	Title               string
	InitialContent      string
	FindSuggestionsFunc func(string) []*types.Suggestion
	HandleConfirm       func(string) error
}

type MenuItem struct {
	DisplayString  string
	DisplayStrings []string
	OnPress        func() error
	// only applies when displayString is used
	OpensMenu bool
}

type RealPopupHandler struct {
	*common.Common
	index int
	sync.Mutex
	createPopupPanelFn  func(CreatePopupPanelOpts) error
	onErrorFn           func() error
	closePopupFn        func() error
	createMenuFn        func(CreateMenuOptions) error
	withWaitingStatusFn func(message string, f func() error) error
	toastFn             func(message string)
	getPromptInputFn    func() string
}

var _ IPopupHandler = &RealPopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(CreatePopupPanelOpts) error,
	onErrorFn func() error,
	closePopupFn func() error,
	createMenuFn func(CreateMenuOptions) error,
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

func (self *RealPopupHandler) Menu(opts CreateMenuOptions) error {
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

	coloredMessage := style.FgRed.Sprint(strings.TrimSpace(message))
	if err := self.onErrorFn(); err != nil {
		return err
	}

	return self.Ask(AskOpts{
		Title:  self.Tr.Error,
		Prompt: coloredMessage,
	})
}

func (self *RealPopupHandler) Ask(opts AskOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(CreatePopupPanelOpts{
		Title:               opts.Title,
		Prompt:              opts.Prompt,
		HandleConfirm:       opts.HandleConfirm,
		HandleClose:         opts.HandleClose,
		HandlersManageFocus: opts.HandlersManageFocus,
	})
}

func (self *RealPopupHandler) Prompt(opts PromptOpts) error {
	self.Lock()
	self.index++
	self.Unlock()

	return self.createPopupPanelFn(CreatePopupPanelOpts{
		Title:               opts.Title,
		Prompt:              opts.InitialContent,
		Editable:            true,
		HandleConfirmPrompt: opts.HandleConfirm,
		FindSuggestionsFunc: opts.FindSuggestionsFunc,
	})
}

func (self *RealPopupHandler) WithLoaderPanel(message string, f func() error) error {
	index := 0
	self.Lock()
	self.index++
	index = self.index
	self.Unlock()

	err := self.createPopupPanelFn(CreatePopupPanelOpts{
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
// asyncronously updating the suggestions list under the prompt.
func (self *RealPopupHandler) GetPromptInput() string {
	return self.getPromptInputFn()
}

type TestPopupHandler struct {
	OnErrorMsg func(message string) error
	OnAsk      func(opts AskOpts) error
	OnPrompt   func(opts PromptOpts) error
}

func (self *TestPopupHandler) Error(err error) error {
	return self.ErrorMsg(err.Error())
}

func (self *TestPopupHandler) ErrorMsg(message string) error {
	return self.OnErrorMsg(message)
}

func (self *TestPopupHandler) Ask(opts AskOpts) error {
	return self.OnAsk(opts)
}

func (self *TestPopupHandler) Prompt(opts PromptOpts) error {
	return self.OnPrompt(opts)
}

func (self *TestPopupHandler) WithLoaderPanel(message string, f func() error) error {
	return f()
}

func (self *TestPopupHandler) WithWaitingStatus(message string, f func() error) error {
	return f()
}

func (self *TestPopupHandler) Menu(opts CreateMenuOptions) error {
	panic("not yet implemented")
}

func (self *TestPopupHandler) Toast(message string) {
	panic("not yet implemented")
}

func (self *TestPopupHandler) CurrentInput() string {
	panic("not yet implemented")
}
