package gui

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
)

type PopupHandler interface {
	Error(message string) error
	Ask(opts askOpts) error
	Prompt(opts promptOpts) error
	Loader(message string) error
}

type RealPopupHandler struct {
	gui *Gui
}

func (self *RealPopupHandler) Error(message string) error {
	gui := self.gui

	coloredMessage := style.FgRed.Sprint(strings.TrimSpace(message))
	if err := gui.refreshSidePanels(refreshOptions{mode: ASYNC}); err != nil {
		return err
	}

	return self.Ask(askOpts{
		title:  gui.Tr.Error,
		prompt: coloredMessage,
	})
}

func (self *RealPopupHandler) Ask(opts askOpts) error {
	gui := self.gui

	return gui.createPopupPanel(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.prompt,
		handleConfirm:       opts.handleConfirm,
		handleClose:         opts.handleClose,
		handlersManageFocus: opts.handlersManageFocus,
	})
}

func (self *RealPopupHandler) Prompt(opts promptOpts) error {
	gui := self.gui

	return gui.createPopupPanel(createPopupPanelOpts{
		title:               opts.title,
		prompt:              opts.initialContent,
		editable:            true,
		handleConfirmPrompt: opts.handleConfirm,
		findSuggestionsFunc: opts.findSuggestionsFunc,
	})
}

func (self *RealPopupHandler) Loader(message string) error {
	gui := self.gui

	return gui.createPopupPanel(createPopupPanelOpts{
		prompt:    message,
		hasLoader: true,
	})
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

func (self *TestPopupHandler) Loader(message string) error {
	return nil
}
