package helpers

import (
	goContext "context"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type ConfirmationHelper struct {
	c *HelperCommon
}

func NewConfirmationHelper(c *HelperCommon) *ConfirmationHelper {
	return &ConfirmationHelper{
		c: c,
	}
}

// This file is for the rendering of confirmation panels along with setting and handling associated
// keybindings.

func (self *ConfirmationHelper) closeAndCallConfirmationFunction(cancel goContext.CancelFunc, function func() error) error {
	cancel()

	self.c.Context().Pop()

	if function != nil {
		if err := function(); err != nil {
			return err
		}
	}

	return nil
}

func (self *ConfirmationHelper) wrappedConfirmationFunction(cancel goContext.CancelFunc, function func() error) func() error {
	return func() error {
		return self.closeAndCallConfirmationFunction(cancel, function)
	}
}

func (self *ConfirmationHelper) wrappedPromptConfirmationFunction(
	cancel goContext.CancelFunc,
	function func(string) error,
	getResponse func() string,
	allowEmptyInput bool,
	preserveWhitespace bool,
) func() error {
	return func() error {
		if self.c.GocuiGui().IsPasting {
			// The user is pasting multi-line text into a prompt; we don't want to handle the
			// line feeds as "confirm" keybindings. Simply ignoring them is the best we can do; this
			// will cause the entire pasted text to appear as a single line in the prompt. Hopefully
			// the user knows that ctrl-u allows them to delete it again...
			return nil
		}

		response := getResponse()
		if !preserveWhitespace {
			response = strings.TrimSpace(response)
		}

		if response == "" && !allowEmptyInput {
			self.c.ErrorToast(self.c.Tr.PromptInputCannotBeEmptyToast)
			return nil
		}

		return self.closeAndCallConfirmationFunction(cancel, func() error {
			return function(response)
		})
	}
}

func (self *ConfirmationHelper) DeactivateConfirmation() {
	self.c.Mutexes().PopupMutex.Lock()
	self.c.State().GetRepoState().SetCurrentPopupOpts(nil)
	self.c.Mutexes().PopupMutex.Unlock()

	self.c.Views().Confirmation.Visible = false

	self.clearConfirmationViewKeyBindings()
}

func (self *ConfirmationHelper) DeactivatePrompt() {
	self.c.Mutexes().PopupMutex.Lock()
	self.c.State().GetRepoState().SetCurrentPopupOpts(nil)
	self.c.Mutexes().PopupMutex.Unlock()

	self.c.Views().Prompt.Visible = false
	self.c.Views().Suggestions.Visible = false

	self.clearPromptViewKeyBindings()
}

func getMessageHeight(wrap bool, editable bool, message string, width int, tabWidth int) int {
	wrappedLines, _, _ := utils.WrapViewLinesToWidth(wrap, editable, message, width, tabWidth)
	return len(wrappedLines)
}

func (self *ConfirmationHelper) getPopupPanelDimensionsForContentHeight(panelWidth, contentHeight int, parentPopupContext types.Context) (int, int, int, int) {
	return self.getPopupPanelDimensionsAux(panelWidth, contentHeight, parentPopupContext)
}

func (self *ConfirmationHelper) getPopupPanelDimensionsAux(panelWidth int, panelHeight int, parentPopupContext types.Context) (int, int, int, int) {
	width, height := self.c.GocuiGui().Size()
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
	}
	if parentPopupContext != nil {
		// If there's already a popup on the screen, offset the new one from its
		// parent so that it's clearly distinguished from the parent
		x0, y0, _, _ := parentPopupContext.GetView().Dimensions()
		x0 += 2
		y0 += 1
		return x0, y0, x0 + panelWidth, y0 + panelHeight + 1
	}
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (self *ConfirmationHelper) getPopupPanelWidth() int {
	width, _ := self.c.GocuiGui().Size()
	// we want a minimum width up to a point, then we do it based on ratio.
	panelWidth := 4 * width / 7
	minWidth := 80
	if panelWidth < minWidth {
		panelWidth = min(width-2, minWidth)
	}

	return panelWidth
}

func (self *ConfirmationHelper) prepareConfirmationPanel(
	opts types.ConfirmOpts,
) {
	self.c.Views().Confirmation.Title = opts.Title
	self.c.Views().Confirmation.FgColor = theme.GocuiDefaultTextColor

	self.c.ResetViewOrigin(self.c.Views().Confirmation)
	self.c.SetViewContent(self.c.Views().Confirmation, style.AttrBold.Sprint(strings.TrimSpace(opts.Prompt)))
}

func (self *ConfirmationHelper) preparePromptPanel(
	opts types.ConfirmOpts,
) {
	self.c.Views().Prompt.Title = opts.Title
	self.c.Views().Prompt.FgColor = theme.GocuiDefaultTextColor
	self.c.Views().Prompt.Mask = runeForMask(opts.Mask)
	self.c.Views().Prompt.SetOrigin(0, 0)

	textArea := self.c.Views().Prompt.TextArea
	textArea.Clear()
	textArea.TypeString(opts.Prompt)
	self.c.Views().Prompt.RenderTextArea()

	if opts.FindSuggestionsFunc != nil {
		suggestionsContext := self.c.Contexts().Suggestions
		suggestionsContext.State.FindSuggestions = opts.FindSuggestionsFunc
		suggestionsView := self.c.Views().Suggestions
		suggestionsView.Wrap = false
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		suggestionsContext.SetSuggestions(opts.FindSuggestionsFunc(""))
		suggestionsView.Visible = true
		suggestionsView.Title = fmt.Sprintf(self.c.Tr.SuggestionsTitle, self.c.UserConfig().Keybinding.Universal.TogglePanel)
		suggestionsView.Subtitle = ""
	}
}

func runeForMask(mask bool) rune {
	if mask {
		return '*'
	}
	return 0
}

func (self *ConfirmationHelper) CreatePopupPanel(ctx goContext.Context, opts types.CreatePopupPanelOpts) {
	self.c.Mutexes().PopupMutex.Lock()
	defer self.c.Mutexes().PopupMutex.Unlock()

	_, cancel := goContext.WithCancel(ctx)

	// we don't allow interruptions of non-loader popups in case we get stuck somehow
	// e.g. a credentials popup never gets its required user input so a process hangs
	// forever.
	// The proper solution is to have a queue of popup options
	currentPopupOpts := self.c.State().GetRepoState().GetCurrentPopupOpts()
	if currentPopupOpts != nil && !currentPopupOpts.HasLoader {
		self.c.Log.Error("ignoring create popup panel because a popup panel is already open")
		cancel()
		return
	}

	// remove any previous keybindings
	self.clearConfirmationViewKeyBindings()
	self.clearPromptViewKeyBindings()

	var context types.Context
	if opts.Editable {
		self.c.Contexts().Suggestions.State.FindSuggestions = opts.FindSuggestionsFunc

		self.preparePromptPanel(
			types.ConfirmOpts{
				Title:               opts.Title,
				Prompt:              opts.Prompt,
				FindSuggestionsFunc: opts.FindSuggestionsFunc,
				Mask:                opts.Mask,
			})

		context = self.c.Contexts().Prompt

		self.setPromptKeyBindings(cancel, opts)
	} else {
		if opts.FindSuggestionsFunc != nil {
			panic("non-editable confirmation views do not support suggestions")
		}

		self.c.Contexts().Suggestions.State.FindSuggestions = nil

		self.prepareConfirmationPanel(
			types.ConfirmOpts{
				Title:  opts.Title,
				Prompt: opts.Prompt,
			})

		context = self.c.Contexts().Confirmation

		self.setConfirmationKeyBindings(cancel, opts)
	}

	self.c.Contexts().Suggestions.State.AllowEditSuggestion = opts.AllowEditSuggestion

	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	self.c.Context().Push(context, types.OnFocusOpts{})
}

func (self *ConfirmationHelper) setConfirmationKeyBindings(cancel goContext.CancelFunc, opts types.CreatePopupPanelOpts) {
	onConfirm := self.wrappedConfirmationFunction(cancel, opts.HandleConfirm)
	onClose := self.wrappedConfirmationFunction(cancel, opts.HandleClose)

	self.c.Contexts().Confirmation.State.OnConfirm = onConfirm
	self.c.Contexts().Confirmation.State.OnClose = onClose
}

func (self *ConfirmationHelper) setPromptKeyBindings(cancel goContext.CancelFunc, opts types.CreatePopupPanelOpts) {
	onConfirm := self.wrappedPromptConfirmationFunction(cancel, opts.HandleConfirmPrompt,
		func() string { return self.c.Views().Prompt.TextArea.GetContent() },
		opts.AllowEmptyInput, opts.PreserveWhitespace)

	onSuggestionConfirm := self.wrappedPromptConfirmationFunction(
		cancel,
		opts.HandleConfirmPrompt,
		self.getSelectedSuggestionValue,
		opts.AllowEmptyInput,
		opts.PreserveWhitespace,
	)

	onClose := self.wrappedConfirmationFunction(cancel, opts.HandleClose)

	onDeleteSuggestion := func() error {
		if opts.HandleDeleteSuggestion == nil {
			return nil
		}

		idx := self.c.Contexts().Suggestions.GetSelectedLineIdx()
		return opts.HandleDeleteSuggestion(idx)
	}

	self.c.Contexts().Prompt.State.OnConfirm = onConfirm
	self.c.Contexts().Prompt.State.OnClose = onClose
	self.c.Contexts().Suggestions.State.OnConfirm = onSuggestionConfirm
	self.c.Contexts().Suggestions.State.OnClose = onClose
	self.c.Contexts().Suggestions.State.OnDeleteSuggestion = onDeleteSuggestion
}

func (self *ConfirmationHelper) clearConfirmationViewKeyBindings() {
	noop := func() error { return nil }
	self.c.Contexts().Confirmation.State.OnConfirm = noop
	self.c.Contexts().Confirmation.State.OnClose = noop
}

func (self *ConfirmationHelper) clearPromptViewKeyBindings() {
	noop := func() error { return nil }
	self.c.Contexts().Prompt.State.OnConfirm = noop
	self.c.Contexts().Prompt.State.OnClose = noop
	self.c.Contexts().Suggestions.State.OnConfirm = noop
	self.c.Contexts().Suggestions.State.OnClose = noop
	self.c.Contexts().Suggestions.State.OnDeleteSuggestion = noop
}

func (self *ConfirmationHelper) getSelectedSuggestionValue() string {
	selectedSuggestion := self.c.Contexts().Suggestions.GetSelected()

	if selectedSuggestion != nil {
		return selectedSuggestion.Value
	}

	return ""
}

func (self *ConfirmationHelper) ResizeCurrentPopupPanels() {
	var parentPopupContext types.Context
	for _, c := range self.c.Context().CurrentPopup() {
		switch c {
		case self.c.Contexts().Menu:
			self.resizeMenu(parentPopupContext)
		case self.c.Contexts().Confirmation:
			self.resizeConfirmationPanel(parentPopupContext)
		case self.c.Contexts().Prompt, self.c.Contexts().Suggestions:
			self.resizePromptPanel(parentPopupContext)
		case self.c.Contexts().CommitMessage, self.c.Contexts().CommitDescription:
			self.ResizeCommitMessagePanels(parentPopupContext)
		}

		parentPopupContext = c
	}
}

func (self *ConfirmationHelper) resizeMenu(parentPopupContext types.Context) {
	// we want the unfiltered length here so that if we're filtering we don't
	// resize the window
	itemCount := self.c.Contexts().Menu.UnfilteredLen()
	offset := 3
	panelWidth := self.getPopupPanelWidth()
	contentWidth := panelWidth - 2 // minus 2 for the frame
	promptLinesCount := self.layoutMenuPrompt(contentWidth)
	x0, y0, x1, y1 := self.getPopupPanelDimensionsForContentHeight(panelWidth, itemCount+offset+promptLinesCount, parentPopupContext)
	menuBottom := y1 - offset
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Menu.Name(), x0, y0, x1, menuBottom, 0)

	tooltipTop := menuBottom + 1
	tooltip := ""
	selectedItem := self.c.Contexts().Menu.GetSelected()
	if selectedItem != nil {
		tooltip = self.TooltipForMenuItem(selectedItem)
	}
	tooltipHeight := getMessageHeight(true, false, tooltip, contentWidth, self.c.Views().Menu.TabWidth) + 2 // plus 2 for the frame
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Tooltip.Name(), x0, tooltipTop, x1, tooltipTop+tooltipHeight-1, 0)
}

