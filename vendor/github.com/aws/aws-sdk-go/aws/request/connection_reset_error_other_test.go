// +build appengine plan9

package request_test

import (
	"errors"
)

var (
	errAcceptConnectionResetStub = errors.New("accept: connection reset")
	errReadConnectionResetStub   = errors.New("read: connection reset")
)
