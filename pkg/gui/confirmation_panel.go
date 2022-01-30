package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) wrappedConfirmationFunction(handlersManageFocus bool, function func() error) func() error {
	return func() error {
		if err := gui.closeConfirmationPrompt(handlersManageFocus); err != nil {
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

func (gui *Gui) wrappedPromptConfirmationFunction(handlersManageFocus bool, function func(string) error, getResponse func() string) func() error {
	return func() error {
		if err := gui.closeConfirmationPrompt(handlersManageFocus); err != nil {
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

func (gui *Gui) closeConfirmationPrompt(handlersManageFocus bool) error {
	// we've already closed it so we can just return
	if !gui.Views.Confirmation.Visible {
		return nil
	}

	if !handlersManageFocus {
		if err := gui.returnFromContext(); err != nil {
			return err
		}
	}

	gui.clearConfirmationViewKeyBindings()
	gui.Views.Confirmation.Visible = false
	gui.Views.Suggestions.Visible = false

	return nil
}

func (gui *Gui) getMessageHeight(wrap bool, message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	// if we need to wrap, calculate height to fit content within view's width
	if wrap {
		for _, line := range lines {
			lineCount += len(line)/width + 1
		}
	} else {
		lineCount = len(lines)
	}
	return lineCount
}

func (gui *Gui) getConfirmationPanelDimensions(wrap bool, prompt string) (int, int, int, int) {
	width, height := gui.g.Size()
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
	panelHeight := gui.getMessageHeight(wrap, prompt, panelWidth)
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) prepareConfirmationPanel(
	title,
	prompt string,
	hasLoader bool,
	findSuggestionsFunc func(string) []*types.Suggestion,
	editable bool,
) error {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(true, prompt)
	// calling SetView on an existing view returns the same view, so I'm not bothering
	// to reassign to gui.Views.Confirmation
	_, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		return err
	}
	gui.Views.Confirmation.HasLoader = hasLoader
	if hasLoader {
		gui.g.StartTicking()
	}
	gui.Views.Confirmation.Title = title
	// for now we do not support wrapping in our editor
	gui.Views.Confirmation.Wrap = !editable
	gui.Views.Confirmation.FgColor = theme.GocuiDefaultTextColor

	gui.findSuggestions = findSuggestionsFunc
	if findSuggestionsFunc != nil {
		suggestionsViewHeight := 11
		suggestionsView, err := gui.g.SetView("suggestions", x0, y1+1, x1, y1+suggestionsViewHeight, 0)
		if err != nil {
			return err
		}
		suggestionsView.Wrap = false
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		gui.setSuggestions(findSuggestionsFunc(""))
		suggestionsView.Visible = true
		suggestionsView.Title = fmt.Sprintf(gui.c.Tr.SuggestionsTitle, gui.c.UserConfig.Keybinding.Universal.TogglePanel)
	}

	return nil
}

func (gui *Gui) createPopupPanel(opts types.CreatePopupPanelOpts) error {
	// remove any previous keybindings
	gui.clearConfirmationViewKeyBindings()

	err := gui.prepareConfirmationPanel(
		opts.Title,
		opts.Prompt,
		opts.HasLoader,
		opts.FindSuggestionsFunc,
		opts.Editable,
	)
	if err != nil {
		return err
	}
	confirmationView := gui.Views.Confirmation
	confirmationView.Editable = opts.Editable
	confirmationView.Editor = gocui.EditorFunc(gui.defaultEditor)

	if opts.Editable {
		textArea := confirmationView.TextArea
		textArea.Clear()
		textArea.TypeString(opts.Prompt)
		confirmationView.RenderTextArea()
	} else {
		if err := gui.renderString(confirmationView, opts.Prompt); err != nil {
			return err
		}
	}

	if err := gui.setKeyBindings(opts); err != nil {
		return err
	}

	return gui.c.PushContext(gui.State.Contexts.Confirmation)
}

func (gui *Gui) setKeyBindings(opts types.CreatePopupPanelOpts) error {
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
		onConfirm = gui.wrappedPromptConfirmationFunction(opts.HandlersManageFocus, opts.HandleConfirmPrompt, func() string { return gui.Views.Confirmation.TextArea.GetContent() })
	} else {
		onConfirm = gui.wrappedConfirmationFunction(opts.HandlersManageFocus, opts.HandleConfirm)
	}

	keybindingConfig := gui.c.UserConfig.Keybinding
	onSuggestionConfirm := gui.wrappedPromptConfirmationFunction(
		opts.HandlersManageFocus,
		opts.HandleConfirmPrompt,
		gui.getSelectedSuggestionValue,
	)

	bindings := []*types.Binding{
		{
			ViewName: "confirmation",
			Contexts: []string{string(context.CONFIRMATION_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.Confirm),
			Handler:  onConfirm,
		},
		{
			ViewName: "confirmation",
			Contexts: []string{string(context.CONFIRMATION_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.ConfirmAlt1),
			Handler:  onConfirm,
		},
		{
			ViewName: "confirmation",
			Contexts: []string{string(context.CONFIRMATION_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.Return),
			Handler:  gui.wrappedConfirmationFunction(opts.HandlersManageFocus, opts.HandleClose),
		},
		{
			ViewName: "confirmation",
			Contexts: []string{string(context.CONFIRMATION_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.TogglePanel),
			Handler: func() error {
				if len(gui.State.Suggestions) > 0 {
					return gui.replaceContext(gui.State.Contexts.Suggestions)
				}
				return nil
			},
		},
		{
			ViewName: "suggestions",
			Contexts: []string{string(context.SUGGESTIONS_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.Confirm),
			Handler:  onSuggestionConfirm,
		},
		{
			ViewName: "suggestions",
			Contexts: []string{string(context.SUGGESTIONS_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.ConfirmAlt1),
			Handler:  onSuggestionConfirm,
		},
		{
			ViewName: "suggestions",
			Contexts: []string{string(context.SUGGESTIONS_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.Return),
			Handler:  gui.wrappedConfirmationFunction(opts.HandlersManageFocus, opts.HandleClose),
		},
		{
			ViewName: "suggestions",
			Contexts: []string{string(context.SUGGESTIONS_CONTEXT_KEY)},
			Key:      gui.getKey(keybindingConfig.Universal.TogglePanel),
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
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.Confirm), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.ConfirmAlt1), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("confirmation", gui.getKey(keybindingConfig.Universal.Return), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("suggestions", gui.getKey(keybindingConfig.Universal.Confirm), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("suggestions", gui.getKey(keybindingConfig.Universal.ConfirmAlt1), gocui.ModNone)
	_ = gui.g.DeleteKeybinding("suggestions", gui.getKey(keybindingConfig.Universal.Return), gocui.ModNone)
}

func (gui *Gui) refreshSuggestions() {
	gui.suggestionsAsyncHandler.Do(func() func() {
		suggestions := gui.findSuggestions(gui.c.GetPromptInput())
		return func() { gui.setSuggestions(suggestions) }
	})
}
