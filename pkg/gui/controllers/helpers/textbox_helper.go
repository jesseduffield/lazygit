package helpers

import (
	goContext "context"

	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TextboxHelper struct {
	c *HelperCommon
}

func NewTextboxHelper(c *HelperCommon) *TextboxHelper {
	return &TextboxHelper{
		c: c,
	}
}

func (self *TextboxHelper) DeactivateTextboxPrompt() {
	self.c.Mutexes().PopupMutex.Lock()
	self.c.State().GetRepoState().SetCurrentPopupOpts(nil)
	self.c.Mutexes().PopupMutex.Unlock()

	self.c.Views().Textbox.Visible = false
	self.clearTextboxViewKeyBindings()
}

func (self *TextboxHelper) clearTextboxViewKeyBindings() {
	noop := func() error { return nil }
	self.c.Contexts().Textbox.State.OnConfirm = noop
	self.c.Contexts().Textbox.State.OnClose = noop
}

func (self *TextboxHelper) CreatePopupPanel(ctx goContext.Context, opts types.CreatePopupPanelOpts) error {
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

	textboxView := self.c.Views().Textbox

	textboxView.Title = opts.Title
	// Introduce confirm key bindings of textbox to users
	textboxView.Subtitle = utils.ResolvePlaceholderString(self.c.Tr.TextboxSubTitle,
		map[string]string{
			"textboxConfirmBinding": keybindings.Label(self.c.UserConfig.Keybinding.Universal.ConfirmInEditor),
		})

	textboxView.Wrap = !opts.Editable
	textboxView.FgColor = theme.GocuiDefaultTextColor
	textboxView.Mask = runeForMask(opts.Mask)

	// Set view position
	width := self.getPopupPanelWidth()
	height := self.getPopupPanelHeight()
	x0, y0, x1, y1 := self.getPosition(width, height)
	self.c.GocuiGui().SetView(textboxView.Name(), x0, y0, x1, y1, 0)

	// Render text in textbox
	textboxView.Editable = opts.Editable
	textArea := textboxView.TextArea
	textArea.Clear()
	textArea.TypeString(opts.Prompt)
	textboxView.RenderTextArea()

	// Setting Handlers
	self.c.Contexts().Textbox.State.OnConfirm = self.wrappedPromptTextboxFunction(cancel, opts.HandleConfirmPrompt, func() string { return self.c.Views().Textbox.TextArea.GetContent() })
	self.c.Contexts().Textbox.State.OnClose = self.wrappedTextboxFunction(cancel, opts.HandleClose)

	// Set text box to current popup
	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	return self.c.PushContext(self.c.Contexts().Textbox)
}

func (self *TextboxHelper) wrappedPromptTextboxFunction(cancel goContext.CancelFunc, function func(string) error, getResponse func() string) func() error {
	return self.wrappedTextboxFunction(cancel, func() error {
		return function(getResponse())
	})
}

func (self *TextboxHelper) wrappedTextboxFunction(cancel goContext.CancelFunc, function func() error) func() error {
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

func (self *TextboxHelper) getPosition(panelWidth int, panelHeight int) (int, int, int, int) {
	width, height := self.c.GocuiGui().Size()
	if panelHeight > height * 3 / 4 {
		panelHeight = height * 3 / 4
	}
	return width / 2 - panelWidth / 2,
		height / 2 - panelHeight / 2 - panelHeight % 2 - 1,
		width / 2 + panelWidth / 2,
		height / 2 + panelHeight / 2
}

func (self *TextboxHelper) getPopupPanelWidth() int {
	width, _ := self.c.GocuiGui().Size()
	panelWidth := 4 * width / 7
	minWidth := 80
	if panelWidth < minWidth {
		if width - 2 < minWidth {
			panelWidth = width - 2
		} else {
			panelWidth = minWidth
		}
	}

	return panelWidth
}

func (self *TextboxHelper) getPopupPanelHeight() int {
	_, height := self.c.GocuiGui().Size()
	var panelHeight int
	maxHeight := 11
	if height - 2 > maxHeight {
		panelHeight = maxHeight
	} else {
		panelHeight = height - 2
	}

	return panelHeight
}
