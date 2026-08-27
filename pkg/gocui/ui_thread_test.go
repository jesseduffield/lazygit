package gocui

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// errStillWaiting stands in for the result of a wait that hasn't produced one.
var errStillWaiting = errors.New("still waiting")

// resultOrTimeout reports what a wait returned, or errStillWaiting if it hasn't
// returned by the time we give up on it.
func resultOrTimeout(result chan error) error {
	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		return errStillWaiting
	}
}

// A worker waiting for the UI thread must not be left parked there once the
// main loop has stopped: nothing will ever run its callback, and the shutdown
// that follows blocks until such workers have finished (see
// tasks.ViewBufferManager.Close).
func TestOnUIThreadAndWaitGivesUpWhenTheLoopExits(t *testing.T) {
	g := newTestGui(t)

	// Closing this is what MainLoop returning does. From here on nothing
	// dequeues user events, so the callback below is never going to run.
	close(g.loopExited)

	result := make(chan error, 1)
	go func() {
		result <- g.OnUIThreadAndWait(func() {})
	}()

	err := resultOrTimeout(result)
	assert.ErrorIs(t, err, ErrLoopExited)
}
