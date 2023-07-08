package components

import (
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type assertionHelper struct {
	gui integrationTypes.GuiDriver
}

func (self *assertionHelper) matchString(matcher *TextMatcher, context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
	})
}

// We no longer assert with retries now that lazygit tells us when it's no longer
// busy. But I'm keeping the function in case we want to re-introduce it later.
func (self *assertionHelper) assertWithRetries(test func() (bool, string)) {
	ok, message := test()
	if !ok {
		self.fail(message)
	}
}

func (self *assertionHelper) fail(message string) {
	self.gui.Fail(message)
}
