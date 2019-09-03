// +build !appengine,!js,!windows,aix

package logrus

import "io"

func checkIfTerminal(w io.Writer) bool {
	return false
}
