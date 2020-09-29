package utils

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// NewDummyLog creates a new dummy Log for testing
func NewDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}
