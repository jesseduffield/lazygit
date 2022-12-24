package components

// for making assertions on string values
type matcher struct {
	// e.g. "contains 'foo'"
	name string
	// returns a bool that says whether the test passed and if it returns false, it
	// also returns a string of the error message
	testFn func(string) (bool, string)
	// this is printed when there's an error so that it's clear what the context of the assertion is
	prefix string
}

func NewMatcher(name string, testFn func(string) (bool, string)) *matcher {
	return &matcher{name: name, testFn: testFn}
}

func (self *matcher) test(value string) (bool, string) {
	ok, message := self.testFn(value)
	if ok {
		return true, ""
	}

	if self.prefix != "" {
		return false, self.prefix + " " + message
	}

	return false, message
}

func (self *matcher) context(prefix string) *matcher {
	self.prefix = prefix

	return self
}
