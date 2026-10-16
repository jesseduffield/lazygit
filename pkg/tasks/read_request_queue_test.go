package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequestQueueHandsBackWhatNoTaskWillServe(t *testing.T) {
	queue := newReadRequestQueue()

	// No task has been asked for, so the request comes straight back to its caller.
	assert.False(t, queue.enqueue(LinesToRead{Total: 1}))

	// Once one has, requests are held for it, and the task serves its own first
	// request ahead of them.
	queue.hold()
	assert.True(t, queue.enqueue(LinesToRead{Total: 2}))
	queue.beginServing(LinesToRead{Total: 1})
	assert.True(t, queue.enqueue(LinesToRead{Total: 3}))

	request, ok := queue.dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, request.Total)

	// A request that a task was serving when it was stopped goes back ahead of the
	// ones made after it.
	queue.requeue(request)
	request, _ = queue.dequeue()
	assert.Equal(t, 1, request.Total)
	request, _ = queue.dequeue()
	assert.Equal(t, 2, request.Total)

	// What no task got to comes back when the queue stops holding requests, and
	// nothing is taken from a caller after that.
	unanswered := queue.stopServing()
	assert.Len(t, unanswered, 1)
	assert.Equal(t, 3, unanswered[0].Total)

	assert.False(t, queue.enqueue(LinesToRead{Total: 4}))
	_, ok = queue.dequeue()
	assert.False(t, ok)
}

func TestReadRequestQueueRingsTheDoorbell(t *testing.T) {
	queue := newReadRequestQueue()
	doorbell := queue.beginServing(LinesToRead{Total: 1})

	// The task's first request is there to be served right away.
	select {
	case <-doorbell:
	default:
		assert.Fail(t, "the doorbell didn't ring for the first request")
	}
	_, _ = queue.dequeue()

	select {
	case <-doorbell:
		assert.Fail(t, "the doorbell rang with nothing queued")
	default:
	}

	// A burst leaves one token, which stands for everything queued.
	queue.enqueue(LinesToRead{Total: 2})
	queue.enqueue(LinesToRead{Total: 3})

	select {
	case <-doorbell:
	default:
		assert.Fail(t, "the doorbell didn't ring")
	}
	select {
	case <-doorbell:
		assert.Fail(t, "the doorbell rang twice for one wake")
	default:
	}
}
