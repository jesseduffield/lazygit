// Copyright © 2016 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package jwalterweatherman

import (
	"sync/atomic"
)

type logCounter struct {
	counter uint64
}

func (c *logCounter) incr() {
	atomic.AddUint64(&c.counter, 1)
}

func (c *logCounter) resetCounter() {
	atomic.StoreUint64(&c.counter, 0)
}

func (c *logCounter) getCount() uint64 {
	return atomic.LoadUint64(&c.counter)
}

func (c *logCounter) Write(p []byte) (n int, err error) {
	c.incr()
	return len(p), nil
}

// LogCountForLevel returns the number of log invocations for a given threshold.
func (n *Notepad) LogCountForLevel(l Threshold) uint64 {
	return n.logCounters[l].getCount()
}

// LogCountForLevelsGreaterThanorEqualTo returns the number of log invocations
// greater than or equal to a given threshold.
func (n *Notepad) LogCountForLevelsGreaterThanorEqualTo(threshold Threshold) uint64 {
	var cnt uint64

	for i := int(threshold); i < len(n.logCounters); i++ {
		cnt += n.LogCountForLevel(Threshold(i))
	}

	return cnt
}

// ResetLogCounters resets the invocation counters for all levels.
func (n *Notepad) ResetLogCounters() {
	for _, np := range n.logCounters {
		np.resetCounter()
	}
}