// Wraps the lines of the menu prompt to the available width and rerenders the
// menu if needed. Returns the number of lines the prompt takes up.
func (self *ConfirmationHelper) layoutMenuPrompt(contentWidth int) int {
	oldPromptLines := self.c.Contexts().Menu.GetPromptLines()
	var promptLines []string
	prompt := self.c.Contexts().Menu.GetPrompt()
	if len(prompt) > 0 {
		promptLines, _, _ = utils.WrapViewLinesToWidth(true, false, prompt, contentWidth, self.c.Views().Menu.TabWidth)
		promptLines = append(promptLines, "")
	}
	self.c.Contexts().Menu.SetPromptLines(promptLines)
	if len(oldPromptLines) != len(promptLines) {
		// The number of lines in the prompt has changed; this happens either
		// because we're now showing a menu that has a prompt, and the previous
		// menu didn't (or vice versa), or because the user is resizing the
		// terminal window while a menu with a prompt is open.

		// We need to rerender to give the menu context a chance to update its
		// non-model items, and reinitialize the data it uses for converting
		// between view index and model index.
		self.c.Contexts().Menu.HandleRender()

		// Then we need to refocus to ensure the cursor is in the right place in
		// the view.
		self.c.Contexts().Menu.HandleFocus(types.OnFocusOpts{})
	}
	return len(promptLines)
}

