//nolint:revive
package xtime

import (
	"time"
)

func NewFakeClock() *FakeClock {
	return NewFakeClockAt(time.Now())
}

func NewFakeClockAt(t time.Time) *FakeClock {
	return &FakeClock{
		time: t,
	}
}

type FakeClock struct {
	_ noCopy

	// Not protected by a mutex. If a warning is thrown in your tests,
	// just disable parallel tests.
	time time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.time
}

func (c *FakeClock) Since(t time.Time) time.Duration {
	return c.time.Sub(t)
}

func (c *FakeClock) Until(t time.Time) time.Duration {
	return t.Sub(c.time)
}

func (c *FakeClock) Sleep(d time.Duration) {
	c.time = c.time.Add(d)
}
