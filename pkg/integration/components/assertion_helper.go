package components

import (
	"time"

	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type assertionHelper struct {
	gui integrationTypes.GuiDriver
}

// milliseconds we'll wait when an assertion fails.
var retryWaitTimes = []int{0, 1, 1, 1, 1, 1, 5, 10, 20, 40, 100, 200, 500, 1000, 2000, 4000}

func (self *assertionHelper) matchString(matcher *matcher, context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
	})
}

func (self *assertionHelper) assertWithRetries(test func() (bool, string)) {
	var message string
	for _, waitTime := range retryWaitTimes {
		time.Sleep(time.Duration(waitTime) * time.Millisecond)

		var ok bool
		ok, message = test()
		if ok {
			return
		}
	}

	self.fail(message)
}

func (self *assertionHelper) fail(message string) {
	self.gui.Fail(message)
}
