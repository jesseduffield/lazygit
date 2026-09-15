//nolint:revive
package xtime

import (
	"time"
)

func NewRealClock() *RealClock {
	return &RealClock{}
}

type RealClock struct {
	_ noCopy
}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

func (c *RealClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (c *RealClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (c *RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
