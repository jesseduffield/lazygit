package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mattn/go-runewidth"
)

// This file is for the rendering of confirmation panels along with setting and handling associated
// keybindings.

func (gui *Gui) wrappedConfirmationFunction(cancel context.CancelFunc, function func() error) func() error {
	return func() error {
		cancel()

		if err := gui.c.PopContext(); err != nil {
			return err
		}

		if function != nil {
			if err := function(); err != nil {
				return gui.c.Error(err)
			}
		}

		return nil
	}
}

func (gui *Gui) wrappedPromptConfirmationFunction(cancel context.CancelFunc, function func(string) error, getResponse func() string) func() error {
	return func() error {
		cancel()

		if err := gui.c.PopContext(); err != nil {
			return err
		}

		if function != nil {
			if err := function(getResponse()); err != nil {
				return gui.c.Error(err)
			}
		}

		return nil
	}
}

func (gui *Gui) deactivateConfirmationPrompt() {
	gui.Mutexes.PopupMutex.Lock()
	gui.State.CurrentPopupOpts = nil
	gui.Mutexes.PopupMutex.Unlock()

	gui.Views.Confirmation.Visible = false
	gui.Views.Suggestions.Visible = false

	gui.clearConfirmationViewKeyBindings()
}

func (gui *Gui) getMessageHeight(wrap bool, message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	// if we need to wrap, calculate height to fit content within view's width
	if wrap {
		for _, line := range lines {
			lineCount += runewidth.StringWidth(line)/width + 1
		}
	} else {
		lineCount = len(lines)
	}
	return lineCount
}

func (gui *Gui) getPopupPanelDimensionsForContentHeight(panelWidth, contentHeight int) (int, int, int, int) {
	return gui.getPopupPanelDimensionsAux(panelWidth, contentHeight)
}

func (gui *Gui) getPopupPanelDimensionsAux(panelWidth int, panelHeight int) (int, int, int, int) {
	width, height := gui.g.Size()
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) getConfirmationPanelWidth() int {
	width, _ := gui.g.Size()
	// we want a minimum width up to a point, then we do it based on ratio.
	panelWidth := 4 * width / 7
	minWidth := 80
	if panelWidth < minWidth {
		if width-2 < minWidth {
			panelWidth = width - 2
		} else {
			panelWidth = minWidth
		}
	}

	return panelWidth
}

func (gui *Gui) prepareConfirmationPanel(
	ctx context.Context,
	opts types.ConfirmOpts,
) error {
	gui.Views.Confirmation.HasLoader = opts.HasLoader
	if opts.HasLoader {
		gui.g.StartTicking(ctx)
	}
	gui.Views.Confirmation.Title = opts.Title
	// for now we do not support wrapping in our editor
	gui.Views.Confirmation.Wrap = !opts.Editable
	gui.Views.Confirmation.FgColor = theme.GocuiDefaultTextColor
	gui.Views.Confirmation.Mask = runeForMask(opts.Mask)
	_ = gui.Views.Confirmation.SetOrigin(0, 0)

	gui.findSuggestions = opts.FindSuggestionsFunc
	if opts.FindSuggestionsFunc != nil {
		suggestionsView := gui.Views.Suggestions
		suggestionsView.Wrap = false
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		gui.setSuggestions(opts.FindSuggestionsFunc(""))
		suggestionsView.Visible = true
		suggestionsView.Title = fmt.Sprintf(gui.c.Tr.SuggestionsTitle, gui.c.UserConfig.Keybinding.Universal.TogglePanel)
	}

	gui.resizeConfirmationPanel()
	return nil
}

func runeForMask(mask bool) rune {
	if mask {
		return '*'
	}
	return 0
}

func (gui *Gui) createPopupPanel(ctx context.Context, opts types.CreatePopupPanelOpts) error {
	gui.Mutexes.PopupMutex.Lock()
	defer gui.Mutexes.PopupMutex.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	// we don't allow interruptions of non-loader popups in case we get stuck somehow
	// e.g. a credentials popup never gets its required user input so a process hangs
	// forever.
	// The proper solution is to have a queue of popup options
	if gui.State.CurrentPopupOpts != nil && !gui.State.CurrentPopupOpts.HasLoader {
		gui.Log.Error("ignoring create popup panel because a popup panel is already open")
		cancel()
		return nil
	}

	// remove any previous keybindings
	gui.clearConfirmationViewKeyBindings()

	err := gui.prepareConfirmationPanel(
		ctx,
		types.ConfirmOpts{
			Title:               opts.Title,
			Prompt:              opts.Prompt,
			HasLoader:           opts.HasLoader,
			FindSuggestionsFunc: opts.FindSuggestionsFunc,
			Editable:            opts.Editable,
			Mask:                opts.Mask,
		})
	if err != nil {
		cancel()
		return err
	}
	confirmationView := gui.Views.Confirmation
	confirmationView.Editable = opts.Editable
	confirmationView.Editor = gocui.EditorFunc(gui.promptEditor)

	if opts.Editable {
		textArea := confirmationView.TextArea
		textArea.Clear()
		textArea.TypeString(opts.Prompt)
		gui.resizeConfirmationPanel()
		confirmationView.RenderTextArea()
	} else {
		if err := gui.renderString(confirmationView, style.AttrBold.Sprint(opts.Prompt)); err != nil {
			cancel()
			return err
		}
	}

	if err := gui.setKeyBindings(cancel, opts); err != nil {
		cancel()
		return err
	}

	gui.State.CurrentPopupOpts = &opts

	return gui.c.PushContext(gui.State.Contexts.Confirmation)
}

