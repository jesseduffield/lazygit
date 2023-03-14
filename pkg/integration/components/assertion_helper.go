package components

import (
	"os"
	"time"

	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type assertionHelper struct {
	gui integrationTypes.GuiDriver
}

// milliseconds we'll wait when an assertion fails.
func retryWaitTimes() []int {
	if os.Getenv("LONG_WAIT_BEFORE_FAIL") == "true" {
		// CI has limited hardware, may be throttled, runs tests in parallel, etc, so we
		// give it more leeway compared to when we're running things locally.
		return []int{0, 1, 1, 1, 1, 1, 5, 10, 20, 40, 100, 200, 500, 1000, 2000, 4000}
	} else {
		return []int{0, 1, 1, 1, 1, 1, 5, 10, 20, 40, 100, 200}
	}
}

func (self *assertionHelper) matchString(matcher *Matcher, context string, getValue func() string) {
	self.assertWithRetries(func() (bool, string) {
		value := getValue()
		return matcher.context(context).test(value)
	})
}

func (self *assertionHelper) assertWithRetries(test func() (bool, string)) {
	var message string
	for _, waitTime := range retryWaitTimes() {
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
