package helpers

import (
	goContext "context"
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"

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

		if err := self.c.PopContext(); err != nil {
			return err
		}

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
	if !wrap {
		return len(strings.Split(message, "\n"))
	}

	lineCount := 0
	lines := strings.Split(message, "\n")

	for _, line := range lines {
		n := 0
		lastWhitespaceIndex := -1
		for i, currChr := range line {
			rw := runewidth.RuneWidth(currChr)
			n += rw

			if n > width {
				if currChr == ' ' {
					n = 0
				} else if currChr == '-' {
					n = rw
				} else if lastWhitespaceIndex != -1 && lastWhitespaceIndex+1 != i {
					if line[lastWhitespaceIndex] == '-' {
						n = i - lastWhitespaceIndex
					} else {
						n = i - lastWhitespaceIndex + 1
					}
				} else {
					n = rw
				}
				lineCount++
				lastWhitespaceIndex = -1
			} else if currChr == ' ' || currChr == '-' {
				lastWhitespaceIndex = i
			}
		}
		lineCount++
	}

	return lineCount
}

func (self *ConfirmationHelper) getPopupPanelDimensions(wrap bool, prompt string) (int, int, int, int) {
	panelWidth := self.getPopupPanelWidth()
	panelHeight := getMessageHeight(wrap, prompt, panelWidth)
	return self.getPopupPanelDimensionsAux(panelWidth, panelHeight)
}

func (self *ConfirmationHelper) getPopupPanelDimensionsForContentHeight(panelWidth, contentHeight int) (int, int, int, int) {
	return self.getPopupPanelDimensionsAux(panelWidth, contentHeight)
}

func (self *ConfirmationHelper) getPopupPanelDimensionsAux(panelWidth int, panelHeight int) (int, int, int, int) {
	width, height := self.c.GocuiGui().Size()
	if panelHeight > height*3/4 {
		panelHeight = height * 3 / 4
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
) error {
	self.c.Views().Confirmation.Title = opts.Title
	// for now we do not support wrapping in our editor
	self.c.Views().Confirmation.Wrap = !opts.Editable
	self.c.Views().Confirmation.FgColor = theme.GocuiDefaultTextColor
	self.c.Views().Confirmation.Mask = runeForMask(opts.Mask)
	_ = self.c.Views().Confirmation.SetOrigin(0, 0)

	suggestionsContext := self.c.Contexts().Suggestions
	suggestionsContext.State.FindSuggestions = opts.FindSuggestionsFunc
	if opts.FindSuggestionsFunc != nil {
		suggestionsView := self.c.Views().Suggestions
		suggestionsView.Wrap = false
		suggestionsView.FgColor = theme.GocuiDefaultTextColor
		suggestionsContext.SetSuggestions(opts.FindSuggestionsFunc(""))
		suggestionsView.Visible = true
		suggestionsView.Title = fmt.Sprintf(self.c.Tr.SuggestionsTitle, self.c.UserConfig.Keybinding.Universal.TogglePanel)
		suggestionsView.Subtitle = ""
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

	_, cancel := goContext.WithCancel(ctx)

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
		types.ConfirmOpts{
			Title:               opts.Title,
			Prompt:              opts.Prompt,
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
		self.c.SetViewContent(confirmationView, style.AttrBold.Sprint(underlineLinks(opts.Prompt)))
	}

	if err := self.setKeyBindings(cancel, opts); err != nil {
		cancel()
		return err
	}

	self.c.Contexts().Suggestions.State.AllowEditSuggestion = opts.AllowEditSuggestion

	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	return self.c.PushContext(self.c.Contexts().Confirmation)
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
		underlinedLink := style.AttrUnderline.Sprint(remaining[linkStart:linkEnd])
		if strings.HasSuffix(underlinedLink, "\x1b[0m") {
			// Replace the "all styles off" code with "underline off" code
			underlinedLink = underlinedLink[:len(underlinedLink)-2] + "24m"
		}
		result += remaining[:linkStart] + underlinedLink
		remaining = remaining[linkEnd:]
	}
	return result + remaining
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

	return nil
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

func (self *ConfirmationHelper) ResizeConfirmationPanel() {
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
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (self *ConfirmationHelper) ResizeCurrentPopupPanel() error {
	c := self.c.CurrentContext()

	switch c {
	case self.c.Contexts().Menu:
		self.resizeMenu()
	case self.c.Contexts().Confirmation, self.c.Contexts().Suggestions:
		self.resizeConfirmationPanel()
	case self.c.Contexts().CommitMessage, self.c.Contexts().CommitDescription:
		self.ResizeCommitMessagePanels()
	}

	return nil
}

func (self *ConfirmationHelper) ResizePopupPanel(v *gocui.View, content string) error {
	x0, y0, x1, y1 := self.getPopupPanelDimensions(v.Wrap, content)
	_, err := self.c.GocuiGui().SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

func (self *ConfirmationHelper) resizeMenu() {
	// we want the unfiltered length here so that if we're filtering we don't
	// resize the window
	itemCount := self.c.Contexts().Menu.UnfilteredLen()
	offset := 3
	panelWidth := self.getPopupPanelWidth()
	x0, y0, x1, y1 := self.getPopupPanelDimensionsForContentHeight(panelWidth, itemCount+offset)
	menuBottom := y1 - offset
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Menu.Name(), x0, y0, x1, menuBottom, 0)

	tooltipTop := menuBottom + 1
	tooltip := ""
	selectedItem := self.c.Contexts().Menu.GetSelected()
	if selectedItem != nil {
		tooltip = self.TooltipForMenuItem(selectedItem)
	}
	tooltipHeight := getMessageHeight(true, tooltip, panelWidth) + 2 // plus 2 for the frame
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Tooltip.Name(), x0, tooltipTop, x1, tooltipTop+tooltipHeight-1, 0)
}

func (self *ConfirmationHelper) resizeConfirmationPanel() {
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
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight)
	confirmationViewBottom := y1 - suggestionsViewHeight
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Confirmation.Name(), x0, y0, x1, confirmationViewBottom, 0)

	suggestionsViewTop := confirmationViewBottom + 1
	_, _ = self.c.GocuiGui().SetView(self.c.Views().Suggestions.Name(), x0, suggestionsViewTop, x1, suggestionsViewTop+suggestionsViewHeight, 0)
}

func (self *ConfirmationHelper) ResizeCommitMessagePanels() {
	panelWidth := self.getPopupPanelWidth()
	content := self.c.Views().CommitDescription.TextArea.GetContent()
	summaryViewHeight := 3
	panelHeight := getMessageHeight(false, content, panelWidth)
	minHeight := 7
	if panelHeight < minHeight {
		panelHeight = minHeight
	}
	x0, y0, x1, y1 := self.getPopupPanelDimensionsAux(panelWidth, panelHeight)

	_, _ = self.c.GocuiGui().SetView(self.c.Views().CommitMessage.Name(), x0, y0, x1, y0+summaryViewHeight-1, 0)
	_, _ = self.c.GocuiGui().SetView(self.c.Views().CommitDescription.Name(), x0, y0+summaryViewHeight, x1, y1+summaryViewHeight, 0)
}

func (self *ConfirmationHelper) IsPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "confirmation" || viewName == "menu"
}

func (self *ConfirmationHelper) IsPopupPanelFocused() bool {
	return self.IsPopupPanel(self.c.CurrentContext().GetViewName())
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
