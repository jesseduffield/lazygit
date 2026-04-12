//go:build !go1.24

package deadlock

type beforeAfter struct {
	before interface{}
	after  interface{}
}

func newBeforeAfter(before, after interface{}) beforeAfter {
	return beforeAfter{before: before, after: after}
}
