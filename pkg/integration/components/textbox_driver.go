package components

type TextboxDriver struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *TextboxDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Textbox()
}

func (self *TextboxDriver) Title(expected *TextMatcher) *TextboxDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

func (self *TextboxDriver) Type(value string) *TextboxDriver {
	self.t.typeContent(value)

	return self
}

func (self *TextboxDriver) NewLine() *TextboxDriver {
	self.getViewDriver().PressEnter()

	return self
}

func (self *TextboxDriver) Confirm() *TextboxDriver {
	self.getViewDriver().PressAltEnter()

	return self
}

func (self *TextboxDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.t.Fail("You must check the title of a menu popup by calling Title() before calling Confirm()/Cancel().")
	}
}
