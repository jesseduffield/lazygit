package gocui

import (
	"testing"

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
