package components

// TODO: soft-code this
const ClearKey = "<c-u>"

type SearchDriver struct {
	t *TestDriver
}

func (self *SearchDriver) getViewDriver() *ViewDriver {
	return self.t.Views().Search()
}

// asserts on the text initially present in the prompt
func (self *SearchDriver) InitialText(expected *Matcher) *SearchDriver {
	self.getViewDriver().Content(expected)

	return self
}

func (self *SearchDriver) Type(value string) *SearchDriver {
	self.t.typeContent(value)

	return self
}

func (self *SearchDriver) Clear() *SearchDriver {
	self.t.press(ClearKey)

	return self
}

func (self *SearchDriver) Confirm() {
	self.getViewDriver().PressEnter()
}

func (self *SearchDriver) Cancel() {
	self.getViewDriver().PressEscape()
}
