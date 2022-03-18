package test

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var _ logrus.FieldLogger = &FakeFieldLogger{}

// for now we're just tracking calls to the Error and Errorf methods
type FakeFieldLogger struct {
	loggedErrors []string
	*logrus.Entry
}

func (self *FakeFieldLogger) Error(args ...interface{}) {
	if len(args) != 1 {
		panic("Expected exactly one argument to FakeFieldLogger.Error")
	}

	switch arg := args[0].(type) {
	case error:
		self.loggedErrors = append(self.loggedErrors, arg.Error())
	case string:
		self.loggedErrors = append(self.loggedErrors, arg)
	}
}

func (self *FakeFieldLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	self.loggedErrors = append(self.loggedErrors, msg)
}

func (self *FakeFieldLogger) AssertErrors(t *testing.T, expectedErrors []string) {
	t.Helper()
	assert.EqualValues(t, expectedErrors, self.loggedErrors)
}
