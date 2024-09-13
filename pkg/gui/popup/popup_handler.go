package popup

import (
	"context"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PopupHandler struct {
	*common.Common
	createPopupPanelFn      func(context.Context, types.CreatePopupPanelOpts)
	onErrorFn               func() error
	popContextFn            func()
	currentContextFn        func() types.Context
	createMenuFn            func(types.CreateMenuOptions) error
	withWaitingStatusFn     func(message string, f func(gocui.Task) error)
	withWaitingStatusSyncFn func(message string, f func() error) error
	toastFn                 func(message string, kind types.ToastKind)
	getPromptInputFn        func() string
	inDemo                  func() bool
}

var _ types.IPopupHandler = &PopupHandler{}

func NewPopupHandler(
	common *common.Common,
	createPopupPanelFn func(context.Context, types.CreatePopupPanelOpts),
	onErrorFn func() error,
	popContextFn func(),
	currentContextFn func() types.Context,
	createMenuFn func(types.CreateMenuOptions) error,
	withWaitingStatusFn func(message string, f func(gocui.Task) error),
	withWaitingStatusSyncFn func(message string, f func() error) error,
	toastFn func(message string, kind types.ToastKind),
	getPromptInputFn func() string,
	inDemo func() bool,
) *PopupHandler {
	return &PopupHandler{
		Common:                  common,
		createPopupPanelFn:      createPopupPanelFn,
		onErrorFn:               onErrorFn,
		popContextFn:            popContextFn,
		currentContextFn:        currentContextFn,
		createMenuFn:            createMenuFn,
		withWaitingStatusFn:     withWaitingStatusFn,
		withWaitingStatusSyncFn: withWaitingStatusSyncFn,
		toastFn:                 toastFn,
		getPromptInputFn:        getPromptInputFn,
		inDemo:                  inDemo,
	}
}

func (self *PopupHandler) Menu(opts types.CreateMenuOptions) error {
	return self.createMenuFn(opts)
}

func (self *PopupHandler) Toast(message string) {
	self.toastFn(message, types.ToastKindStatus)
}

func (self *PopupHandler) ErrorToast(message string) {
	self.toastFn(message, types.ToastKindError)
}

func (self *PopupHandler) SetToastFunc(f func(string, types.ToastKind)) {
	self.toastFn = f
}

func (self *PopupHandler) WithWaitingStatus(message string, f func(gocui.Task) error) error {
	self.withWaitingStatusFn(message, f)
	return nil
}

func (self *PopupHandler) WithWaitingStatusSync(message string, f func() error) error {
	return self.withWaitingStatusSyncFn(message, f)
}

func (self *PopupHandler) ErrorHandler(err error) error {
	// Need to set bold here explicitly; otherwise it gets cancelled by the red colouring.
	coloredMessage := style.FgRed.SetBold().Sprint(strings.TrimSpace(err.Error()))
	if err := self.onErrorFn(); err != nil {
		return err
	}

	self.Alert(self.Tr.Error, coloredMessage)

	return nil
}

func (self *PopupHandler) Alert(title string, message string) {
	self.Confirm(types.ConfirmOpts{Title: title, Prompt: message})
}

func (self *PopupHandler) Confirm(opts types.ConfirmOpts) {
	self.createPopupPanelFn(context.Background(), types.CreatePopupPanelOpts{
		Title:         opts.Title,
		Prompt:        opts.Prompt,
		HandleConfirm: opts.HandleConfirm,
		HandleClose:   opts.HandleClose,
	})
}

func (self *PopupHandler) Prompt(opts types.PromptOpts) {
	self.createPopupPanelFn(context.Background(), types.CreatePopupPanelOpts{
		Title:                  opts.Title,
		Prompt:                 opts.InitialContent,
		Editable:               true,
		HandleConfirmPrompt:    opts.HandleConfirm,
		HandleClose:            opts.HandleClose,
		HandleDeleteSuggestion: opts.HandleDeleteSuggestion,
		FindSuggestionsFunc:    opts.FindSuggestionsFunc,
		AllowEditSuggestion:    opts.AllowEditSuggestion,
		Mask:                   opts.Mask,
	})
}

// returns the content that has currently been typed into the prompt. Useful for
// asynchronously updating the suggestions list under the prompt.
func (self *PopupHandler) GetPromptInput() string {
	return self.getPromptInputFn()
}