func (self *ConfirmationHelper) resizeConfirmationPanel(parentPopupContext types.Context) {
	panelWidth := self.getPopupPanelWidth()
	contentWidth := panelWidth - 2 // minus 2 for the frame
	confirmationView := self.c.Views().Confirmation
	prompt := confirmationView.Buffer()
	panelHeight := getMessageHeight(true, false, prompt, contentWidth, confirmationView.TabWidth)
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight, parentPopupContext)
	_, _ = self.c.GocuiGui().SetView(confirmationView.Name(), x0, y0, x1, y1, 0)
}

func (self *ConfirmationHelper) resizePromptPanel(parentPopupContext types.Context) {
	suggestionsViewHeight := 0
	if self.c.Views().Suggestions.Visible {
		suggestionsViewHeight = 11
	}
	panelWidth := self.getPopupPanelWidth()
	contentWidth := panelWidth - 2 // minus 2 for the frame
	promptView := self.c.Views().Prompt
	prompt := promptView.TextArea.GetContent()
	panelHeight := getMessageHeight(false, true, prompt, contentWidth, promptView.TabWidth) + suggestionsViewHeight
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight, parentPopupContext)
	promptViewBottom := y1 - suggestionsViewHeight
	_, _ = self.c.GocuiGui().SetView(promptView.Name(), x0, y0, x1, promptViewBottom, 0)

	suggestionsViewTop := promptViewBottom + 1
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (self *ConfirmationHelper) ResizeCommitMessagePanels(parentPopupContext types.Context) {
	panelWidth := self.getPopupPanelWidth()
	content := self.c.Views().CommitDescription.TextArea.GetContent()
	summaryViewHeight := 3
	panelHeight := getMessageHeight(false, true, content, panelWidth, self.c.Views().CommitDescription.TabWidth)
	minHeight := 7
	if panelHeight < minHeight {
		panelHeight = minHeight
	}
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight, parentPopupContext)

	_, _ = self.c.GocuiGui().SetView(self.c.Views().CommitMessage.Name(), x0, y0, x1, y0+summaryViewHeight-1, 0)
	_, _ = self.c.GocuiGui().SetView(self.c.Views().CommitDescription.Name(), x0, y0+summaryViewHeight, x1, y1+summaryViewHeight, 0)
}

func (self *ConfirmationHelper) IsPopupPanel(context types.Context) bool {
	return context.GetKind() == types.PERSISTENT_POPUP || context.GetKind() == types.TEMPORARY_POPUP
}

func (self *ConfirmationHelper) IsPopupPanelFocused() bool {
	return self.IsPopupPanel(self.c.Context().Current())
}

func (self *ConfirmationHelper) TooltipForMenuItem(menuItem *types.MenuItem) string {
	tooltip := menuItem.Tooltip
	if menuItem.DisabledReason != nil && menuItem.DisabledReason.Text != "" {
		if tooltip != "" {
			tooltip += "\n\n"
		}
		tooltip += style.FgRed.Sprintf(self.c.Tr.DisabledMenuItemPrefix) + menuItem.DisabledReason.Text
	}
	return tooltip
}
