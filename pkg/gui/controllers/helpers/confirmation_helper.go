package helpers

import (
	goContext "context"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/mattn/go-runewidth"
)

type ConfirmationHelper struct {
	c        *types.HelperCommon
	contexts *context.ContextTree
}

func NewConfirmationHelper(
	c *types.HelperCommon,
	contexts *context.ContextTree,
) *ConfirmationHelper {
	return &ConfirmationHelper{
		c:        c,
		contexts: contexts,
	}
}

// This file is for the rendering of confirmation panels along with setting and handling associated
// keybindings.

func (self *ConfirmationHelper) wrappedConfirmationFunction(cancel goContext.CancelFunc, function func() error) func() error {
	return func() error {
		cancel()

		if err := self.c.PopContext(); err != nil {
			return err
		}

		if function != nil {
			if err := function(); err != nil {
				return self.c.Error(err)
			}
		}

		return nil
	}
}

func (self *ConfirmationHelper) wrappedPromptConfirmationFunction(cancel goContext.CancelFunc, function func(string) error, getResponse func() string) func() error {
	return self.wrappedConfirmationFunction(cancel, func() error {
		return function(getResponse())
	})
}

func (self *ConfirmationHelper) DeactivateConfirmationPrompt() {
	self.c.Mutexes().PopupMutex.Lock()
	self.c.State().GetRepoState().SetCurrentPopupOpts(nil)
	self.c.Mutexes().PopupMutex.Unlock()

	self.c.Views().Confirmation.Visible = false
	self.c.Views().Suggestions.Visible = false

	self.clearConfirmationViewKeyBindings()
}

