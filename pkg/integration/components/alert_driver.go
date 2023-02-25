package components

type AlertDriver struct {
	t                 *TestDriver
	hasCheckedTitle   bool
	hasCheckedContent bool
}

func (self *AlertDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Confirmation()
}

// asserts that the alert view has the expected title
func (self *AlertDriver) Title(expected *Matcher) *AlertDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

// asserts that the alert view has the expected content
func (self *AlertDriver) Content(expected *Matcher) *AlertDriver {
	self.getViewDriver().Content(expected)

	self.hasCheckedContent = true

	return self
}

func (self *AlertDriver) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEnter()
}

func (self *AlertDriver) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEscape()
}

func (self *AlertDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedContent || !self.hasCheckedTitle {
		self.t.Fail("You must both check the content and title of a confirmation popup by calling Title()/Content() before calling Confirm()/Cancel().")
	}
}
