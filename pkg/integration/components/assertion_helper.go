package components

import (
	"time"

	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

const eventuallyTimeout = 2 * time.Second

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

func (self *assertionHelper) assertEventually(test func() (bool, string)) {
	deadline := time.Now().Add(eventuallyTimeout)
	for {
		ok, message := test()
		if ok {
			return
		}
		if time.Now().After(deadline) {
			self.fail(message)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (self *assertionHelper) fail(message string) {
	self.gui.Fail(message)
}
