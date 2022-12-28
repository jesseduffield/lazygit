package components

type MenuDriver struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *MenuDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Menu()
}

// asserts that the popup has the expected title
func (self *MenuDriver) Title(expected *matcher) *MenuDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

func (self *MenuDriver) Confirm() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEnter()
}

func (self *MenuDriver) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEscape()
}

func (self *MenuDriver) Select(option *matcher) *MenuDriver {
	self.getViewDriver().NavigateToListItem(option)

	return self
}

func (self *MenuDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.t.Fail("You must check the title of a menu popup by calling Title() before calling Confirm()/Cancel().")
	}
}
