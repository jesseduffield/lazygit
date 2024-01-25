package components

type MenuDriver struct {
	t               *TestDriver
	hasCheckedTitle bool
}

func (self *MenuDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Menu()
}

// asserts that the popup has the expected title
func (self *MenuDriver) Title(expected *TextMatcher) *MenuDriver {
	self.getViewDriver().Title(expected)

	self.hasCheckedTitle = true

	return self
}

func (self *MenuDriver) Confirm() *MenuDriver {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEnter()

	return self
}

func (self *MenuDriver) Cancel() {
	self.checkNecessaryChecksCompleted()

	self.getViewDriver().PressEscape()
}

func (self *MenuDriver) Select(option *TextMatcher) *MenuDriver {
	self.getViewDriver().NavigateToLine(option)

	return self
}

func (self *MenuDriver) Lines(matchers ...*TextMatcher) *MenuDriver {
	self.getViewDriver().Lines(matchers...)

	return self
}

func (self *MenuDriver) TopLines(matchers ...*TextMatcher) *MenuDriver {
	self.getViewDriver().TopLines(matchers...)

	return self
}

func (self *MenuDriver) Filter(text string) *MenuDriver {
	self.getViewDriver().FilterOrSearch(text)

	return self
}

func (self *MenuDriver) LineCount(matcher *IntMatcher) *MenuDriver {
	self.getViewDriver().LineCount(matcher)

	return self
}

func (self *MenuDriver) Wait(milliseconds int) *MenuDriver {
	self.getViewDriver().Wait(milliseconds)

	return self
}

func (self *MenuDriver) Tooltip(option *TextMatcher) *MenuDriver {
	self.t.Views().Tooltip().Content(option)

	return self
}

func (self *MenuDriver) Tap(f func()) *MenuDriver {
	self.getViewDriver().Tap(f)
	return self
}

func (self *MenuDriver) checkNecessaryChecksCompleted() {
	if !self.hasCheckedTitle {
		self.t.Fail("You must check the title of a menu popup by calling Title() before calling Confirm()/Cancel().")
	}
}