func getMessageHeight(wrap bool, message string, width int) int {
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

func (self *ConfirmationHelper) getConfirmationPanelDimensions(wrap bool, prompt string) (int, int, int, int) {
	panelWidth := self.getConfirmationPanelWidth()
	panelHeight := getMessageHeight(wrap, prompt, panelWidth)
	return self.getConfirmationPanelDimensionsAux(panelWidth, panelHeight)
}

func (self *ConfirmationHelper) getConfirmationPanelDimensionsForContentHeight(panelWidth, contentHeight int) (int, int, int, int) {
	return self.getConfirmationPanelDimensionsAux(panelWidth, contentHeight)
}

func (self *ConfirmationHelper) getConfirmationPanelDimensionsAux(panelWidth int, panelHeight int) (int, int, int, int) {
	width, height := self.c.GocuiGui().Size()
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (self *ConfirmationHelper) getConfirmationPanelWidth() int {
	width, _ := self.c.GocuiGui().Size()
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

func (self *ConfirmationHelper) prepareConfirmationPanel(
	ctx goContext.Context,
	opts types.ConfirmOpts,
) error {
	self.c.Views().Confirmation.HasLoader = opts.HasLoader
	if opts.HasLoader {
		self.c.GocuiGui().StartTicking(ctx)
	}
	self.c.Views().Confirmation.Title = opts.Title
	// for now we do not support wrapping in our editor
	self.c.Views().Confirmation.Wrap = !opts.Editable
	self.c.Views().Confirmation.FgColor = theme.GocuiDefaultTextColor
	self.c.Views().Confirmation.Mask = runeForMask(opts.Mask)
	_ = self.c.Views().Confirmation.SetOrigin(0, 0)

	suggestionsContext := self.contexts.Suggestions
	suggestionsContext.State.FindSuggestions = opts.FindSuggestionsFunc
	if opts.FindSuggestionsFunc != nil {
		suggestionsView := self.c.Views().Suggestions
		suggestionsView.Wrap = false
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		suggestionsContext.SetSuggestions(opts.FindSuggestionsFunc(""))
		suggestionsView.Visible = true
		suggestionsView.Title = fmt.Sprintf(self.c.Tr.SuggestionsTitle, self.c.UserConfig.Keybinding.Universal.TogglePanel)
	}

	self.ResizeConfirmationPanel()
	return nil
}

func runeForMask(mask bool) rune {
	if mask {
		return '*'
	}
	return 0
}

func (self *ConfirmationHelper) CreatePopupPanel(ctx goContext.Context, opts types.CreatePopupPanelOpts) error {
	self.c.Mutexes().PopupMutex.Lock()
	defer self.c.Mutexes().PopupMutex.Unlock()

	ctx, cancel := goContext.WithCancel(ctx)

	// we don't allow interruptions of non-loader popups in case we get stuck somehow
	// e.g. a credentials popup never gets its required user input so a process hangs
	// forever.
	// The proper solution is to have a queue of popup options
	currentPopupOpts := self.c.State().GetRepoState().GetCurrentPopupOpts()
	if currentPopupOpts != nil && !currentPopupOpts.HasLoader {
		self.c.Log.Error("ignoring create popup panel because a popup panel is already open")
		cancel()
		return nil
	}

	// remove any previous keybindings
	self.clearConfirmationViewKeyBindings()

	err := self.prepareConfirmationPanel(
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
	confirmationView := self.c.Views().Confirmation
	confirmationView.Editable = opts.Editable

	if opts.Editable {
		textArea := confirmationView.TextArea
		textArea.Clear()
		textArea.TypeString(opts.Prompt)
		self.ResizeConfirmationPanel()
		confirmationView.RenderTextArea()
	} else {
		self.c.ResetViewOrigin(confirmationView)
		self.c.SetViewContent(confirmationView, style.AttrBold.Sprint(opts.Prompt))
	}

	if err := self.setKeyBindings(cancel, opts); err != nil {
		cancel()
		return err
	}

	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	return self.c.PushContext(self.contexts.Confirmation)
}

func (self *ConfirmationHelper) setKeyBindings(cancel goContext.CancelFunc, opts types.CreatePopupPanelOpts) error {
	var onConfirm func() error
	if opts.HandleConfirmPrompt != nil {
		onConfirm = self.wrappedPromptConfirmationFunction(cancel, opts.HandleConfirmPrompt, func() string { return self.c.Views().Confirmation.TextArea.GetContent() })
	} else {
		onConfirm = self.wrappedConfirmationFunction(cancel, opts.HandleConfirm)
	}

	onSuggestionConfirm := self.wrappedPromptConfirmationFunction(
		cancel,
		opts.HandleConfirmPrompt,
		self.getSelectedSuggestionValue,
	)

	onClose := self.wrappedConfirmationFunction(cancel, opts.HandleClose)

	self.contexts.Confirmation.State.OnConfirm = onConfirm
	self.contexts.Confirmation.State.OnClose = onClose
	self.contexts.Suggestions.State.OnConfirm = onSuggestionConfirm
	self.contexts.Suggestions.State.OnClose = onClose

	return nil
}

func (self *ConfirmationHelper) clearConfirmationViewKeyBindings() {
	noop := func() error { return nil }
	self.contexts.Confirmation.State.OnConfirm = noop
	self.contexts.Confirmation.State.OnClose = noop
	self.contexts.Suggestions.State.OnConfirm = noop
	self.contexts.Suggestions.State.OnClose = noop
}

func (self *ConfirmationHelper) getSelectedSuggestionValue() string {
	selectedSuggestion := self.contexts.Suggestions.GetSelected()

	if selectedSuggestion != nil {
		return selectedSuggestion.Value
	}

	return ""
}

func (self *ConfirmationHelper) ResizeConfirmationPanel() {
	suggestionsViewHeight := 0
	if self.c.Views().Suggestions.Visible {
		suggestionsViewHeight = 11
	}
	panelWidth := self.getConfirmationPanelWidth()
	prompt := self.c.Views().Confirmation.Buffer()
	wrap := true
	if self.c.Views().Confirmation.Editable {
		prompt = self.c.Views().Confirmation.TextArea.GetContent()
		wrap = false
	}
	panelHeight := getMessageHeight(wrap, prompt, panelWidth) + suggestionsViewHeight
	x0, y0, x1, y1 := self.getConfirmationPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (self *ConfirmationHelper) ResizeCurrentPopupPanel() error {
	v := self.c.GocuiGui().CurrentView()
	if v == nil {
		return nil
	}

	if v == self.c.Views().Menu {
		self.resizeMenu()
	} else if v == self.c.Views().Confirmation || v == self.c.Views().Suggestions {
		self.ResizeConfirmationPanel()
	} else if self.IsPopupPanel(v.Name()) {
		return self.ResizePopupPanel(v, v.Buffer())
	}

	return nil
}

func (self *ConfirmationHelper) ResizePopupPanel(v *gocui.View, content string) error {
	x0, y0, x1, y1 := self.getConfirmationPanelDimensions(v.Wrap, content)
	_, err := self.c.GocuiGui().SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

func (self *ConfirmationHelper) resizeMenu() {
	itemCount := self.contexts.Menu.GetList().Len()
	offset := 3
	panelWidth := self.getConfirmationPanelWidth()
	x0, y0, x1, y1 := self.getConfirmationPanelDimensionsForContentHeight(panelWidth, itemCount+offset)
	menuBottom := y1 - offset
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Menu.Name(), x0, y0, x1, menuBottom, 0)

	tooltipTop := menuBottom + 1
	tooltipHeight := getMessageHeight(true, self.contexts.Menu.GetSelected().Tooltip, panelWidth) + 2 // plus 2 for the frame
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Tooltip.Name(), x0, tooltipTop, x1, tooltipTop+tooltipHeight-1, 0)
}

func (self *ConfirmationHelper) IsPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "confirmation" || viewName == "menu"
}

func (self *ConfirmationHelper) IsPopupPanelFocused() bool {
	return self.IsPopupPanel(self.c.CurrentContext().GetViewName())
}
