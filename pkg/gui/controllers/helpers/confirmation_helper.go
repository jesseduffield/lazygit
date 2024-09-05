package helpers

import (
	goContext "context"
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/mattn/go-runewidth"
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

func (self *ConfirmationHelper) wrappedConfirmationFunction(cancel goContext.CancelFunc, function func() error) func() error {
	return func() error {
		cancel()

		self.c.Context().Pop()

		if function != nil {
			if err := function(); err != nil {
				return err
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

// Temporary hack: we're just duplicating the logic in `gocui.lineWrap`
func getMessageHeight(wrap bool, message string, width int) int {
	return len(wrapMessageToWidth(wrap, message, width))
}

func wrapMessageToWidth(wrap bool, message string, width int) []string {
	lines := strings.Split(message, "\n")
	if !wrap {
		return lines
	}

	wrappedLines := make([]string, 0, len(lines))

	for _, line := range lines {
		n := 0
		offset := 0
		lastWhitespaceIndex := -1
		for i, currChr := range line {
			rw := runewidth.RuneWidth(currChr)
			n += rw

			if n > width {
				if currChr == ' ' {
					wrappedLines = append(wrappedLines, line[offset:i])
					offset = i + 1
					n = 0
				} else if currChr == '-' {
					wrappedLines = append(wrappedLines, line[offset:i])
					offset = i
					n = rw
				} else if lastWhitespaceIndex != -1 && lastWhitespaceIndex+1 != i {
					if line[lastWhitespaceIndex] == '-' {
						wrappedLines = append(wrappedLines, line[offset:lastWhitespaceIndex+1])
						offset = lastWhitespaceIndex + 1
						n = i - lastWhitespaceIndex
					} else {
						wrappedLines = append(wrappedLines, line[offset:lastWhitespaceIndex])
						offset = lastWhitespaceIndex + 1
						n = i - lastWhitespaceIndex + 1
					}
				} else {
					wrappedLines = append(wrappedLines, line[offset:i])
					offset = i
					n = rw
				}
				lastWhitespaceIndex = -1
			} else if currChr == ' ' || currChr == '-' {
				lastWhitespaceIndex = i
			}
		}

		wrappedLines = append(wrappedLines, line[offset:])
	}

	return wrappedLines
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
		if width-2 < minWidth {
			panelWidth = width - 2
		} else {
			panelWidth = minWidth
		}
	}

	return panelWidth
}

func (self *ConfirmationHelper) prepareConfirmationPanel(
	opts types.ConfirmOpts,
) {
	self.c.Views().Confirmation.Title = opts.Title
	// for now we do not support wrapping in our editor
	self.c.Views().Confirmation.Wrap = !opts.Editable
	self.c.Views().Confirmation.FgColor = theme.GocuiDefaultTextColor
	self.c.Views().Confirmation.Mask = runeForMask(opts.Mask)
	self.c.Views().Confirmation.SetOrigin(0, 0)

	suggestionsContext := self.c.Contexts().Suggestions
	suggestionsContext.State.FindSuggestions = opts.FindSuggestionsFunc
	if opts.FindSuggestionsFunc != nil {
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

	self.prepareConfirmationPanel(
		types.ConfirmOpts{
			Title:               opts.Title,
			Prompt:              opts.Prompt,
			FindSuggestionsFunc: opts.FindSuggestionsFunc,
			Editable:            opts.Editable,
			Mask:                opts.Mask,
		})
	confirmationView := self.c.Views().Confirmation
	confirmationView.Editable = opts.Editable

	if opts.Editable {
		textArea := confirmationView.TextArea
		textArea.Clear()
		textArea.TypeString(opts.Prompt)
		confirmationView.RenderTextArea()
	} else {
		self.c.ResetViewOrigin(confirmationView)
		self.c.SetViewContent(confirmationView, style.AttrBold.Sprint(underlineLinks(opts.Prompt)))
	}

	self.setKeyBindings(cancel, opts)

	self.c.Contexts().Suggestions.State.AllowEditSuggestion = opts.AllowEditSuggestion

	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	self.c.Context().Push(self.c.Contexts().Confirmation)
}

func underlineLinks(text string) string {
	result := ""
	remaining := text
	for {
		linkStart := strings.Index(remaining, "https://")
		if linkStart == -1 {
			break
		}

		linkEnd := strings.IndexAny(remaining[linkStart:], " \n>")
		if linkEnd == -1 {
			linkEnd = len(remaining)
		} else {
			linkEnd += linkStart
		}
		underlinedLink := style.PrintSimpleHyperlink(remaining[linkStart:linkEnd])
		result += remaining[:linkStart] + underlinedLink
		remaining = remaining[linkEnd:]
	}
	return result + remaining
}

func (self *ConfirmationHelper) setKeyBindings(cancel goContext.CancelFunc, opts types.CreatePopupPanelOpts) {
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

	onDeleteSuggestion := func() error {
		if opts.HandleDeleteSuggestion == nil {
			return nil
		}

		idx := self.c.Contexts().Suggestions.GetSelectedLineIdx()
		return opts.HandleDeleteSuggestion(idx)
	}

	self.c.Contexts().Confirmation.State.OnConfirm = onConfirm
	self.c.Contexts().Confirmation.State.OnClose = onClose
	self.c.Contexts().Suggestions.State.OnConfirm = onSuggestionConfirm
	self.c.Contexts().Suggestions.State.OnClose = onClose
	self.c.Contexts().Suggestions.State.OnDeleteSuggestion = onDeleteSuggestion
}

func (self *ConfirmationHelper) clearConfirmationViewKeyBindings() {
	noop := func() error { return nil }
	self.c.Contexts().Confirmation.State.OnConfirm = noop
	self.c.Contexts().Confirmation.State.OnClose = noop
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
		case self.c.Contexts().Confirmation, self.c.Contexts().Suggestions:
			self.resizeConfirmationPanel(parentPopupContext)
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
	tooltipHeight := getMessageHeight(true, tooltip, contentWidth) + 2 // plus 2 for the frame
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Tooltip.Name(), x0, tooltipTop, x1, tooltipTop+tooltipHeight-1, 0)
}

// Wraps the lines of the menu prompt to the available width and rerenders the
// menu if needed. Returns the number of lines the prompt takes up.
func (self *ConfirmationHelper) layoutMenuPrompt(contentWidth int) int {
	oldPromptLines := self.c.Contexts().Menu.GetPromptLines()
	var promptLines []string
	prompt := self.c.Contexts().Menu.GetPrompt()
	if len(prompt) > 0 {
		promptLines = wrapMessageToWidth(true, prompt, contentWidth)
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
	suggestionsViewHeight := 0
	if self.c.Views().Suggestions.Visible {
		suggestionsViewHeight = 11
	}
	panelWidth := self.getPopupPanelWidth()
	prompt := self.c.Views().Confirmation.Buffer()
	wrap := true
	if self.c.Views().Confirmation.Editable {
		prompt = self.c.Views().Confirmation.TextArea.GetContent()
		wrap = false
	}
	panelHeight := getMessageHeight(wrap, prompt, panelWidth) + suggestionsViewHeight
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight, parentPopupContext)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (self *ConfirmationHelper) ResizeCommitMessagePanels(parentPopupContext types.Context) {
	panelWidth := self.getPopupPanelWidth()
	content := self.c.Views().CommitDescription.TextArea.GetContent()
	summaryViewHeight := 3
	panelHeight := getMessageHeight(false, content, panelWidth)
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
	if menuItem.DisabledReason != nil {
		if tooltip != "" {
			tooltip += "\n\n"
		}
		tooltip += style.FgRed.Sprintf(self.c.Tr.DisabledMenuItemPrefix) + menuItem.DisabledReason.Text
	}
	return tooltip
}
