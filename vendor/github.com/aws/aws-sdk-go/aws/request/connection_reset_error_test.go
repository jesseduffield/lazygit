// +build !appengine,!plan9

package request_test

import (
	"net"
	"syscall"
)

var (
	errAcceptConnectionResetStub = &net.OpError{Op: "accept", Err: syscall.ECONNRESET}
	errReadConnectionResetStub   = &net.OpError{Op: "read", Err: syscall.ECONNRESET}
)
