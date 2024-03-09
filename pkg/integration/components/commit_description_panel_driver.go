package components

type CommitDescriptionPanelDriver struct {
	t *TestDriver
}

func (self *CommitDescriptionPanelDriver) getViewDriver() *ViewDriver {
	return self.t.Views().CommitDescription()
}

// asserts on the current context of the description
func (self *CommitDescriptionPanelDriver) Content(expected *TextMatcher) *CommitDescriptionPanelDriver {
	self.getViewDriver().Content(expected)

	return self
}

func (self *CommitDescriptionPanelDriver) Type(value string) *CommitDescriptionPanelDriver {
	self.t.typeContent(value)

	return self
}

func (self *CommitDescriptionPanelDriver) SwitchToSummary() *CommitMessagePanelDriver {
	self.getViewDriver().PressTab()
	return &CommitMessagePanelDriver{t: self.t}
}

func (self *CommitDescriptionPanelDriver) AddNewline() *CommitDescriptionPanelDriver {
	self.t.pressFast(self.t.keys.Universal.Confirm)
	return self
}

func (self *CommitDescriptionPanelDriver) GoToBeginning() *CommitDescriptionPanelDriver {
	numLines := len(self.getViewDriver().getView().BufferLines())
	for i := 0; i < numLines; i++ {
		self.t.pressFast("<up>")
	}

	self.t.pressFast("<c-a>")
	return self
}

func (self *CommitDescriptionPanelDriver) Title(expected *TextMatcher) *CommitDescriptionPanelDriver {
	self.getViewDriver().Title(expected)

	return self
}

func (self *CommitDescriptionPanelDriver) Cancel() {
	self.getViewDriver().PressEscape()
}
