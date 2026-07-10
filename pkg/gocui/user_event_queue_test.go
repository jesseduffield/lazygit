package gocui

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Enqueuing far more events than the old fixed 256-slot buffer, without the
// main loop draining them, used to panic ("userEvents channel full"). It must
// not: producers can legitimately burst faster than a stalled UI loop drains
// (e.g. one command-log entry per git command when adding a large directory to
// a custom patch, or any producer while the loop is blocked in a subprocess).
// The events must also stay in FIFO order.
func TestUpdateIsUnboundedAndPreservesOrder(t *testing.T) {
	g := newTestGui(t)

	const n = 1000
	var got []int
	for i := range n {
		g.Update(func(*Gui) error {
			got = append(got, i)
			return nil
		})
	}

	// Drain the whole queue the way the main loop's inner drain does.
	_, err := g.processRemainingEvents()
	assert.NoError(t, err)

	want := make([]int, n)
	for i := range want {
		want[i] = i
	}
	assert.Equal(t, want, got)
}

// Concurrent producers must be able to enqueue safely (run under -race). Only
// same-goroutine order is guaranteed, so we check that every event is delivered
// exactly once and that each producer's own events stay in order.
func TestUpdateConcurrentProducers(t *testing.T) {
	g := newTestGui(t)

	const producers = 8
	const perProducer = 500

	type item struct{ producer, seq int }
	var got []item

	var wg sync.WaitGroup
	for p := range producers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seq := range perProducer {
				g.Update(func(*Gui) error {
					got = append(got, item{p, seq})
					return nil
				})
			}
		}()
	}
	// Update is a synchronous, non-blocking enqueue, so once every producer has
	// returned, every event is in the queue and a single drain sees them all.
	wg.Wait()

	_, err := g.processRemainingEvents()
	assert.NoError(t, err)

	assert.Len(t, got, producers*perProducer)
	lastSeq := make([]int, producers)
	for p := range lastSeq {
		lastSeq[p] = -1
	}
	for _, it := range got {
		assert.Equal(t, lastSeq[it.producer]+1, it.seq, "producer %d events out of order", it.producer)
		lastSeq[it.producer] = it.seq
	}
}
