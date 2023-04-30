package components

type CommitMessagePanelDriver struct {
	t *TestDriver
}

func (self *CommitMessagePanelDriver) getViewDriver() *ViewDriver {
	return self.t.Views().CommitMessage()
}

// asserts on the text initially present in the prompt
func (self *CommitMessagePanelDriver) InitialText(expected *Matcher) *CommitMessagePanelDriver {
	return self.Content(expected)
}

// asserts on the current context in the prompt
func (self *CommitMessagePanelDriver) Content(expected *Matcher) *CommitMessagePanelDriver {
	self.getViewDriver().Content(expected)

	return self
}

// asserts that the confirmation view has the expected title
func (self *CommitMessagePanelDriver) Title(expected *Matcher) *CommitMessagePanelDriver {
	self.getViewDriver().Title(expected)

	return self
}

func (self *CommitMessagePanelDriver) Type(value string) *CommitMessagePanelDriver {
	self.t.typeContent(value)

	return self
}

func (self *CommitMessagePanelDriver) SwitchToDescription() *CommitDescriptionPanelDriver {
	self.getViewDriver().PressTab()
	return &CommitDescriptionPanelDriver{t: self.t}
}

func (self *CommitMessagePanelDriver) AddNewline() *CommitMessagePanelDriver {
	self.t.press(self.t.keys.Universal.Confirm)

	return self
}

func (self *CommitMessagePanelDriver) Clear() *CommitMessagePanelDriver {
	// clearing multiple times in case there's multiple lines
	//  (the clear button only clears a single line at a time)
	maxAttempts := 100
	for i := 0; i < maxAttempts+1; i++ {
		if self.getViewDriver().getView().Buffer() == "" {
			break
		}

		self.t.press(ClearKey)
		if i == maxAttempts {
			panic("failed to clear commit message panel")
		}
	}

	return self
}

func (self *CommitMessagePanelDriver) Confirm() {
	self.getViewDriver().PressEnter()
}

func (self *CommitMessagePanelDriver) Close() {
	self.getViewDriver().PressEscape()
}

func (self *CommitMessagePanelDriver) Cancel() {
	self.getViewDriver().PressEscape()
}

func (self *CommitMessagePanelDriver) SelectPreviousMessage() *CommitMessagePanelDriver {
	self.getViewDriver().SelectPreviousItem()
	return self
}

func (self *CommitMessagePanelDriver) SelectNextMessage() *CommitMessagePanelDriver {
	self.getViewDriver().SelectNextItem()
	return self
}
