package gocui

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskManagerHasBusyForegroundTaskExcept(t *testing.T) {
	t.Run("no tasks", func(t *testing.T) {
		tm := newTaskManager()
		assert.False(t, tm.hasBusyForegroundTaskExcept(nil))
	})

	t.Run("a busy foreground task counts", func(t *testing.T) {
		tm := newTaskManager()
		tm.NewTask(false)
		assert.True(t, tm.hasBusyForegroundTaskExcept(nil))
	})

	t.Run("a busy background task does not count", func(t *testing.T) {
		tm := newTaskManager()
		tm.NewTask(true)
		assert.False(t, tm.hasBusyForegroundTaskExcept(nil))
	})

	t.Run("a done foreground task does not count", func(t *testing.T) {
		tm := newTaskManager()
		task := tm.NewTask(false)
		task.Done()
		assert.False(t, tm.hasBusyForegroundTaskExcept(nil))
	})

	t.Run("a paused foreground task does not count", func(t *testing.T) {
		tm := newTaskManager()
		task := tm.NewTask(false)
		task.Pause()
		assert.False(t, tm.hasBusyForegroundTaskExcept(nil))
	})

	t.Run("the ignored task does not count", func(t *testing.T) {
		tm := newTaskManager()
		task := tm.NewTask(false)
		assert.False(t, tm.hasBusyForegroundTaskExcept(task))
	})

	t.Run("another foreground task counts even when one is ignored", func(t *testing.T) {
		tm := newTaskManager()
		ignored := tm.NewTask(false)
		tm.NewTask(false)
		assert.True(t, tm.hasBusyForegroundTaskExcept(ignored))
	})

	t.Run("only a background task alongside the ignored current event", func(t *testing.T) {
		// This is the repo-switch case: the switch is handled as the current
		// event (ignored) while a background refresh is in flight; it must not
		// be considered busy.
		tm := newTaskManager()
		current := tm.NewTask(false)
		tm.NewTask(true)
		assert.False(t, tm.hasBusyForegroundTaskExcept(current))
	})
}

func TestTaskManagerWaitUntilIdle(t *testing.T) {
	// returnsWithin reports whether f returns within the given duration.
	returnsWithin := func(d time.Duration, f func()) bool {
		done := make(chan struct{})
		go func() {
			f()
			close(done)
		}()
		select {
		case <-done:
			return true
		case <-time.After(d):
			return false
		}
	}

	t.Run("returns immediately when no task was ever created", func(t *testing.T) {
		tm := newTaskManager()
		assert.True(t, returnsWithin(time.Second, tm.WaitUntilIdle))
	})

	t.Run("blocks while a task is busy", func(t *testing.T) {
		tm := newTaskManager()
		tm.NewTask(false)
		assert.False(t, returnsWithin(50*time.Millisecond, tm.WaitUntilIdle))
	})

	t.Run("wakes up when the last busy task completes", func(t *testing.T) {
		tm := newTaskManager()
		task := tm.NewTask(false)
		go func() {
			time.Sleep(10 * time.Millisecond)
			task.Done()
		}()
		assert.True(t, returnsWithin(time.Second, tm.WaitUntilIdle))
	})

	t.Run("a paused task counts as idle", func(t *testing.T) {
		tm := newTaskManager()
		task := tm.NewTask(false)
		task.Pause()
		assert.True(t, returnsWithin(time.Second, tm.WaitUntilIdle))
	})

	t.Run("a task completing while nobody waits must not block", func(t *testing.T) {
		// This is the deadlock case: the waiter (the integration-test runner)
		// is between waits, and itself needs the task manager's mutex (it
		// creates a task whenever it enqueues work) before it waits again. The
		// idle notification must neither block the completing task while it
		// holds the mutex, nor get lost.
		tm := newTaskManager()
		assert.True(t, returnsWithin(time.Second, func() {
			// the program goes idle with nobody waiting...
			tm.NewTask(true).Done()

			// ...and creating and completing more tasks afterwards must still
			// be possible
			task := tm.NewTask(false)
			tm.NewTask(false).Done()
			task.Done()
		}))
		assert.True(t, returnsWithin(time.Second, tm.WaitUntilIdle))
	})
}