func (gui *Gui) setKeyBindings(cancel context.CancelFunc, opts types.CreatePopupPanelOpts) error {
	actions := utils.ResolvePlaceholderString(
		gui.c.Tr.CloseConfirm,
		map[string]string{
			"keyBindClose":   "esc",
			"keyBindConfirm": "enter",
		},
	)

	_ = gui.renderString(gui.Views.Options, actions)
	var onConfirm func() error
	if opts.HandleConfirmPrompt != nil {
		onConfirm = gui.wrappedPromptConfirmationFunction(cancel, opts.HandleConfirmPrompt, func() string { return gui.Views.Confirmation.TextArea.GetContent() })
	} else {
		onConfirm = gui.wrappedConfirmationFunction(cancel, opts.HandleConfirm)
	}

	keybindingConfig := gui.c.UserConfig.Keybinding
	onSuggestionConfirm := gui.wrappedPromptConfirmationFunction(
		cancel,
		opts.HandleConfirmPrompt,
		gui.getSelectedSuggestionValue,
	)

	bindings := []*types.Binding{
		{
			ViewName: "confirmation",
			Key:      keybindings.GetKey(keybindingConfig.Universal.Confirm),
			Handler:  onConfirm,
		},
		{
			ViewName: "confirmation",
			Key:      keybindings.GetKey(keybindingConfig.Universal.Return),
			Handler:  gui.wrappedConfirmationFunction(cancel, opts.HandleClose),
		},
		{
			ViewName: "confirmation",
			Key:      keybindings.GetKey(keybindingConfig.Universal.TogglePanel),
			Handler: func() error {
				if len(gui.State.Suggestions) > 0 {
					return gui.replaceContext(gui.State.Contexts.Suggestions)
				}
				return nil
			},
		},
		{
			ViewName: "suggestions",
			Key:      keybindings.GetKey(keybindingConfig.Universal.Confirm),
			Handler:  onSuggestionConfirm,
		},
		{
			ViewName: "suggestions",
			Key:      keybindings.GetKey(keybindingConfig.Universal.Return),
			Handler:  gui.wrappedConfirmationFunction(cancel, opts.HandleClose),
		},
		{
			ViewName: "suggestions",
			Key:      keybindings.GetKey(keybindingConfig.Universal.TogglePanel),
			Handler:  func() error { return gui.replaceContext(gui.State.Contexts.Confirmation) },
		},
	}

	for _, binding := range bindings {
		if err := gui.SetKeybinding(binding); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) clearConfirmationViewKeyBindings() {
	keybindingConfig := gui.c.UserConfig.Keybinding
	_ = gui.g.DeleteKeybinding("confirmation", keybindings.GetKey(keybindingConfig.Universal.Confirm), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("confirmation", keybindings.GetKey(keybindingConfig.Universal.Return), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("suggestions", keybindings.GetKey(keybindingConfig.Universal.Confirm), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("suggestions", keybindings.GetKey(keybindingConfig.Universal.Return), gocui.ModNone)
}

func (gui *Gui) refreshSuggestions() {
	gui.suggestionsAsyncHandler.Do(func() func() {
		findSuggestionsFn := gui.findSuggestions
		if findSuggestionsFn != nil {
			suggestions := gui.findSuggestions(gui.c.GetPromptInput())
			return func() { gui.setSuggestions(suggestions) }
		} else {
			return func() {}
		}
	})
}

func (gui *Gui) handleAskFocused() error {
	keybindingConfig := gui.c.UserConfig.Keybinding

	message := utils.ResolvePlaceholderString(
		gui.c.Tr.CloseConfirm,
		map[string]string{
			"keyBindClose":   keybindings.Label(keybindingConfig.Universal.Return),
			"keyBindConfirm": keybindings.Label(keybindingConfig.Universal.Confirm),
		},
	)

	return gui.renderString(gui.Views.Options, message)
}
